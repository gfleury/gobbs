package users

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return List.RunE(cmd, args)
		}
		return fmt.Errorf("Commnand not found")
	},
}
