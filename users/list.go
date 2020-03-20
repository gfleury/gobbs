package users

import (
	"context"
	"errors"
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

// List is the cmd implementation for Listing Users
var List = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List users",
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var users []bitbucketv1.User

		opts := map[string]interface{}{
			// "state": *listState,
			"limit": 1000,
		}

		for {
			var hasNext bool
			apiClient, cancel, err := common.APIClient(cmd)
			defer cancel()

			if err != nil {
				log.Critical(err.Error())
				return err
			}

			response, err := apiClient.DefaultApi.GetUsers(opts)

			if netError, ok := err.(net.Error); (!ok || (ok && !netError.Timeout())) &&
				!errors.Is(err, context.Canceled) &&
				!errors.Is(err, context.DeadlineExceeded) &&
				response != nil && response.Response != nil &&
				response.Response.StatusCode >= http.StatusMultipleChoices {
				common.PrintApiError(response.Values)
				return err
			} else if err != nil {
				log.Critical(err.Error())
				return err
			}

			pagedUsers, err := bitbucketv1.GetUsersResponse(response)
			if err != nil {
				log.Critical(err.Error())
				return err
			}
			users = append(users, pagedUsers...)

			hasNext, opts["start"] = bitbucketv1.HasNextPage(response)
			if !hasNext {
				break
			}
		}

		sort.Slice(users, func(i, j int) bool {
			return users[i].Slug < users[j].Slug
		})

		header := []string{"ID", "Slug", "Email", "Name", "Last Login"}
		table := common.Table(header)

		for _, user := range users {
			table.Append([]string{
				fmt.Sprintf("%d", user.ID),
				user.Slug,
				user.EmailAddress,
				user.DisplayName,
				fmt.Sprint(time.Unix(user.LastAuthenticationTimestamp/1000, 0).Format("2006-01-02T15:04:05-0700")),
			})
		}
		table.Render()

		return nil
	},
}
