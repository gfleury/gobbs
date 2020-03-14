package pullrequests

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/gfleury/gobbs/common"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// PullRequestRoot cmd root for cobra
var PullRequestRoot = &cobra.Command{
	Use:   "pr",
	Short: "Interact with pull requests",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			List.Run(cmd, args)
		}
	},
}

//List is the cmd implementation for Listing Pull Requests
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
			apiClient := common.APIClient(cmd)
			stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)
			response, err := apiClient.DefaultApi.GetPullRequestsPage(*stashInfo.Project(), *stashInfo.Repo(), opts)
			if err != nil {
				log.Fatal(err.Error())
			}

			pagedPrs, err := bitbucketv1.GetPullRequestsResponse(response)
			if err != nil {
				log.Fatalln(err.Error())
			}
			prs = append(prs, pagedPrs...)

			if response.Values["isLastPage"].(bool) {
				break
			} else {
				opts["start"] = int(response.Values["nextPageStart"].(float64))
			}
		}

		sort.Slice(prs, func(i, j int) bool {
			return prs[i].ID > prs[j].ID
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "State", "Created", "Author", "Short Desc.", "Reviewers"})
		table.SetAutoWrapText(true)
		table.SetAutoFormatHeaders(true)
		table.SetColWidth(40)
		table.SetRowSeparator("-")
		table.SetRowLine(true)
		for _, pr := range prs {
			table.Append([]string{
				fmt.Sprintf("%d", pr.ID),
				pr.State,
				fmt.Sprint(time.Unix(pr.CreatedDate/1000, 0).Format("2006-01-02T15:04:05-0700")),
				pr.Author.User.Name,
				fmt.Sprintf("[%s -> %s] %s", pr.FromRef.DisplayID, pr.ToRef.DisplayID, pr.Title),
				func() (r string) {
					for _, reviewer := range pr.Reviewers {
						r = fmt.Sprintf("%s%s %s\n", r, reviewer.User.Name,
							func() string {
								if reviewer.Approved {
									return "(V)"
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

var (
	listState *string
)

func init() {
	listState = List.Flags().StringP("state", "s", "OPEN", "List only PR's in that state (ALL, OPEN, DECLINED or MERGED)")
}
