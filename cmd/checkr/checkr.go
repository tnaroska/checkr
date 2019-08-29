package main

// with go modules disabled

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/analogj/checkr/pkg/actions"
	"github.com/analogj/checkr/pkg/config"
	"github.com/analogj/checkr/pkg/utils"
	"github.com/analogj/checkr/pkg/version"
	"github.com/fatih/color"
	"github.com/urfave/cli"
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
		Authors: []cli.Author{
			{
				Name:  "Jason Kulatunga",
				Email: "jason@thesparktree.com",
			},
		},
		Before: func(c *cli.Context) error {

			ghcsRepo := "github.com/AnalogJ/checkr"

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

		Commands: []cli.Command{
			{
				Name:  "create",
				Usage: "Create a Github Check Run",
				//UsageText:   "doo - does the dooing",
				Action: func(c *cli.Context) error {
					fmt.Fprintln(c.App.Writer, c.Command.Usage)

					if !c.IsSet("pr") && !c.IsSet("sha") {
						return errors.New("Required flag \"pr\" or \"sha\" is not set")
					} else if c.IsSet("pr") {
						config.Set("pr", c.Int("pr"))
					} else {
						config.Set("sha", c.String("sha"))
					}

					config.Set("org", c.String("org"))
					config.Set("repo", c.String("repo"))

					runAction := actions.RunAction{Config: config}
					return runAction.Create()
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "org, o",
						Usage:    "Github repository owner/organization name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "repo, r",
						Usage:    "Github repository name",
						Required: true,
					},

					&cli.StringFlag{
						Name:  "pr",
						Usage: "Github pull request number (required if sha is not provided)",
					},
					&cli.StringFlag{
						Name:  "sha",
						Usage: "Github pull request head SHA (required if pr is not provided)",
					},

					&cli.StringFlag{
						Name:  "url",
						Usage: "Provide an optional link that will be set in the Check run as the `detail_url`",
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(color.HiRedString("ERROR: %v", err))
	}
}

//func main2() {
//
//
//
//	var createMessage = "new file added"
//	var createBranch = "AnalogJ-patch-3"
//
//	//author ghcs-test[bot] <54657905+ghcs-test[bot]@users.noreply.github.com> 1567062088 +0000
//	var authorName = "ghcs-test"                                //name can be anything
//	var authorEmail = "ghcs-test[bot]@users.noreply.github.com" //[bot]required here, but prefix number not requirerd.
//	var author = github.CommitAuthor{
//		Name:  &authorName,
//		Email: &authorEmail,
//	}
//
//	/// TOOD: write a file to a github branch using the commits api
//	created, resp, err := appClient.Repositories.CreateFile(ctx, "AnalogJ", "golang_analogj_test", "netnew11.txt", &github.RepositoryContentFileOptions{
//		Message:   &createMessage,
//		Content:   []byte("This is my content in a byte array,"),
//		Branch:    &createBranch,
//		Author:    &author,
//		Committer: &author,
//	})
//	if err != nil {
//		fmt.Printf("error: %s", err)
//	}
//	fmt.Print(created)
//
//}
