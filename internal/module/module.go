package module

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/eunanhardy/terrapak-action/internal/config"
	"github.com/eunanhardy/terrapak-action/internal/fileutils"
	"github.com/eunanhardy/terrapak-action/internal/github/store"
	"github.com/eunanhardy/terrapak-action/internal/http_client"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/gofrs/uuid"
)

type ModuleModel struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name        string `json:"name"`
	Provider    string `json:"provider"`
	Namespace 	string `json:"namespace"`
	Version     string `json:"version"`
	DownloadCount  int `json:"download_count"`
	PublishedAt time.Time `json:"published_at"`
	Readme      string `json:"readme"`
}

type UploadRequestBody struct {
	Readme string `json:"readme" form:"readme"`
}

// Returns a ModuleModel struct and the status code of the request
func Read(hostname string,config *config.ModuleConfig) (data ModuleModel, status int, response_err error) {
	endpoint := fmt.Sprintf("%s/v1/api/%s/%s/%s/%s",hostname,config.GetNamespace(config.Namespace),config.Name,config.Provider,config.Version)
	client := http_client.Default()
	resp, err := client.Get(endpoint); if err != nil {
		color.Red("Error reading module")
		return data,0,err
	}

	err1 := json.NewDecoder(resp.Body).Decode(&data); if err1 != nil {
		color.Red("Communication error with terrapak service :: Error decoding response")
		return data,0,err1
	}


	status = resp.StatusCode
	switch(resp.StatusCode){
		case 200:
			return data, status, nil
		case 404:
			return data, status, nil
	}

	return data, status, nil
}

func Upload(hostname string,config *config.ModuleConfig) error {
	endpoint := fmt.Sprintf("%s/v1/api/%s/%s/%s/%s/upload",hostname,config.GetNamespace(config.Namespace),config.Name,config.Provider,config.Version)
	requestid := uuid.Must(uuid.NewV4())
	localpath := fmt.Sprintf("/tmp/%s/",requestid)
	filepath := fmt.Sprintf("%s/%s.zip",localpath,config.Name)
	readme_path := fmt.Sprintf("%s/README.md",config.Path)
	fmt.Println("readme path: ",readme_path)
	uploadRequestBody := UploadRequestBody{}
	err := os.MkdirAll(localpath,os.ModePerm); if err != nil {
		fmt.Println(err)
		return err
	}

	if fileutils.FileExists(readme_path) {
		fmt.Println("README.md exists for ",config.Name)
		bytesReadme, err := os.ReadFile(readme_path); if err != nil {
			color.Red("Error reading README.md")
			return err
		}

		uploadRequestBody.Readme = string(bytesReadme)
	}

	client := resty.New()
	client.SetAuthToken(http_client.DefaultToken())
	err = fileutils.ZipDir(config.Path,filepath); if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := client.R().SetFile("file",filepath).SetBody(uploadRequestBody).Post(endpoint); if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode() == 200 {
		result := store.ResultStore{
			Name: config.Name,
			Version: config.Version,
			Change: "Synced",
		}
		result.Add()
		os.Remove(filepath)
	}
	
	return nil
}

func ModuleDraftCheck(hostname string, config *config.ModuleConfig, data ModuleModel) {
	has_chanaged := fileutils.HasPreviousChanges(config.Path)

	if data.PublishedAt.Year() < 2000 {
		if has_chanaged {
			err := Upload(hostname,config); if err != nil {
				color.Red("[LOG] - Error syncing module changes:%s\n",err)
			}
		}
	} else {
		fmt.Printf("[LOG] - Changes detected in %s, but the module is already published, Create a new version to apply changes",config.Name)
	}

}

func PublishModule(module *config.ModuleConfig){
	gc := config.Default()
	client := http_client.Default()
	endpoint := fmt.Sprintf("%s/v1/api/%s/%s/%s/%s/publish",gc.Terrapak.Hostname,module.GetNamespace(module.Namespace),module.Name,module.Provider,module.Version)
	resp, err := client.Get(endpoint); if err != nil {
		color.Red("failed to publish module")
		os.Exit(1)
	}

	if resp.StatusCode == 200 {
		fmt.Printf("module published: %s",module.Name)
	}

}

func RemoveDraft(module *config.ModuleConfig){
	gc := config.Default()
	endpoint := fmt.Sprintf("%s/v1/api/%s/%s/%s/%s/close",gc.Terrapak.Hostname,module.GetNamespace(module.Namespace),module.Name,module.Provider,module.Version)
	client := http_client.Default()
	resp, err := client.Get(endpoint); if err != nil {
		fmt.Println("failed to remove draft")
		os.Exit(1)
	}

	if resp.StatusCode == 200 {
		fmt.Printf("[DEBUG] - %s draft removed",module.Name)
	}

}