package common

type contextKey string

var StashInfoKey = contextKey("StashInfoKey")

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
