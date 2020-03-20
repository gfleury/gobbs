package pullrequests

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
)

var (
	description, title *string
)

func init() {
	description = Create.Flags().StringP("description", "d", "", "PR description")
	title = Create.Flags().StringP("title", "T", "", "PR title")
}

// Create is the cmd implementation for Creating Pull Requests
var Create = &cobra.Command{
	Use:     "create fromBranch toBranch",
	Aliases: []string{"cr"},
	Short:   "Create pull requests for repository",
	Args:    cobra.MinimumNArgs(2),
	PreRunE: mustHaveProjectRepo,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, cancel, err := common.APIClient(cmd)
		defer cancel()

		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)

		if *title == "" {
			*title = titleFromBranch(args[0], args[1])
		}

		if *description == "" {
			*description = titleFromBranch(args[0], args[1])
		}

		pr := bitbucketv1.PullRequest{
			Title:       *title,
			Description: *description,
			Version:     0,
			State:       "OPEN",
			Open:        true,
			Closed:      false,
			FromRef: bitbucketv1.PullRequestRef{
				ID: fmt.Sprintf("refs/heads/%s", args[0]),
				Repository: bitbucketv1.Repository{
					Slug: *stashInfo.Repo(),
					Project: &bitbucketv1.Project{
						Key: *stashInfo.Project(),
					},
				},
			},
			ToRef: bitbucketv1.PullRequestRef{
				ID: fmt.Sprintf("refs/heads/%s", args[1]),
				Repository: bitbucketv1.Repository{
					Slug: *stashInfo.Repo(),
					Project: &bitbucketv1.Project{
						Key: *stashInfo.Project(),
					},
				},
			},
			Locked: false,
			//"reviewers": reviewers,
		}

		response, err := apiClient.DefaultApi.CreatePullRequest(*stashInfo.Project(), *stashInfo.Repo(), pr)

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

		pr, err = bitbucketv1.GetPullRequestResponse(response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		header := []string{"ID", "State", "Created", "Short Desc.", "Reviewers"}
		table := common.Table(header)

		table.Append([]string{
			fmt.Sprintf("%d", pr.ID),
			pr.State,
			fmt.Sprint(time.Unix(pr.CreatedDate/1000, 0).Format("2006-01-02T15:04:05-0700")),
			//pr.Author.User.Name,
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

func titleFromBranch(from, to string) string {
	return fmt.Sprintf("Merge '%s' into '%s'", from, to)
}
