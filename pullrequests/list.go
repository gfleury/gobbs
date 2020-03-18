package pullrequests

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"time"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
)

var (
	listState *string
)

func init() {
	listState = List.Flags().StringP("state", "s", "OPEN", "List only PR's in that state (ALL, OPEN, DECLINED or MERGED)")
}

// List is the cmd implementation for Listing Pull Requests
var List = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List pull requests for repository",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var prs []bitbucketv1.PullRequest

		opts := map[string]interface{}{
			"state": *listState,
			"limit": 1000,
		}

		for {
			var hasNext bool
			apiClient, cancel, err := common.APIClient(cmd)
			defer cancel()

			if err != nil {
				log.Fatal(err.Error())
			}

			stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)
			response, err := apiClient.DefaultApi.GetPullRequestsPage(*stashInfo.Project(), *stashInfo.Repo(), opts)

			if netError, ok := err.(net.Error); (!ok || (ok && !netError.Timeout())) && response != nil && response.Response.StatusCode >= http.StatusMultipleChoices {
				common.PrintApiError(response.Values)
			}

			if err != nil {
				log.Fatal(err.Error())
			}

			pagedPrs, err := bitbucketv1.GetPullRequestsResponse(response)
			if err != nil {
				log.Fatal(err.Error())
			}
			prs = append(prs, pagedPrs...)

			hasNext, opts["start"] = bitbucketv1.HasNextPage(response)
			if !hasNext {
				break
			}
		}

		sort.Slice(prs, func(i, j int) bool {
			return prs[i].ID > prs[j].ID
		})

		header := []string{"ID", "State (Version)", "Created", "Author", "Short Desc.", "Reviewers"}
		table := common.Table(header)

		for _, pr := range prs {
			table.Append([]string{
				fmt.Sprintf("%d", pr.ID),
				fmt.Sprintf("%s (%v)", pr.State, pr.Version),
				fmt.Sprint(time.Unix(pr.CreatedDate/1000, 0).Format("2006-01-02T15:04:05-0700")),
				pr.Author.User.Name,
				//fmt.Sprintf("%d / %d", pr.Properties.OpenTaskCount, pr.Properties.ResolvedTaskCount),
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
		}
		table.Render()
	},
}
