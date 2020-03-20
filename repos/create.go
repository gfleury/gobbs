package repos

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
)

var (
	forkable, private *bool
)

func init() {
	forkable = Create.Flags().BoolP("forkable", "f", true, "Is repository forkable")
	private = Create.Flags().BoolP("private", "j", false, "Is repository private")
}

// Create is the cmd implementation for Creating Pull Requests
var Create = &cobra.Command{
	Use:     "create repoName",
	Aliases: []string{"cr"},
	Short:   "Create repository",
	Args:    cobra.MinimumNArgs(1),
	PreRunE: mustHaveProject,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, cancel, err := common.APIClient(cmd)
		defer cancel()

		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)

		repository := bitbucketv1.Repository{
			Name:     args[0],
			Forkable: *forkable,
			Public:   !*private,
		}

		response, err := apiClient.DefaultApi.CreateRepository(*stashInfo.Project(), repository)

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

		repo, err := bitbucketv1.GetRepositoryResponse(response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		header := []string{"ID", "Slug", "Links", "Status"}
		table := common.Table(header)

		table.Append([]string{
			fmt.Sprintf("%d", repo.ID),
			repo.Slug,
			func() (r string) {
				if repo.Links != nil {
					for _, link := range repo.Links.Clone {
						r = fmt.Sprintf("%s%s: %s\n", r, link.Name, link.Href)
					}
				}
				return
			}(),
			repo.StatusMessage,
		})
		table.Render()

		return nil
	},
}
