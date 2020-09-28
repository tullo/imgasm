package file

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dgraph-io/badger/v2"
	"github.com/go-chi/chi"
	"github.com/h2non/bimg"
	"github.com/h2non/filetype"
	"github.com/packago/config"
	"github.com/tullo/cookie"
	"github.com/tullo/imgasm/backblaze"
	"github.com/tullo/imgasm/db"
	"github.com/tullo/imgasm/models"
	"github.com/tullo/imgasm/ui/templates"
)

const maxFileSize int64 = 1024 * 1024 * 10 // 10 MB

type File struct {
	log *log.Logger
}

func New() *File {
	f := File{
		log: log.New(os.Stdout, "IMGASM : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
	}
	return &f
}

// Retrieve knows how to load a file from the backblaze cloud storrage.
func (File) Retrieve(w http.ResponseWriter, r *http.Request) {
	sess, _ := cookie.GetSession(r, config.File().GetString("session.key"))
	commonData := templates.ReadCommonData(w, r)
	fileDataBytes, err := db.BadgerDB.Get([]byte(chi.URLParam(r, "fileid")))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			commonData.MetaTitle = "404"
			templates.Render(w, "not-found.html", map[string]interface{}{
				"Common": commonData,
			})
			return
		}
		sess.AddFlash(err.Error())
		sess.Save(r, w)

	}
	var fileData models.FileData
	if err = json.Unmarshal(fileDataBytes, &fileData); err != nil {
		sess.AddFlash(err.Error())
		sess.Save(r, w)
	}

	fileServerURL := config.File().GetString("FileServerURL")
	templates.Render(w, "file.html", map[string]interface{}{
		"Common":        commonData,
		"Filename":      fileData.Filename,
		"FileServerURL": fileServerURL,
	})
}

// Upload knows how to save a file to the backblaze cloud storrage.
func (f File) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	file, _, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			err = errors.New("remember to select a file to upload")
		}
		renderTemplateWithError(w, r, err, "index.html")
		return
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	f.log.Println("trouble reading file:", err)
	if len(buf) == 0 {
		renderTemplateWithError(w, r, errors.New("trouble reading file, empty body"), "index.html")
		return
	}

	// Get file type and check if it's supported
	kind, err := filetype.Match(buf)
	if err != nil {
		renderTemplateWithError(w, r, err, "index.html")
		return
	}
	if !filetype.IsImage(buf) {
		renderTemplateWithError(w, r, errors.New("this filetype is currently not supported"), "index.html")
		return
	}

	// check if the file has already been processed and uploaded
	originalMD5 := fmt.Sprintf("%x", md5.Sum(buf))
	if savedMD5Bytes, err := db.BadgerDB.Get([]byte(originalMD5)); err != badger.ErrKeyNotFound {
		var fileMD5 models.FileMD5
		if err = json.Unmarshal(savedMD5Bytes, &fileMD5); err != nil {
			renderTemplateWithError(w, r, err, "index.html")
			return
		}
		f.log.Println("original:", originalMD5)
		f.log.Println("processed:", fileMD5.Processed)
		if originalMD5 != fileMD5.Processed {
			if savedMD5Bytes, err := db.BadgerDB.Get([]byte(fileMD5.Processed)); err != badger.ErrKeyNotFound {
				if err = json.Unmarshal(savedMD5Bytes, &fileMD5); err != nil {
					renderTemplateWithError(w, r, err, "index.html")
					return
				}
			}
		}
		fileData, err := saveFileData(fmt.Sprintf("%x.%s", fileMD5.Processed, fileMD5.Extension))
		if err != nil {
			renderTemplateWithError(w, r, err, "index.html")
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/%s", fileData.ID), http.StatusSeeOther)
		return
	}

	// convert image to PNG unless it's a JPEG
	if kind.MIME.Type != "image/png" && kind.MIME.Type != "image/jpeg" {
		buf, err = bimg.NewImage(buf).Convert(bimg.PNG)
		if err != nil {
			renderTemplateWithError(w, r, err, "index.html")
			return
		}
	}

	// resize the image if its width or heigh is greater than 1500
	size, err := bimg.NewImage(buf).Size()
	if err != nil {
		renderTemplateWithError(w, r, err, "index.html")
		return
	}
	if size.Height > 1500 || size.Width > 1500 {
		fitWidth, fitHeigh := calculateFitDimension(size.Width, size.Height, 1500, 1500)
		buf, err = bimg.NewImage(buf).Resize(fitWidth, fitHeigh)
		if err != nil {
			renderTemplateWithError(w, r, err, "index.html")
			return
		}
		kind, err = filetype.Match(buf)
		if err != nil {
			renderTemplateWithError(w, r, err, "index.html")
			return
		}
	}

	// upload file to backblaze
	image := models.File{
		Body:      buf,
		MD5Hash:   fmt.Sprintf("%x", md5.Sum(buf)),
		MimeType:  kind.MIME.Value,
		Extension: kind.Extension,
	}
	if err = backblaze.Upload(f.log, image); err != nil {
		renderTemplateWithError(w, r, err, "index.html")
		return
	}

	fileData, err := saveFileData(fmt.Sprintf("%x.%s", image.MD5Hash, image.Extension))
	if err != nil {
		renderTemplateWithError(w, r, err, "index.html")
		return
	}
	saveMD5Hashes(originalMD5, image.MD5Hash, image.Extension)

	http.Redirect(w, r, fmt.Sprintf("/%s", fileData.ID), http.StatusSeeOther)
}

func renderTemplateWithError(w http.ResponseWriter, r *http.Request, err error, template string) {
	sess, _ := cookie.GetSession(r, config.File().GetString("session.key"))
	sess.AddFlash(err.Error())
	sess.Save(r, w)
	commonData := templates.ReadCommonData(w, r)
	templates.Render(w, template, map[string]interface{}{
		"Common": commonData,
	})
}
