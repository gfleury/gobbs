package common

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

var (
	isSchemeRegExp   = regexp.MustCompile(`^[^:]+://`)
	scpLikeURLRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5})(?:\/|:))?(?P<path>[^\\].*\/[^\\].*)$`)
)

// GitInfo try to probe stash/bitbucket information from git repository
func GitInfo() (host, project, repo string, err error) {
	gitRepo, err := openRepo()
	if err != nil {
		err = fmt.Errorf("Unable to open git repository: %s", err.Error())
		return
	}
	remotes, err := gitRepo.Remotes()
	if err != nil {
		err = fmt.Errorf("Unable to get remote references from git repository: %s", err.Error())
		return
	}

	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			return parseGitURL(remote.Config().URLs[0])
		}
	}

	if len(remotes) > 0 {
		return parseGitURL(remotes[0].Config().URLs[0])
	}

	return
}

func parseGitURL(rawURL string) (host, project, repo string, err error) {
	if !matchesScheme(rawURL) && matchesScpLike(rawURL) {
		_, host, _, path := findScpLikeComponents(rawURL)
		idx1, idx2 := 1, 2
		if !strings.HasPrefix(path, "/") {
			idx1, idx2 = 0, 1
		}
		urlPath := strings.Split(path, "/")
		if len(urlPath) > 1 {
			return host, urlPath[idx1], strings.Split(urlPath[idx2], ".")[0], err
		}
	}

	gitURL, err := url.Parse(rawURL)
	if err != nil {
		err = fmt.Errorf("Unable to parse remote git URL from git repository: %s", err.Error())
		return
	}
	urlPath := strings.Split(gitURL.Path, "/")
	if len(urlPath) > 1 {
		return gitURL.Hostname(), urlPath[1], strings.Split(urlPath[2], ".")[0], err
	}
	return
}

// From go-git internals
func matchesScpLike(url string) bool {
	is := scpLikeURLRegExp.MatchString(url)
	return is
}

// From go-git internals
func findScpLikeComponents(url string) (user, host, port, path string) {
	m := scpLikeURLRegExp.FindStringSubmatch(url)
	return m[1], m[2], m[3], m[4]
}

// From go-git internals
func matchesScheme(url string) bool {
	is := isSchemeRegExp.MatchString(url)
	return is
}

func openRepo() (*git.Repository, error) {
	cwd, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("Unable to get current directory: %s", err.Error())
		return nil, err
	}
	return git.PlainOpen(cwd)
}
