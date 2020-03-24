package pullrequests

import (
	"context"
	"errors"
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

// Merge is the cmd implementation for Merging Pull Requests
var Merge = &cobra.Command{
	Use:     "merge pullRequestID",
	Aliases: []string{"mer"},
	Short:   "Merge pull requests for repository",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Critical("Argument must be a pull request ID. Err: %s", err.Error())
			return err
		}

		apiClient, cancel, err := common.APIClient(cmd)
		defer cancel()

		if err != nil {
			log.Critical("Argument must be a pull request ID. Err: %s", err.Error())
			return err
		}

		stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)
		err = mustHaveProjectRepo(stashInfo)
		if err != nil {
			return err
		}

		opts := map[string]interface{}{
			"version": int32(*prVersion),
		}

		response, err := apiClient.DefaultApi.Merge(*stashInfo.Project(), *stashInfo.Repo(), prID, opts, nil, []string{"application/json"})

		if netError, ok := err.(net.Error); (!ok || (ok && !netError.Timeout())) &&
			!errors.Is(err, context.Canceled) &&
			!errors.Is(err, context.DeadlineExceeded) &&
			response != nil && response.Response != nil &&
			response.Response.StatusCode >= http.StatusMultipleChoices {
			common.PrintApiError(response.Values)
			cmd.SilenceUsage = true
			log.Debugf(err.Error())
			return fmt.Errorf("Unable to process request, API Error")
		} else if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		if response.StatusCode == http.StatusNoContent {
			log.Infof("Pull request ID: %v sucessfully MERGED.", prID)
		}

		pr, err := bitbucketv1.GetPullRequestResponse(response)
		if err != nil {
			log.Critical("Argument must be a pull request ID. Err: %s", err.Error())
			return err
		}

		header := []string{"ID", "State", "Updated", "Result", "Tasks Resolv. / Done", "Short Desc.", "Reviewers"}
		table := common.Table(header)

		table.Append([]string{
			fmt.Sprintf("%d", pr.ID),
			pr.State,
			fmt.Sprint(time.Unix(pr.UpdatedDate/1000, 0).Format("2006-01-02T15:04:05-0700")),
			pr.Properties.MergeResult.Outcome,
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

		return nil
	},
}
