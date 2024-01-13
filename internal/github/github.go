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
Terrapak has detected changes in the following modules:
| Module | Version | Action |
| :---: | :---: | :---: |
`

func AddPRComment(markdown string) {
	token := os.Getenv("INPUT_GITHUB_TOKEN")
	repoowner := os.Getenv("INPUT_REPO_NAME")
	owner := strings.Split(repoowner, "/")[0]
	repo := strings.Split(repoowner, "/")[1]
	pr_number, err := strconv.Atoi(os.Getenv("INPUT_ISSUE_NUMBER")); if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	client := gh.NewTokenClient(ctx, token)
	list,_,err := client.Issues.ListComments(ctx,owner, repo, pr_number, nil); if err != nil {
		fmt.Println(err)
	}
	currentComment := gh.IssueComment{}
	for _, comment := range list {
		if strings.Contains(*comment.Body, "Terrapak Sync") {
			currentComment = *comment
			fmt.Println(currentComment.ID)
			return
		}
	}

	if currentComment.Body == nil {
		input := &gh.IssueComment{Body: &markdown}
		_,_, err = client.Issues.CreateComment(ctx,owner, repo, pr_number, input); if err != nil {
			fmt.Println(err)
		}
	} else {
		markdown = markdown + " \n Edited"
		input := &gh.IssueComment{Body: &markdown}
		_,_, err = client.Issues.EditComment(ctx,owner, repo, *currentComment.ID,input); if err != nil {
			fmt.Println(err)
		}
	}

	
}

func DisplayPRResults(){
	results_template := TABLE_TEMPLATE + store.Print()
	fmt.Println(results_template)
	AddPRComment(results_template)
}