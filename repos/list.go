package repos

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sort"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
)

// var (
// 	listState *string
// )

// func init() {
// 	listState = List.Flags().StringP("state", "s", "OPEN", "List only PR's in that state (ALL, OPEN, DECLINED or MERGED)")
// }

// List is the cmd implementation for Listing Repositories
var List = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List repositories",
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var repos []bitbucketv1.Repository

		opts := map[string]interface{}{
			// "state": *listState,
			"limit": 1000,
		}

		for {
			var hasNext bool
			apiClient, cancel, err := common.APIClient(cmd)
			defer cancel()

			if err != nil {
				cmd.SilenceUsage = true
				return err
			}

			stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)
			err = mustHaveProject(stashInfo)
			if err != nil {
				return err
			}

			response, err := apiClient.DefaultApi.GetRepositoriesWithOptions(*stashInfo.Project(), opts)

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

			pagedRepos, err := bitbucketv1.GetRepositoriesResponse(response)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}
			repos = append(repos, pagedRepos...)

			hasNext, opts["start"] = bitbucketv1.HasNextPage(response)
			if !hasNext {
				break
			}
		}

		sort.Slice(repos, func(i, j int) bool {
			return repos[i].Slug > repos[j].Slug
		})

		header := []string{"ID", "Slug", "Links", "Status"}
		table := common.Table(header)

		for _, repo := range repos {
			table.Append([]string{
				fmt.Sprintf("%d", repo.ID),
				repo.Slug,
				func() (r string) {
					for _, link := range repo.Links.Clone {
						r = fmt.Sprintf("%s%s: %s\n", r, link.Name, link.Href)
					}
					return
				}(),
				repo.StatusMessage,
			})
		}
		table.Render()

		return nil
	},
}
