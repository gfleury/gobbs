package main

import (
	"log"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/pullrequests"

	"github.com/spf13/cobra"
)

func main() {

	var stashInfo common.StashInfo

	ctx, cancel := common.APIClientContext(&stashInfo)
	defer cancel()

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.PersistentFlags().StringVar(stashInfo.Host(), "host", *stashInfo.Host(), "Stash host (https://stash.example.com)")
	rootCmd.PersistentFlags().StringVar(stashInfo.Project(), "project", *stashInfo.Project(), "Stash project slug (PRJ)")
	rootCmd.PersistentFlags().StringVar(stashInfo.Repo(), "repository", *stashInfo.Repo(), "Stash repository slug (repo01)")

	rootCmd.PersistentFlags().StringVar(stashInfo.Credential().User(), "user", *stashInfo.Credential().User(), "Stash username")
	rootCmd.PersistentFlags().StringVar(stashInfo.Credential().Passwd(), "passwd", *stashInfo.Credential().Passwd(), "Stash username password")

	pullrequests.PullRequestRoot.AddCommand(pullrequests.List)

	rootCmd.AddCommand(pullrequests.PullRequestRoot)

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}
