package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/analogj/checkr/pkg/client"
	"github.com/analogj/checkr/pkg/config"
	"github.com/google/go-github/v50/github"
)

type RunAction struct {
	Config config.Interface
}

func (r *RunAction) Create(payloadPath string) error {

	payloadData, err := ioutil.ReadFile(payloadPath)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}

	var checkRun github.CreateCheckRunOptions

	err = json.Unmarshal(payloadData, &checkRun)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}

	ctx := context.Background()

	//create jwt Client
	jwtClient, err := client.GetJwtClient(r.Config)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}

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
	appClient, err := client.GetAppClient(r.Config, *appInst.ID)
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}

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

	checkRun.HeadSHA = headSha
	if r.Config.IsSet("details_url") {
		var detailsUrl = r.Config.GetString("details_url")
		checkRun.DetailsURL = &detailsUrl
	}

	if len(checkRun.Output.Annotations) > 50 {
		//we need to chunk the check run annotations into groups of 50

		chunkedAnnotations := chunkAnnotations(checkRun.Output.Annotations, 50)

		//get the first chunk
		firstAnnotationChunk, chunkedAnnotations := chunkedAnnotations[0], chunkedAnnotations[1:]

		//get the last chunk
		lastAnnotationChunk, chunkedAnnotations := chunkedAnnotations[len(chunkedAnnotations)-1], chunkedAnnotations[:len(chunkedAnnotations)-1]

		//save original values
		origCompletedAt := checkRun.GetCompletedAt()
		origStatus := checkRun.GetStatus()
		origConclusion := checkRun.GetConclusion()

		//set overrides
		inProgressStatus := "in_progress"
		checkRun.CompletedAt = nil
		checkRun.Status = &inProgressStatus
		checkRun.Conclusion = nil

		checkRun.Output.Annotations = firstAnnotationChunk
		firstCheckRun, err := createCheckRun(appClient, ctx, r.Config.GetString("org"), r.Config.GetString("repo"), checkRun)

		if err != nil {
			return err
		}

		for _, chunk := range chunkedAnnotations {
			updateCheckRunOptions := github.UpdateCheckRunOptions{
				Name:        checkRun.Name,
				DetailsURL:  checkRun.DetailsURL,
				ExternalID:  checkRun.ExternalID,
				Status:      checkRun.Status,
				Conclusion:  checkRun.Conclusion,
				CompletedAt: nil,
				Output:      checkRun.Output,
			}

			updateCheckRunOptions.Output.Annotations = chunk

			_, err := updateCheckRun(appClient, ctx, r.Config.GetString("org"), r.Config.GetString("repo"), firstCheckRun.GetID(), updateCheckRunOptions)

			if err != nil {
				return err
			}
		}

		//now that we've finished uploading all chunks, lets submit the last peice and restore the correct complete status info.

		lastCheckRunOptions := github.UpdateCheckRunOptions{
			Name:        checkRun.Name,
			DetailsURL:  checkRun.DetailsURL,
			ExternalID:  checkRun.ExternalID,
			Status:      &origStatus,
			Conclusion:  &origConclusion,
			CompletedAt: &origCompletedAt,
			Output:      checkRun.Output,
		}
		lastCheckRunOptions.Output.Annotations = lastAnnotationChunk

		_, err = updateCheckRun(appClient, ctx, r.Config.GetString("org"), r.Config.GetString("repo"), firstCheckRun.GetID(), lastCheckRunOptions)

		return err

	} else {
		_, err := createCheckRun(appClient, ctx, r.Config.GetString("org"), r.Config.GetString("repo"), checkRun)
		return err
	}
}

func createCheckRun(appClient *github.Client, ctx context.Context, owner, repo string, checkRun github.CreateCheckRunOptions) (*github.CheckRun, error) {
	createdCheckRun, _, err := appClient.Checks.CreateCheckRun(ctx, owner, repo, checkRun)
	return createdCheckRun, err
}

func updateCheckRun(appClient *github.Client, ctx context.Context, owner, repo string, checkRunId int64, checkRunUpdate github.UpdateCheckRunOptions) (*github.CheckRun, error) {
	updatedCheckRun, _, err := appClient.Checks.UpdateCheckRun(ctx, owner, repo, checkRunId, checkRunUpdate)
	return updatedCheckRun, err
}

func chunkAnnotations(annotations []*github.CheckRunAnnotation, chunkSize int) [][]*github.CheckRunAnnotation {
	var divided [][]*github.CheckRunAnnotation

	for i := 0; i < len(annotations); i += chunkSize {
		end := i + chunkSize

		if end > len(annotations) {
			end = len(annotations)
		}

		divided = append(divided, annotations[i:end])
	}

	return divided
}
