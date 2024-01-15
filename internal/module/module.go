package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/eunanhardy/terrapak-action/internal/config"
	"github.com/eunanhardy/terrapak-action/internal/fileutils"
	"github.com/eunanhardy/terrapak-action/internal/github/store"
	"github.com/eunanhardy/terrapak-action/internal/http_client"

	"github.com/fatih/color"
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
	Hash		string `json:"hash"`
	Readme      string `json:"readme"`
}

type UploadRequestBody struct {
	Readme string `json:"readme" form:"readme"`
	Hash   string   `json:"hash" form:"hash"`
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

func Pack(config *config.ModuleConfig)(string,error){
	requestid := uuid.Must(uuid.NewV4())
	localpath := fmt.Sprintf("/tmp/%s/",requestid)
	filepath := fmt.Sprintf("%s/%s.zip",localpath,config.Name)
	err := os.MkdirAll(localpath,os.ModePerm); if err != nil {
		fmt.Println(err)
		return "",err
	}
	err = fileutils.ZipDir(config.Path,filepath); if err != nil {
		fmt.Println(err)
		return "",err
	}

	return filepath,nil
}

func Upload(hostname string,config *config.ModuleConfig) error {
	endpoint := fmt.Sprintf("%s/v1/api/%s/%s/%s/%s/upload",hostname,config.GetNamespace(config.Namespace),config.Name,config.Provider,config.Version)
	readme_path := fmt.Sprintf("%s/README.md",config.Path)
	uploadRequestBody := UploadRequestBody{}
	if fileutils.FileExists(readme_path) {
		bytesReadme, err := os.ReadFile(readme_path); if err != nil {
			color.Red("Error reading README.md")
			return err
		}

		uploadRequestBody.Readme = string(bytesReadme)
	}

	filepath,err := Pack(config); if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	form := multipart.NewWriter(buf)

	file, err := os.Open(filepath); if err != nil {
		return err
	}
	defer file.Close()

	part, err := form.CreateFormFile("file",filepath); if err != nil {
		return err
	}

	_, err = io.Copy(part,file); if err != nil {
		return err
	}

	err = form.WriteField("readme",uploadRequestBody.Readme); if err != nil {
		return err
	}

	err = form.Close(); if err != nil {
		return err
	}

	req, err := http.NewRequest("POST",endpoint,buf); if err != nil {
		return err

	}

	req.Header.Set("Content-Type",form.FormDataContentType())

	client := http_client.Default()
	res, err := client.Do(req); if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		fmt.Println("File Uploaded")
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
			result := store.ResultStore{Name: config.Name, Version: config.Version, Change: "Changes applied"}
			result.Add()
		}
	} else {
		fmt.Printf("[LOG] - Changes detected in %s, but the module is already published, Create a new version to apply changes",config.Name)
		result := store.ResultStore{Name: config.Name, Version: config.Version, Change: "New Version Required"}
		result.Add()
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