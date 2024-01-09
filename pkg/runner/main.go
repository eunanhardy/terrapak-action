package runner

import (
	"fmt"
	"os"

	ms "github.com/eunanhardy/terrapak-action/internal/module"

	"github.com/eunanhardy/terrapak-action/internal/config"
	"github.com/eunanhardy/terrapak-action/internal/fileutils"
	"github.com/eunanhardy/terrapak-action/internal/github"
	"github.com/eunanhardy/terrapak-action/internal/http_client"
	"github.com/eunanhardy/terrapak-action/internal/module"

	"github.com/fatih/color"
)

func Run(){
	token, set := os.LookupEnv("INPUT_TOKEN"); if !set {
		fmt.Println("[ERROR] - Terrapak token not found")
		return
	}
	http_client.New(token)

	color.NoColor = false
	action, set := os.LookupEnv("INPUT_ACTION"); if !set {
		fmt.Println("[ERROR] - No action found")
		return
	}

	switch(action){
		case "sync":
			onOpenedPR()
			break;
		case "merged":
			onMergedPR()
			break;
		case "closed":
			onClosedPR()
			break;
	}
}

func onOpenedPR(){

	current_config, err := config.Load(); if err != nil {
		fmt.Println("[ERROR] - Could not load config file")
		os.Exit(1)
	}

	if !healthCheck(current_config.Terrapak.Hostname) {
		fmt.Println("[ERROR] - Terrapak cannot be reached")
		os.Exit(1)
	}
	github.NewDefaultResultSet()
	for _, mod := range current_config.Modules {
		changed := fileutils.HasChanges(mod.Path)
		fmt.Printf("[DEBUG] - %s has changes: %t\n",mod.Name,changed)
		if changed {
			data , status, err := ms.Read(current_config.Terrapak.Hostname,&mod); if err != nil {
				fmt.Println(err)
			}

			switch(status){
				case 404:
					fmt.Println("Module not found")
					module.Upload(current_config.Terrapak.Hostname,&mod)
				break;
				case 200:
					module.ModuleDraftCheck(current_config.Terrapak.Hostname,&mod,data)
				break;

			}

		}
	}

	github.DisplayPRResults()
}

func onMergedPR(){
	gc := config.Default()
	current_config, err := config.Load(); if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !healthCheck(current_config.Terrapak.Hostname) {
		fmt.Println("[ERROR] - Terrapak service is down")
		os.Exit(1)
	}

	for _, mod := range current_config.Modules {
		module,status, err := ms.Read(current_config.Terrapak.Hostname,&mod); if err != nil {
			fmt.Println(err)
		}

		if status == 200 {
			if module.PublishedAt.Year() < 2000 {
				comment := fmt.Sprintf("Module Published: `%s/%s/%s/%s/%s`",gc.Terrapak.Hostname,mod.GetNamespace(mod.Namespace),mod.Name,mod.Provider,mod.Version)
				fmt.Println("[DEBUG] - Module is a draft, publishing module")
				ms.PublishModule(&mod)
				github.AddPRComment(comment)
			}
		}
	}

}

func onClosedPR(){
	current_config, err := config.Load(); if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !healthCheck(current_config.Terrapak.Hostname) {
		fmt.Println("[ERROR] - Terrapak service is down, cannot cleanup module")
		os.Exit(1)
	}

	for _, mod := range current_config.Modules {
		module,status, err := ms.Read(current_config.Terrapak.Hostname,&mod); if err != nil {
			fmt.Println(err)
		}

		if status == 200 {
			if module.PublishedAt.Year() < 2000 {
				color.Magenta("[DELETE] - Module is a draft, removing module")
				ms.RemoveDraft(&mod)
			}
		}
	}

}

func healthCheck(hostname string) bool {
	endpont := fmt.Sprintf("%s/ping",hostname)
	client := http_client.Default()
	resp, err := client.Get(endpont); if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode == 200 {
		color.Green("[INFO] - Terrapak service is up")
		return true
	}

	return false
}