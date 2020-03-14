package common

import "github.com/spf13/viper"

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

func Config() *viper.Viper {
	return v
}

func SetConfig(c *viper.Viper) {
	v = c
}
