package github

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/eunanhardy/terrapak-action/internal/github/store"
)

const TABLE_TEMPLATE = `## Terrapak Sync
Changes detected in the following modules.
| Module | Version | Action |
| --- | --- | --- |\n
`

func AddPRComment(markdown string) {
	token := os.Getenv("INPUT_GITHUB_TOKEN")
	repo := os.Getenv("INPUT_REPO_NAME")
	pr_number := os.Getenv("INPUT_ISSUE_NUMBER")
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s/comments", repo, pr_number)
	fmt.Println(endpoint)
	body := fmt.Sprintf(`{"body": "%s"}`, markdown)
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(body)); if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	fmt.Println(req)
	client := &http.Client{}
	resp, err := client.Do(req); if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Status)


	defer resp.Body.Close()

}

func DisplayPRResults(){
	results_template := TABLE_TEMPLATE + store.Print()
	fmt.Println(results_template)
	AddPRComment(results_template)
}