package pullrequests

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/gfleury/gobbs/common"
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
		apiClient := common.APIClient(cmd)
		stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)
		response, err := apiClient.DefaultApi.GetPullRequestsPage(*stashInfo.Project(), *stashInfo.Repo(), map[string]interface{}{})
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(response)
	},
}
