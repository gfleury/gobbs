package pullrequests

import (
	"fmt"

	"github.com/gfleury/gobbs/common"

	"github.com/spf13/cobra"
)

func init() {
	PullRequestRoot.AddCommand(List)
	PullRequestRoot.AddCommand(Create)
	PullRequestRoot.AddCommand(Delete)
	PullRequestRoot.AddCommand(Approve)
	PullRequestRoot.AddCommand(Merge)
	PullRequestRoot.AddCommand(Info)
	PullRequestRoot.AddCommand(Diff)
}

// PullRequestRoot cmd root for cobra
var PullRequestRoot = &cobra.Command{
	Use:     "pullrequest",
	Aliases: []string{"pr"},
	Short:   "Interact with pull requests",
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return List.RunE(cmd, args)
		}
		return fmt.Errorf("Commnand not found")
	},
}

func mustHaveProjectRepo(stashInfo *common.StashInfo) error {
	if *stashInfo.Repo() == "" ||
		*stashInfo.Project() == "" {
		return fmt.Errorf("Unable to identify Project and Repository")
	}
	return nil
}
