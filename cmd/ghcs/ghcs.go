package main

// with go modules disabled

import (
	"context"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"net/http"
)

func main() {

	ctx := context.Background()

	// Set here you application ID and installation id
	appId := 39792
	installationId := 1703770 //installation for Analogj/golang_analogj_test

	// Wrap the shared transport for use with defined application and installation IDs
	jwtTransport, err := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, appId, "/Users/jason/repos/gopath/src/github.com/analogj/ghcs/ghcs-test.2019-08-28.private-key.pem")
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	// Use installation transport with jwtClient
	// NewClient returns a new GitHub API jwtClient.
	// If a nil httpClient is provided, http.DefaultClient will be used. To use API methods which require authentication,
	// provide an http.Client that will perform the authentication for you.
	jwtClient := github.NewClient(&http.Client{Transport: jwtTransport})

	appTransport, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appId, installationId, "/Users/jason/repos/gopath/src/github.com/analogj/ghcs/ghcs-test.2019-08-28.private-key.pem")
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	appClient := github.NewClient(&http.Client{Transport: appTransport})

	access_token, err := appTransport.Token()
	if err != nil {
		fmt.Printf("err: %s", err)
	}

	fmt.Printf("Installation access token: %s", access_token)

	// Get org installation
	appService := jwtClient.Apps
	appInst, resp, err := appService.FindRepositoryInstallation(ctx, "AnalogJ", "golang_analogj_test")
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	//installations, resp, err := appService.ListInstallations(ctx, nil)
	fmt.Print(resp)
	fmt.Print(appInst)

	//Get Pull requests

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

}
