package main

import (
	"fmt"
	"log"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/pullrequests"
	"github.com/mitchellh/go-homedir"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

func main() {

	var stashInfo common.StashInfo

	cobra.OnInitialize(initConfig)

	ctx, cancel := common.APIClientContext(&stashInfo)
	defer cancel()

	var rootCmd = &cobra.Command{Use: "gobbs"}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gobbs.yaml)")

	rootCmd.PersistentFlags().StringVarP(stashInfo.Host(), "host", "H", *stashInfo.Host(), "Stash host (https://stash.example.com)")
	rootCmd.PersistentFlags().StringVarP(stashInfo.Project(), "project", "P", *stashInfo.Project(), "Stash project slug (PRJ)")
	rootCmd.PersistentFlags().StringVarP(stashInfo.Repo(), "repository", "r", *stashInfo.Repo(), "Stash repository slug (repo01)")

	rootCmd.PersistentFlags().StringVarP(stashInfo.Credential().User(), "user", "u", *stashInfo.Credential().User(), "Stash username")
	rootCmd.PersistentFlags().StringVarP(stashInfo.Credential().Passwd(), "passwd", "p", *stashInfo.Credential().Passwd(), "Stash username password")

	pullrequests.PullRequestRoot.AddCommand(pullrequests.List)

	rootCmd.AddCommand(pullrequests.PullRequestRoot)

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = common.Config().BindPFlag(fmt.Sprintf(common.UserNameKey, *stashInfo.Host()), rootCmd.PersistentFlags().Lookup("user"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	if common.Config().GetString("password_method") == "" {
		err = common.Config().BindPFlag(fmt.Sprintf(common.PasswdKey, *stashInfo.Host()), rootCmd.PersistentFlags().Lookup("passwd"))
		if err != nil {
			log.Fatalln(err.Error())
		}
	} else if common.Config().GetString("password_method") == "gopass" {
		err = common.SavePasswdExternal(*stashInfo.Host(), rootCmd.PersistentFlags().Lookup("passwd").Value.String())
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	err = common.Config().WriteConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func initConfig() {
	common.SetConfig(viper.NewWithOptions(viper.KeyDelimiter("::")))

	if cfgFile != "" {
		// Use config file from the flag.
		common.Config().SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalln(err.Error())
		}

		// Search config in home directory with name ".gobbs" (without extension).
		common.Config().AddConfigPath(home)
		common.Config().SetConfigName(".gobbs")
	}

	common.Config().AutomaticEnv()

	if err := common.Config().ReadInConfig(); err == nil {
		fmt.Println("Using config file:", common.Config().ConfigFileUsed())
	}
}
