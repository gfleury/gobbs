package pullrequests

import "github.com/spf13/cobra"

func init() {
	PullRequestRoot.AddCommand(List)
	PullRequestRoot.AddCommand(Create)
	PullRequestRoot.AddCommand(Delete)
}

// PullRequestRoot cmd root for cobra
var PullRequestRoot = &cobra.Command{
	Use:     "pullrequest",
	Aliases: []string{"pr"},
	Short:   "Interact with pull requests",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			List.Run(cmd, args)
		}
	},
}
