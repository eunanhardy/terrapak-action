package github

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

const TABLE_TEMPLATE = `## Terrapak Sync
Changes detected in the following modules.
| Module | Version | Action |
| --- | --- | --- |\n
`

func AddPRComment(markdown string) {
	token := os.Getenv("INPUT_GITHUB_TOKEN")
	owner := os.Getenv("INPUT_REPO_NAME")
	pr_number := os.Getenv("INPUT_ISSUE_NUMBER")
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s/comments", owner, pr_number)
	body := fmt.Sprintf(`{"body": "%s"}`, markdown)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(body)); if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req); if err != nil {
		fmt.Println(err)
	}


	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		
		color.Green("Comment added successfully")
	}
}

func DisplayPRResults(){
	results_template := TABLE_TEMPLATE+Print()
	AddPRComment(results_template)
}