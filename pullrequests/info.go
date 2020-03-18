package pullrequests

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
)

// Info is the cmd implementation for Merging Pull Requests
var Info = &cobra.Command{
	Use:     "info pullRequestID",
	Aliases: []string{"inf"},
	Short:   "Info pull requests for repository",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Argument must be a pull request ID. Err: %s", err.Error())
		}

		apiClient, cancel, err := common.APIClient(cmd)
		defer cancel()

		if err != nil {
			log.Fatal(err.Error())
		}

		stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)

		response, err := apiClient.DefaultApi.GetPullRequest(*stashInfo.Project(), *stashInfo.Repo(), prID)

		if netError, ok := err.(net.Error); (!ok || (ok && !netError.Timeout())) && response != nil && response.Response.StatusCode >= http.StatusMultipleChoices {
			common.PrintApiError(response.Values)
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		if response.StatusCode == http.StatusNoContent {
			log.Infof("Pull request ID: %v sucessfully GOT.", prID)
		}

		pr, err := bitbucketv1.GetPullRequestResponse(response)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Get Build Status
		apiClient, cancel, err = common.APIClient(cmd)
		defer cancel()

		if err != nil {
			log.Fatal(err.Error())
		}

		response, err = apiClient.DefaultApi.GetCommitBuildStatuses(pr.FromRef.LatestCommit)

		if netError, ok := err.(net.Error); (!ok || (ok && !netError.Timeout())) && response != nil && response.Response.StatusCode >= http.StatusMultipleChoices {
			common.PrintApiError(response.Values)
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		if response.StatusCode == http.StatusNoContent {
			log.Infof("Build Status for commit ID: %v sucessfully GOT.", pr.FromRef.LatestCommit)
		}

		buildStatuses, err := bitbucketv1.GetBuildStatusesResponse(response)
		if err != nil {
			log.Fatal(err.Error())
		}

		header := []string{"ID", "State", "Updated", "Builds", "Tasks Resolv. / Done", "Short Desc.", "Reviewers"}
		table := common.Table(header)

		table.Append([]string{
			fmt.Sprintf("%d", pr.ID),
			pr.State,
			fmt.Sprint(time.Unix(pr.UpdatedDate/1000, 0).Format("2006-01-02T15:04:05-0700")),
			func() (r string) {
				for _, buildStatus := range buildStatuses {
					r = fmt.Sprintf("%sName: %s\n\tState: %s\n\tURL: %s\n\tDesc.:  %s\n\tAdded: %s\n",
						r,
						buildStatus.Name,
						buildStatus.State,
						buildStatus.Url,
						buildStatus.Description,
						time.Unix(buildStatus.DateAdded/1000, 0).Format("2006-01-02T15:04:05-0700"))
				}
				return
			}(),
			fmt.Sprintf("%d / %d", pr.Properties.OpenTaskCount, pr.Properties.ResolvedTaskCount),
			fmt.Sprintf("[%s -> %s] %s", pr.FromRef.DisplayID, pr.ToRef.DisplayID, pr.Title),
			func() (r string) {
				for _, reviewer := range pr.Reviewers {
					r = fmt.Sprintf("%s%s %s\n", r, reviewer.User.Name,
						func() string {
							if reviewer.Approved {
								return "(A)"
							}
							return "( )"
						}())
				}
				return
			}(),
		})
		table.Render()
	},
}
