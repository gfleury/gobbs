package common

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
)

func (s *S) TestAPIClient(c *check.C) {
	savedArgs := os.Args
	os.Args = []string{"help"}

	SetConfig(viper.NewWithOptions(viper.KeyDelimiter("::")))
	Config().SetConfigType("yaml")
	var yamlExample = []byte(`
https://github.com:
  user: gfleury
  passwd: 123
`)
	Config().ReadConfig(bytes.NewBuffer(yamlExample))

	a := cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TestAPIClient")
		},
	}
	ctx := APIClientContext(&StashInfo{})

	err := a.ExecuteContext(ctx)
	c.Assert(err, check.IsNil)

	err = os.Chdir("..")
	api, cancel, err := APIClient(&a)
	defer cancel()
	c.Assert(err, check.IsNil)
	err = os.Chdir("common")
	c.Assert(api, check.NotNil)
	os.Args = savedArgs
}
