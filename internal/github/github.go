package github

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eunanhardy/terrapak-action/internal/github/store"
	gh "github.com/google/go-github/v58/github"
)

const TABLE_TEMPLATE = `## Terrapak Sync
Changes detected in the following modules.
| Module | Version | Action |
| :---: | :---: | :---: |\n
`

func AddPRComment(markdown string) {
	token := os.Getenv("INPUT_GITHUB_TOKEN")
	repoowner := os.Getenv("INPUT_REPO_NAME")
	owner := strings.Split(repoowner, "/")[0]
	repo := strings.Split(repoowner, "/")[1]
	pr_number, err := strconv.Atoi(os.Getenv("INPUT_ISSUE_NUMBER")); if err != nil {
		fmt.Println(err)
	}
	//endpoint := fmt.Sprintf("https://api.github.com/repos/%s/issues/%s/comments", repo, pr_number)
	//fmt.Println(endpoint)
	input := &gh.IssueComment{Body: &markdown}
	ctx := context.Background()
	client := gh.NewTokenClient(ctx, token)
	comment,resp, err := client.Issues.CreateComment(ctx,owner, repo, pr_number, input); if err != nil {
		fmt.Println(err)
	}
	fmt.Println(comment)
	fmt.Println(resp.Status)
	// client := http_client.New(token)
	// resp, err := client.Post(endpoint,"application/json",strings.NewReader(body)); if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(resp.Status)
	// fmt.Println(strings.NewReader(body))
	// req, err := http.NewRequest("POST", endpoint, strings.NewReader(body)); if err != nil {
	// 	fmt.Println(err)
	// }
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Accept", "application/vnd.github.v3+json")
	// req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	// fmt.Println(req)
	// client := &http.Client{}
	// resp, err := client.Do(req); if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(resp.Status)


	// defer resp.Body.Close()

}

func DisplayPRResults(){
	results_template := TABLE_TEMPLATE + store.Print()
	fmt.Println(results_template)
	AddPRComment(results_template)
}