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
	err := Config().ReadConfig(bytes.NewBuffer(yamlExample))
	c.Assert(err, check.IsNil)

	a := cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("TestAPIClient")
		},
	}
	ctx := APIClientContext(&StashInfo{})

	err = a.ExecuteContext(ctx)
	c.Assert(err, check.IsNil)

	err = os.Chdir("..")
	c.Assert(err, check.IsNil)
	api, cancel, err := APIClient(&a)
	defer cancel()
	c.Assert(err, check.IsNil)
	err = os.Chdir("common")
	c.Assert(err, check.IsNil)
	c.Assert(api, check.NotNil)
	os.Args = savedArgs
}

func (s *S) TestAPIClientWithHost(c *check.C) {
	os.Args = []string{"help"}
	Config().SetConfigType("yaml")
	var yamlExample = []byte(`
https://localhost:
  user: gfleury
  passwd: 123
`)
	err := Config().ReadConfig(bytes.NewBuffer(yamlExample))
	c.Assert(err, check.IsNil)

	a := cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	si := StashInfo{
		host: "localhost",
	}
	ctx := APIClientContext(&si)

	err = a.ExecuteContext(ctx)
	c.Assert(err, check.IsNil)

	err = os.Chdir("..")
	c.Assert(err, check.IsNil)
	_, cancel, err := APIClient(&a)
	defer cancel()
	c.Assert(err, check.IsNil)

	c.Assert(si.host, check.Equals, "https://localhost")
	err = os.Chdir("common")
	c.Assert(err, check.IsNil)
}

func (s *S) TestAPIClientWithArguments(c *check.C) {
	os.Args = []string{"help"}
	Config().SetConfigType("yaml")
	var yamlExample = []byte(`
http://localhost:
  user: gfleury
  passwd: 123
`)
	err := Config().ReadConfig(bytes.NewBuffer(yamlExample))
	c.Assert(err, check.IsNil)

	a := cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	si := StashInfo{
		host: "http://localhost",
	}
	ctx := APIClientContext(&si)

	err = a.ExecuteContext(ctx)
	c.Assert(err, check.IsNil)

	err = os.Chdir("..")
	c.Assert(err, check.IsNil)
	_, cancel, err := APIClient(&a)
	defer cancel()
	c.Assert(err, check.IsNil)

	c.Assert(si.host, check.Equals, "http://localhost")
	err = os.Chdir("common")
	c.Assert(err, check.IsNil)
}
