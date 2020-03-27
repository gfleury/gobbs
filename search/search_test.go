package search

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gfleury/gobbs/common"

	bbmock "github.com/gfleury/go-bitbucket-v1/test/bb-mock-server/go"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	host, project, repo        *string
	stashInfo                  common.StashInfo
	originalStdout, readStdout *os.File
}

var _ = check.Suite(&S{})

func (s *S) TearDownSuite(c *check.C) {
}

func (s *S) SetUpSuite(c *check.C) {
	s.stashInfo = common.StashInfo{}
	s.host = s.stashInfo.Host()
	s.project = s.stashInfo.Project()
	s.repo = s.stashInfo.Repo()

	*s.stashInfo.Timeout() = 2

	*s.project = "PRJ1"
	*s.repo = "my-repository"

	*s.stashInfo.Credential().Passwd() = "passwd"
	*s.stashInfo.Credential().User() = "test"

	go func() {
		err := bbmock.RunServer(7995)
		if err != nil {
			c.Assert(err, check.IsNil)
		}
	}()

	for err := fmt.Errorf(""); err != nil; _, err = http.Get("http://localhost:7995/") {

		time.Sleep(time.Second)
	}

}

func (s *S) mockStdout() {
	// Create Fake Buffer for Stdout
	s.originalStdout = os.Stdout
	var w *os.File
	s.readStdout, w, _ = os.Pipe()
	os.Stdout = w
}

func (s *S) closeMockStdout() {
	os.Stdout.Close()
	os.Stdout = s.originalStdout
}

func (s *S) TestListInvalidHost(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"users", "list"}

	ctx := common.APIClientContext(&s.stashInfo)
	err := Search.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*connect: connection refused")
}

func (s *S) TestListInvalidHostTimeouted(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"users", "list"}

	*s.stashInfo.Timeout() = 0

	ctx := common.APIClientContext(&s.stashInfo)
	err := Search.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*context deadline exceeded")
	*s.stashInfo.Timeout() = 2
}

func (s *S) TestListValidHost(c *check.C) {
	*s.host = "http://localhost:7995"

	os.Args = []string{"search", "-N", "git", "clone"}

	s.mockStdout()

	ctx := common.APIClientContext(&s.stashInfo)
	err := Search.ExecuteContext(ctx)
	c.Assert(err, check.IsNil)

	s.closeMockStdout()

	got, err := ioutil.ReadAll(s.readStdout)
	c.Assert(err, check.IsNil)

	want, err := ioutil.ReadFile("mocked_responses/search.output")
	c.Assert(err, check.IsNil)

	c.Assert(string(got), check.Equals, string(want))
}
