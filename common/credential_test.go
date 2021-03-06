package common

import (
	"bytes"

	"github.com/spf13/viper"
	"gopkg.in/check.v1"
)

func (s *S) TestSavePasswdExternal(c *check.C) {
	SetConfig(viper.NewWithOptions(viper.KeyDelimiter("::")))
	Config().SetConfigType("yaml")
	var yamlExample = []byte(`
password_method: gopass
`)
	err := Config().ReadConfig(bytes.NewBuffer(yamlExample))
	c.Assert(err, check.IsNil)

	binary = "/bin/echo"
	err = SavePasswdExternal("stash.example.com", "password")
	c.Assert(err, check.IsNil)
}

func (s *S) TestGetUser(c *check.C) {
	SetConfig(viper.NewWithOptions(viper.KeyDelimiter("::")))

	Config().SetConfigType("yaml")
	var yamlExample = []byte(`
localhost:
  user: test
  passwd: passwd
`)

	err := Config().ReadConfig(bytes.NewBuffer(yamlExample))
	c.Assert(err, check.IsNil)

	cred := Credential{}
	user := cred.GetUser("localhost")
	c.Assert(user, check.Equals, "test")

	passwd := cred.GetPasswd()
	c.Assert(passwd, check.Equals, "passwd")
}

func (s *S) TestGetUserGoPass(c *check.C) {
	SetConfig(viper.NewWithOptions(viper.KeyDelimiter("::")))

	Config().SetConfigType("yaml")
	var yamlExample = []byte(`
password_method: gopass
localhost:
  user: test
  passwd: passwd
`)

	err := Config().ReadConfig(bytes.NewBuffer(yamlExample))
	c.Assert(err, check.IsNil)

	binary = "/bin/echo"

	cred := Credential{}
	user := cred.GetUser("localhost")
	c.Assert(user, check.Equals, "test")

	passwd := cred.GetPasswd()
	c.Assert(passwd, check.Equals, "show -o gobbs/localhost")
}

func (s *S) TestNormalizeGoPassKey(c *check.C) {
	key := normalizeGoPassKey("https://www.google.com.br:443/")
	c.Assert(key, check.Equals, "https-www.google.com.br#443/")
}
