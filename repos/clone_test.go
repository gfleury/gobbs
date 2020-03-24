package repos

import (
	"io/ioutil"
	"os"

	"github.com/gfleury/gobbs/common"

	"gopkg.in/check.v1"
)

func (s *S) TestCloneInvalidHost(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"repos", "clone", "newRepo"}

	ctx := common.APIClientContext(&s.stashInfo)
	err := Clone.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*connect: connection refused")
}

func (s *S) TestCloneInvalidHostTimeouted(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"repos", "clone", "newRepo"}

	*s.stashInfo.Timeout() = 0

	ctx := common.APIClientContext(&s.stashInfo)
	err := Clone.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*context deadline exceeded")
	*s.stashInfo.Timeout() = 2
}

func (s *S) TestCloneValidHost(c *check.C) {
	*s.host = "http://localhost:7991"

	os.Args = []string{"repos", "clone", "newRepo"}

	s.mockStdout()

	ctx := common.APIClientContext(&s.stashInfo)
	err := Clone.ExecuteContext(ctx)
	c.Assert(err, check.NotNil)

	s.closeMockStdout()

	got, err := ioutil.ReadAll(s.readStdout)
	c.Assert(err, check.IsNil)

	want, err := ioutil.ReadFile("mocked_responses/repos_clone.output")
	c.Assert(err, check.IsNil)

	c.Assert(string(got), check.Equals, string(want))

}
