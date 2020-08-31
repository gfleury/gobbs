package cmd

import (
	"fmt"
	"os"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"
	"github.com/gfleury/gobbs/pullrequests"
	"github.com/gfleury/gobbs/repos"
	"github.com/gfleury/gobbs/search"
	"github.com/gfleury/gobbs/users"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	verbosity      int
	stashInfo      common.StashInfo
	storePWDGoPass bool
)
var rootCmd = &cobra.Command{Use: common.AppName}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gobbs.yaml)")
	rootCmd.PersistentFlags().IntVarP(stashInfo.Timeout(), "timeout", "t", 30, "timeout for api requests")
	rootCmd.PersistentFlags().IntVarP(&verbosity, "verbose", "v", 1, "Increase verbosity for debugging")

	rootCmd.PersistentFlags().StringVarP(stashInfo.Host(), "host", "H", *stashInfo.Host(), "Stash host (https://stash.example.com)")
	rootCmd.PersistentFlags().StringVarP(stashInfo.Project(), "project", "P", *stashInfo.Project(), "Stash project slug (PRJ)")
	rootCmd.PersistentFlags().StringVarP(stashInfo.Repo(), "repository", "r", *stashInfo.Repo(), "Stash repository slug (repo01)")

	rootCmd.PersistentFlags().BoolVarP(&storePWDGoPass, "gopass", "g", false, "Enable password to be stored with gopass (gopass must be installed, https://www.gopass.pw/)")
	rootCmd.PersistentFlags().StringVarP(stashInfo.Credential().User(), "user", "u", *stashInfo.Credential().User(), "Stash username")
	rootCmd.PersistentFlags().StringVarP(stashInfo.Credential().Passwd(), "passwd", "p", *stashInfo.Credential().Passwd(), "Stash username password")

	rootCmd.AddCommand(pullrequests.PullRequestRoot)
	rootCmd.AddCommand(repos.ReposRoot)
	rootCmd.AddCommand(users.UserRoot)
	rootCmd.AddCommand(search.Search)
}

func initConfig() {
	os.Setenv("GOBBS_DEBUG", fmt.Sprintf("%d", verbosity))
	log.IncreaseLogLevel(verbosity)

	common.SetConfig(viper.NewWithOptions(viper.KeyDelimiter("::")))

	common.Config().SetEnvPrefix("GOBBS")
	common.Config().AutomaticEnv()

	if cfgFile == "" {
		cfgFile = common.Config().GetString("config")
	}

	if *stashInfo.Host() == "" {
		stashInfo.SetHost(common.Config().GetString("host"))
	}
	if *stashInfo.Credential().User() == "" {
		stashInfo.Credential().SetUser(common.Config().GetString("user"))
	}
	if *stashInfo.Credential().Passwd() == "" {
		stashInfo.Credential().SetPasswd(common.Config().GetString("passwd"))
	}

	if cfgFile != "" {
		// Use config file from the flag.
		common.Config().SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err.Error())
		}

		// Search config in home directory with name ".gobbs" (without extension).
		common.Config().AddConfigPath(home)
		common.Config().SetConfigName(fmt.Sprintf(".%s", common.AppName))
		cfgFile = fmt.Sprintf("%s/.%s.yaml", home, common.AppName)
	}

	if err := common.Config().ReadInConfig(); err == nil {
		log.Debugf("Using config file:", common.Config().ConfigFileUsed())
	} else {
		log.Debugf("Config file not found, creating empty file:", cfgFile)
		file, err := os.Create(cfgFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
	}
}

func persistConfig() error {
	// common.Config() is only initialized if rootCmd ran (false if only help mode was shown)
	if common.Config() != nil {
		err := common.Config().BindPFlag(fmt.Sprintf(common.UserNameKey, *stashInfo.Host()), rootCmd.PersistentFlags().Lookup("user"))
		if err != nil {
			return err
		}
		if common.Config().GetString("password_method") == "" && !storePWDGoPass && stashInfo.Credential().IsNew() {
			err = common.Config().BindPFlag(fmt.Sprintf(common.PasswdKey, *stashInfo.Host()), rootCmd.PersistentFlags().Lookup("passwd"))
			if err != nil {
				return err
			}
		} else if (storePWDGoPass || common.Config().GetString("password_method") == "gopass") && stashInfo.Credential().IsNew() {
			err = common.SavePasswdExternal(*stashInfo.Host(), rootCmd.PersistentFlags().Lookup("passwd").Value.String())
			if err != nil {
				return err
			}
		}
		if storePWDGoPass {
			common.Config().Set("password_method", "gopass")
		}
		err = common.Config().WriteConfig()
		return err
	}
	return nil
}

func Execute() error {
	ctx := common.APIClientContext(&stashInfo)

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		return err
	}

	return persistConfig()
}
