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
| Module | Version | Info |
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
	input := &gh.IssueComment{Body: &markdown}
	removePreviousComment(client,owner, repo, pr_number)
	_,_, err = client.Issues.CreateComment(ctx,owner, repo, pr_number, input); if err != nil {
		fmt.Println(err)
	}

	// list,_,err := client.Issues.ListComments(ctx,owner, repo, pr_number, nil); if err != nil {
	// 	fmt.Println(err)
	// }
	// currentComment := gh.IssueComment{}
	// fmt.Println(list)
	// for _, comment := range list {
	// 	if strings.Contains(*comment.Body, "Terrapak Sync") {
	// 		currentComment = *comment
	// 		fmt.Println(currentComment.ID)
	// 		return
	// 	}
	// }

	// fmt.Printf("%d:%v",*currentComment.ID,currentComment.Body)
	// if currentComment.Body == nil {
	// 	input := &gh.IssueComment{Body: &markdown}
	// 	_,_, err = client.Issues.CreateComment(ctx,owner, repo, pr_number, input); if err != nil {
	// 		fmt.Println(err)
	// 	}
	// } else {
	// 	markdown = markdown + " \n Edited"
	// 	input := &gh.IssueComment{Body: &markdown}
	// 	_,_, err = client.Issues.EditComment(ctx,owner, repo, *currentComment.ID,input); if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
}

func removePreviousComment(client *gh.Client, owner string, repo string, pr_number int) {
	ctx := context.Background()
	list,_,err := client.Issues.ListComments(ctx,owner, repo, pr_number, nil); if err != nil {
		fmt.Println(err)
	}
	currentComment := gh.IssueComment{}
	fmt.Println(list)
	for _, comment := range list {
		if strings.Contains(*comment.Body, "## Terrapak Sync") {
			currentComment = *comment
			fmt.Println(currentComment.ID)
			return
		}
	}
	if currentComment.ID != nil {
		_,err := client.Issues.DeleteComment(ctx,owner, repo, *currentComment.ID); if err != nil {
			fmt.Println(err)
		}
	}
}

func DisplayPRResults(){
	output := store.Print()
	if output != "" {
		results_template := TABLE_TEMPLATE + output
		AddPRComment(results_template)
	} else {
		results_template := "## Terrapak Sync \n No Changes Detected"
		AddPRComment(results_template)
	}
}