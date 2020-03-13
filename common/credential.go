package common

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

type Credential struct {
	user, passwd string
}

func (c *Credential) User() *string {
	return &c.user
}

func (c *Credential) Passwd() *string {
	return &c.passwd
}

func (c *Credential) GetUser(host string) string {
	if c.user == "" {
		cachedCredentials := Config().GetStringMapString(host)
		c.user = cachedCredentials["user"]
		if Config().GetString("password_method") == "gopass" {
			// Got from go-jira
			log.Printf("Querying gopass password source.")
			binary := "gopass"

			if passDir := Config().GetString("password_storage"); passDir != "" {
				orig := os.Getenv("PASSWORD_STORE_DIR")
				log.Printf("using password-directory: %s", passDir)
				os.Setenv("PASSWORD_STORE_DIR", passDir)
				defer os.Setenv("PASSWORD_STORE_DIR", orig)
			}
			if passDir := os.Getenv("PASSWORD_STORE_DIR"); passDir != "" {
				log.Printf("using PASSWORD_STORE_DIR=%s", passDir)
			}
			if bin, err := exec.LookPath(binary); err == nil {
				log.Printf("found gopass at: %s", bin)
				buf := bytes.NewBufferString("")
				cmd := exec.Command(bin, "show", "-o", fmt.Sprintf("Gobbs/%s", strings.Replace(host, "://", "-", -1)))
				cmd.Stdout = buf
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err == nil {
					c.passwd = strings.TrimSpace(buf.String())
				} else {
					log.Printf("gopass command failed with:\n%s", buf.String())
				}
			} else {
				log.Printf("Gopass binary was not found! Fallback to default password behaviour!")
			}
		} else {
			c.passwd = cachedCredentials["passwd"]
		}
	}
	if c.user == "" {
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
		log.Printf("Saving gopass password source.")
		binary := "gopass"

		if passDir := Config().GetString("password_storage"); passDir != "" {
			orig := os.Getenv("PASSWORD_STORE_DIR")
			log.Printf("using password-directory: %s", passDir)
			os.Setenv("PASSWORD_STORE_DIR", passDir)
			defer os.Setenv("PASSWORD_STORE_DIR", orig)
		}
		if passDir := os.Getenv("PASSWORD_STORE_DIR"); passDir != "" {
			log.Printf("using PASSWORD_STORE_DIR=%s", passDir)
		}
		if bin, err := exec.LookPath(binary); err == nil {
			log.Printf("found gopass at: %s", bin)
			buf := bytes.NewBufferString("")
			stdin := bytes.NewBufferString(passwd)
			cmd := exec.Command(bin, "insert", "-f", fmt.Sprintf("Gobbs/%s", strings.Replace(host, "://", "-", -1)))
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
