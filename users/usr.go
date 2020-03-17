package users

import "github.com/spf13/cobra"

func init() {
	UserRoot.AddCommand(List)
	//UserRoot.AddCommand(Create)
}

// UserRoot cmd root for cobra
var UserRoot = &cobra.Command{
	Use:     "user",
	Aliases: []string{"us"},
	Short:   "Interact with users",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			List.Run(cmd, args)
		}
	},
}
