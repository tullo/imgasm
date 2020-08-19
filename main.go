package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/packago/config"
	"github.com/tullo/imgasm.com/file"
	"github.com/tullo/imgasm.com/templates"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

const (
	viewLimit   = "150-H"
	uploadLimit = "100-D"
)

func main() {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	viewStore := memory.NewStore()
	rate, err := limiter.NewRateFromFormatted(viewLimit)
	if err != nil {
		panic(err)
	}
	view := stdlib.NewMiddleware(limiter.New(viewStore, rate, limiter.WithTrustForwardHeader(true)))
	view.OnLimitReached = rateLimitHandler

	file := file.New()

	r.Group(func(r chi.Router) {
		r.Use(view.Handler)
		r.Get("/", index)
		// r.Get("/about", about)
		r.Get("/protect-your-privacy", privacy)
		r.Get("/privacy-policy", policy)
		r.NotFound(notFound)
		r.Get("/{fileid}", file.Retrieve)
	})

	us := memory.NewStore()
	ur, err := limiter.NewRateFromFormatted(uploadLimit)
	if err != nil {
		panic(err)
	}
	upload := stdlib.NewMiddleware(limiter.New(us, ur, limiter.WithTrustForwardHeader(true)))
	upload.OnLimitReached = rateLimitHandler
	r.Group(func(r chi.Router) {
		r.Use(upload.Handler)
		r.Post("/", file.Upload)
	})

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filesDir := filepath.Join(currentDir, "static")
	fileServer(r, "/static", http.Dir(filesDir))

	switch config.File().GetString("environment") {
	case "development":
		panic(http.ListenAndServe(config.File().GetString("port.development"), r))
	case "production":
		panic(http.ListenAndServe(config.File().GetString("port.production"), r))
	default:
		panic("Environment not set")
	}
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}
	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	commonData.MetaTitle = "404"
	templates.Render(w, "not-found.html", map[string]interface{}{
		"Common": commonData,
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	templates.Render(w, "index.html", map[string]interface{}{
		"Common": commonData,
	})
}

func about(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	commonData.MetaTitle = "About"
	templates.Render(w, "about.html", map[string]interface{}{
		"Common": commonData,
	})
}

func privacy(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	commonData.MetaTitle = "Protect Your Privacy"
	templates.Render(w, "protect-your-privacy.html", map[string]interface{}{
		"Common": commonData,
	})
}

func policy(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	commonData.MetaTitle = "Privacy Policy"
	templates.Render(w, "privacy-policy.html", map[string]interface{}{
		"Common": commonData,
	})
}

func rateLimitHandler(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	commonData.MetaTitle = "Rate limit reached"
	templates.Render(w, "rate-limit.html", map[string]interface{}{
		"Common":      commonData,
		"UploadLimit": strings.Replace(uploadLimit, "-D", "", -1),
		"ViewLimit":   strings.Replace(viewLimit, "-H", "", -1),
	})
}
