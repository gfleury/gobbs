package pullrequests

import (
	"net"
	"net/http"
	"strconv"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
)

// Approve is the cmd implementation for Approving Pull Requests
var Approve = &cobra.Command{
	Use:     "approve pullRequestID",
	Aliases: []string{"app"},
	Short:   "Approve pull requests for repository",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatalf("Argument must be a pull request ID. Err: %s", err.Error())
		}

		apiClient, cancel, err := common.APIClient(cmd)
		defer cancel()

		if err != nil {
			log.Fatal(err.Error())
		}

		stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)

		participant := bitbucketv1.UserWithMetadata{
			User: bitbucketv1.UserWithLinks{
				Name: *stashInfo.Credential().User(),
			},
			Approved: true,
			Status:   "APPROVED",
		}

		response, err := apiClient.DefaultApi.UpdateStatus(*stashInfo.Project(), *stashInfo.Repo(), prID, *stashInfo.Credential().User(), participant)

		if netError, ok := err.(net.Error); (!ok || (ok && !netError.Timeout())) && response != nil && response.Response.StatusCode >= http.StatusMultipleChoices {
			common.PrintApiError(response.Values)
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		if netError, ok := err.(net.Error); (!ok || (ok && !netError.Timeout())) && response.Response.StatusCode >= http.StatusMultipleChoices {
			common.PrintApiError(response.Values)
		}

		participant, err = bitbucketv1.GetUserWithMetadataResponse(response)
		if err != nil {
			log.Fatal(err.Error())
		}

		if response.StatusCode == http.StatusCreated {
			log.Infof("Pull request ID: %v sucessfully APPROVED, last commit %s", prID, participant.LastReviewedCommit)
		}

	},
}
