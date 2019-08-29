package actions

import (
	"context"
	"fmt"

	"github.com/analogj/checkr/pkg/client"
	"github.com/analogj/checkr/pkg/config"
	"github.com/google/go-github/github"
)

type RunAction struct {
	Config config.Interface
}

func (r *RunAction) Create() error {

	ctx := context.Background()

	//create jwt Client
	jwtClient, err := client.GetJwtClient(r.Config)

	// get App installation information for this repo.
	appService := jwtClient.Apps
	appInst, resp, err := appService.FindRepositoryInstallation(ctx, r.Config.GetString("org"), r.Config.GetString("repo"))
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}
	fmt.Print(resp)
	fmt.Print(appInst)

	//create app client
	appClient, err := client.GetAppClient(r.Config, int(*appInst.ID))

	// if we don't know the SHA that we're adding the check run to, lets get it from the PR
	var headSha string
	if !r.Config.IsSet("headSha") {
		prData, _, err := appClient.PullRequests.Get(ctx, r.Config.GetString("org"), r.Config.GetString("repo"), r.Config.GetInt("pr"))
		if err != nil {
			fmt.Printf("error: %s", err)
			return err
		}

		headSha = *prData.Head.SHA
		fmt.Printf("Found HEAD SHA: %s", headSha)
	} else {
		headSha = r.Config.GetString("headSha")
	}

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

	check, resp, err := appClient.Checks.CreateCheckRun(ctx, r.Config.GetString("org"), r.Config.GetString("repo"), github.CreateCheckRunOptions{
		Name:       "test-app4",
		HeadSHA:    headSha,
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

	return nil
}
