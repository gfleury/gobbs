package common

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gfleury/gobbs/common/log"

	"gopkg.in/AlecAivazis/survey.v1"
)

var binary = "gopass"

type Credential struct {
	user, passwd string
	new          bool
}

func (c *Credential) User() *string {
	return &c.user
}

func (c *Credential) Passwd() *string {
	return &c.passwd
}

func (c *Credential) IsNew() bool {
	return c.new
}

func (c *Credential) New() {
	c.new = true
}

func (c *Credential) GetUser(host string) string {
	if c.user == "" {
		c.new = false
		cachedCredentials := Config().GetStringMapString(host)
		c.user = cachedCredentials["user"]
		if Config().GetString("password_method") == "gopass" {
			// Got from go-jira
			log.Debugf("Querying gopass password source.")

			if passDir := Config().GetString("password_storage"); passDir != "" {
				orig := os.Getenv("PASSWORD_STORE_DIR")
				log.Debugf("using password-directory: %s", passDir)
				os.Setenv("PASSWORD_STORE_DIR", passDir)
				defer os.Setenv("PASSWORD_STORE_DIR", orig)
			}
			if passDir := os.Getenv("PASSWORD_STORE_DIR"); passDir != "" {
				log.Debugf("using PASSWORD_STORE_DIR=%s", passDir)
			}
			if bin, err := exec.LookPath(binary); err == nil {
				log.Debugf("found gopass at: %s", bin)
				buf := bytes.NewBufferString("")
				cmd := exec.Command(bin, "show", "-o", fmt.Sprintf("%s/%s", AppName, strings.Replace(host, "://", "-", -1)))
				cmd.Stdout = buf
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err == nil {
					c.passwd = strings.TrimSpace(buf.String())
				} else {
					log.Debugf("gopass command failed with:\n%s", buf.String())
				}
			} else {
				log.Debugf("Gopass binary was not found! Fallback to default password behaviour!")
			}
		} else {
			c.passwd = cachedCredentials["passwd"]
		}
	}
	if c.user == "" {
		c.new = true
		prompt := fmt.Sprintf("Stash Username [%s]:", host)
		help := ""
		err := survey.AskOne(
			&survey.Input{
				Message: prompt,
				Help:    help,
			},
			&c.user,
			nil,
		)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	}
	return c.user
}

func (c *Credential) GetPasswd() string {
	if c.passwd == "" {
		prompt := fmt.Sprintf("Stash Password [%s]:", c.user)
		help := ""
		err := survey.AskOne(
			&survey.Password{
				Message: prompt,
				Help:    help,
			},
			&c.passwd,
			nil,
		)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	}
	return c.passwd
}

func SavePasswdExternal(host, passwd string) error {
	if Config().GetString("password_method") == "gopass" {
		log.Debugf("Saving gopass password source.")

		if passDir := Config().GetString("password_storage"); passDir != "" {
			orig := os.Getenv("PASSWORD_STORE_DIR")
			log.Debugf("using password-directory: %s", passDir)
			os.Setenv("PASSWORD_STORE_DIR", passDir)
			defer os.Setenv("PASSWORD_STORE_DIR", orig)
		}
		if passDir := os.Getenv("PASSWORD_STORE_DIR"); passDir != "" {
			log.Debugf("using PASSWORD_STORE_DIR=%s", passDir)
		}
		if bin, err := exec.LookPath(binary); err == nil {
			log.Debugf("found gopass at: %s", bin)
			buf := bytes.NewBufferString("")
			stdin := bytes.NewBufferString(passwd)
			cmd := exec.Command(bin, "insert", "-f", fmt.Sprintf("%s/%s", AppName, strings.Replace(host, "://", "-", -1)))
			cmd.Stdin = stdin
			cmd.Stdout = buf
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Gopass binary was not found! Fallback to default password behaviour")
		}
	}
	return nil
}
