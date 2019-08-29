package main

// with go modules disabled

import (
	"context"
	"fmt"
	"github.com/analogj/ghcs/pkg/client"
	"github.com/analogj/ghcs/pkg/config"
	"github.com/google/go-github/github"
)

func main() {

	ctx := context.Background()

	appConfig, err := config.Create()

	jwtClient, err := client.GetJwtClient(appConfig)

	// Get org installation
	appService := jwtClient.Apps
	appInst, resp, err := appService.FindRepositoryInstallation(ctx, "AnalogJ", "golang_analogj_test")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	fmt.Print(resp)
	fmt.Print(appInst)

	//Get Pull requests
	appClient, err := client.GetAppClient(appConfig, int(*appInst.ID))

	list, _, err := appClient.PullRequests.List(ctx, "AnalogJ", "golang_analogj_test", &github.PullRequestListOptions{})

	if err != nil {
		fmt.Printf("error: %s", err)
	}

	fmt.Print(list)

	var detailsUrl = "https://github.com/AnalogJ/ghcs"
	var externalID = "externalUId"
	var status = "in_progress"

	var outputTitle = "output TItle2"
	var outputSummary = "output Summary2"
	var outputText = "# output Text2\n This is the output details2"

	var annPath = "README.md"
	var annStartLine = 1
	var annLevel = "notice"
	var annMessage = "message"
	var annTitle = "title"
	var annRawDetails = "Raw Details"

	var annPath2 = "cmd/test.go"
	var annStartLine2 = 5
	var annLevel2 = "notice"
	var annMessage2 = "this file is not changed, but has comment"
	var annTitle2 = "comment on unchanged file"
	var annRawDetails2 = "Raw Details WAZZUP"

	var outputAnnotations = []*github.CheckRunAnnotation{
		{
			Path:            &annPath,
			StartLine:       &annStartLine,
			EndLine:         &annStartLine,
			AnnotationLevel: &annLevel,
			Message:         &annMessage,
			Title:           &annTitle,
			RawDetails:      &annRawDetails,
		},
		{
			Path:            &annPath2,
			StartLine:       &annStartLine2,
			EndLine:         &annStartLine2,
			AnnotationLevel: &annLevel2,
			Message:         &annMessage2,
			Title:           &annTitle2,
			RawDetails:      &annRawDetails2,
		},
	}

	check, resp, err := appClient.Checks.CreateCheckRun(ctx, "AnalogJ", "golang_analogj_test", github.CreateCheckRunOptions{
		Name:       "test-app3",
		HeadSHA:    "1806d6eb80881cd2adad546f5a52e8d6489557cb",
		DetailsURL: &detailsUrl,
		ExternalID: &externalID,
		Status:     &status,
		Output: &github.CheckRunOutput{
			Title:       &outputTitle,
			Summary:     &outputSummary,
			Text:        &outputText,
			Annotations: outputAnnotations,
		},
	})

	if err != nil {
		fmt.Printf("error: %s", err)
	}

	fmt.Print(check)

	var createMessage = "new file added"
	var createBranch = "AnalogJ-patch-3"

	//author ghcs-test[bot] <54657905+ghcs-test[bot]@users.noreply.github.com> 1567062088 +0000
	var authorName = "ghcs-test"                                //name can be anything
	var authorEmail = "ghcs-test[bot]@users.noreply.github.com" //[bot]required here, but prefix number not requirerd.
	var author = github.CommitAuthor{
		Name:  &authorName,
		Email: &authorEmail,
	}

	/// TOOD: write a file to a github branch using the commits api
	created, resp, err := appClient.Repositories.CreateFile(ctx, "AnalogJ", "golang_analogj_test", "netnew11.txt", &github.RepositoryContentFileOptions{
		Message:   &createMessage,
		Content:   []byte("This is my content in a byte array,"),
		Branch:    &createBranch,
		Author:    &author,
		Committer: &author,
	})
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	fmt.Print(created)

}
