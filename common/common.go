package common

import (
	"github.com/gfleury/gobbs/common/log"

	"github.com/spf13/viper"
)

type contextKey string

const (
	StashInfoKey = contextKey("StashInfoKey")

	UserNameKey = "%s::user"
	PasswdKey   = "%s::passwd"

	AppName = "gobbs"
)

var v *viper.Viper

type StashInfo struct {
	host, project, repo string
	credential          Credential
	timeout             int
}

// Host return host from stashInfo struct
func (s *StashInfo) Host() *string {
	return &s.host
}

// Project return project from stashInfo struct
func (s *StashInfo) Project() *string {
	return &s.project
}

// Repo return repository from stashInfo struct
func (s *StashInfo) Repo() *string {
	return &s.repo
}

// Credential return credential from stashInfo struct
func (s *StashInfo) Credential() *Credential {
	return &s.credential
}

// Timeout return timeout from stashInfo struct
func (s *StashInfo) Timeout() *int {
	return &s.timeout
}

func Config() *viper.Viper {
	return v
}

func SetConfig(c *viper.Viper) {
	v = c
}

func PrintApiError(values map[string]interface{}) {
	if _, ok := values["errors"]; ok {
		if errors, ok := values["errors"].([]interface{}); ok {
			if len(errors) > 0 {
				if msgs, ok := errors[0].(map[string]interface{}); ok {
					for idx := range msgs {
						if _, ok := msgs[idx].(string); ok {
							log.Critical("%s: %s", idx, msgs[idx])
						}
						if _, ok := msgs[idx].(float64); ok {
							log.Critical("%s: %v", idx, msgs[idx])
						}
					}
				}
			}
		}
	}
}
