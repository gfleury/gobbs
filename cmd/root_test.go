package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gfleury/gobbs/common"
	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
}

var _ = check.Suite(&S{})

func (s *S) TestInitConfig(c *check.C) {
	cfgFile = fmt.Sprintf("%s/config.yaml", os.TempDir())
	defer os.Remove(cfgFile)
	// Run once to create file
	initConfig()

	// Run twice with file already created
	initConfig()
}

func (s *S) TestPersistConfig(c *check.C) {
	cfgFile = fmt.Sprintf("%s/config.yaml", os.TempDir())
	defer os.Remove(cfgFile)

	host := stashInfo.Host()
	user := stashInfo.Credential().User()
	passwd := stashInfo.Credential().Passwd()
	stashInfo.Credential().New()

	*host = "stash.example.com"
	*user = "user"
	*passwd = "password"
	// Run once to create file
	initConfig()

	err := persistConfig()
	c.Assert(err, check.IsNil)

	content, err := ioutil.ReadFile(cfgFile)
	c.Assert(err, check.IsNil)
	c.Assert(string(content), check.Equals, ""+"stash.example.com:\n"+"  passwd: password\n"+"  user: user\n")
}

func (s *S) TestGetContext(c *check.C) {
	ctx := common.APIClientContext(&stashInfo)

	c.Check(ctx.Value(common.StashInfoKey), check.Equals, &stashInfo)
}

func (s *S) TestConfigWithEnvVars(c *check.C) {
	os.Setenv("GOBBS_USER", "user")
	initConfig()
	c.Assert(common.Config().Get("user"), check.Equals, "user")
}
