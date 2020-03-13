package common

import (
	"fmt"
	"log"

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

func (c *Credential) GetUser() string {
	if c.user == "" {
		prompt := fmt.Sprintf("Stash Username: ")
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
		prompt := fmt.Sprintf("Stash Password [%s]: ", c.user)
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
