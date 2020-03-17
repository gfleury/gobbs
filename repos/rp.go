package repos

import "github.com/spf13/cobra"

func init() {
	ReposRoot.AddCommand(List)
	//ReposRoot.AddCommand(Create)
}

// ReposRoot cmd root for cobra
var ReposRoot = &cobra.Command{
	Use:     "repo",
	Aliases: []string{"rp"},
	Short:   "Interact with repositories",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			List.Run(cmd, args)
		}
	},
}
