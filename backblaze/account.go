package backblaze

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/packago/config"
)

// Account holds backblaze specific data.
type Account struct {
	AbsoluteMinimumPartSize int    `json:"absoluteMinimumPartSize"`
	AccountID               string `json:"accountId"`
	Allowed                 struct {
		BucketID     string      `json:"bucketId"`
		BucketName   string      `json:"bucketName"`
		Capabilities []string    `json:"capabilities"`
		NamePrefix   interface{} `json:"namePrefix"`
	} `json:"allowed"`
	APIURL              string `json:"apiUrl"`
	AuthorizationToken  string `json:"authorizationToken"`
	DownloadURL         string `json:"downloadUrl"`
	RecommendedPartSize int    `json:"recommendedPartSize"`
}

// authorizeAccount retrieves backblaze account data connected to a
// configured keyID and applicationKey
func authorizeAccount() (Account, error) {
	var a Account
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/b2api/v2/b2_authorize_account", config.File().GetString("backblaze.rootUrl")), nil)
	if err != nil {
		return a, err
	}
	req.SetBasicAuth(config.File().GetString("backblaze.keyID"), config.File().GetString("backblaze.applicationKey"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return a, err
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&a); err != nil {
		return a, err
	}
	return a, nil
}
