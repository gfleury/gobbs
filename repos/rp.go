package repos

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	ReposRoot.AddCommand(List)
	ReposRoot.AddCommand(Create)
}

// ReposRoot cmd root for cobra
var ReposRoot = &cobra.Command{
	Use:     "repo",
	Aliases: []string{"rp"},
	Short:   "Interact with repositories",
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return List.RunE(cmd, args)
		}
		return fmt.Errorf("Commnand not found")
	},
}
