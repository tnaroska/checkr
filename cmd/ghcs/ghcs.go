package main

// with go modules disabled

import (
	"context"
	"fmt"
	"github.com/analogj/drawbridge/pkg/actions"
	"github.com/analogj/drawbridge/pkg/project"
	"github.com/analogj/ghcs/pkg/client"
	"github.com/analogj/ghcs/pkg/config"
	"github.com/analogj/ghcs/pkg/utils"
	"github.com/analogj/ghcs/pkg/version"
	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"strings"
	"time"
)

var goos string
var goarch string

func main() {

	config, err := config.Create()
	if err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		os.Exit(1)
	}

	app := &cli.App{
		Name:     "ghcs",
		Usage:    "Github Check Suite CLI",
		Version:  version.VERSION,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Jason Kulatunga",
				Email: "jason@thesparktree.com",
			},
		},
		Before: func(c *cli.Context) error {

			ghcsRepo := "github.com/AnalogJ/ghcs"

			var versionInfo string
			if len(goos) > 0 && len(goarch) > 0 {
				versionInfo = fmt.Sprintf("%s.%s-%s", goos, goarch, version.VERSION)
			} else {
				versionInfo = fmt.Sprintf("dev-%s", version.VERSION)
			}

			subtitle := ghcsRepo + utils.LeftPad2Len(versionInfo, " ", 65-len(ghcsRepo))

			color.New(color.FgGreen).Fprintf(c.App.Writer, fmt.Sprintf(utils.StripIndent(
				`
			  oooooooo8 oooo                             oooo                    
			o888     88  888ooooo   ooooooooo8  ooooooo   888  ooooo oo oooooo   
			888          888   888 888oooooo8 888     888 888o888     888    888 
			888o     oo  888   888 888        888         8888 88o    888        
			 888oooo88  o888o o888o  88oooo888  88ooo888 o888o o888o o888o
			%s

			`), subtitle))

			return nil
		},

		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a Github Check run",
				//UsageText:   "doo - does the dooing",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					projectList, err := project.CreateProjectListFromProvidedAnswers(config)
					if err != nil {
						return err
					}

					answerData := map[string]interface{}{}
					if projectList.Length() > 0 && utils.StdinQueryBoolean(fmt.Sprintf("Would you like to create a Drawbridge config using preconfigured answers? (%v available). [yes/no]", projectList.Length())) {

						answerData, err = projectList.Prompt("Enter number to base your configuration from")
						if err != nil {
							return err
						}
					}

					//extend current answerData with CLI provided options.
					cliAnswers, err := createFlagHandler(config, answerData, c.FlagNames(), c)
					if err != nil {
						return err
					}

					createAction := actions.CreateAction{Config: config}
					return createAction.Start(cliAnswers, c.Bool("dryrun"))
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(color.HiRedString("ERROR: %v", err))
	}
}

func main2() {

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
