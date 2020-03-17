package common

import (
	"os"
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
}

var _ = check.Suite(&S{})

func (s *S) TestGitInfo(c *check.C) {
	err := os.Chdir("..")
	c.Assert(err, check.IsNil)
	host, project, repo, err := GitInfo()
	c.Assert(err, check.IsNil)
	c.Assert(host, check.Equals, "github.com")
	c.Assert(project, check.Equals, "gfleury")
	c.Assert(repo, check.Equals, "gobbs")
	err = os.Chdir("common")
	c.Assert(err, check.IsNil)
}

func (s *S) TestGitInfoErr(c *check.C) {
	_, _, _, err := GitInfo()
	c.Assert(err, check.NotNil)
}

func (s *S) TestparseGitURL(c *check.C) {
	host, project, repo, err := parseGitURL("https://github.com/gfleury/gobbs.git")
	c.Assert(err, check.IsNil)
	c.Assert(host, check.Equals, "github.com")
	c.Assert(project, check.Equals, "gfleury")
	c.Assert(repo, check.Equals, "gobbs")
}

func (s *S) TestparseGitURLSSHURL(c *check.C) {
	host, project, repo, err := parseGitURL("ssh://git@stash.example.com:7999/project/repo01.git")
	c.Assert(err, check.IsNil)
	c.Assert(host, check.Equals, "stash.example.com")
	c.Assert(project, check.Equals, "project")
	c.Assert(repo, check.Equals, "repo01")
}

func (s *S) TestparseGitURLSCPSSHURL(c *check.C) {
	host, project, repo, err := parseGitURL("git@github.com:gfleury/gobbs.git")
	c.Assert(err, check.IsNil)
	c.Assert(host, check.Equals, "github.com")
	c.Assert(project, check.Equals, "gfleury")
	c.Assert(repo, check.Equals, "gobbs")
}

func (s *S) TestparseGitURLSCPSSHURL_2(c *check.C) {
	host, project, repo, err := parseGitURL("git@github.com:/gfleury/gobbs.git")
	c.Assert(err, check.IsNil)
	c.Assert(host, check.Equals, "github.com")
	c.Assert(project, check.Equals, "gfleury")
	c.Assert(repo, check.Equals, "gobbs")
}
