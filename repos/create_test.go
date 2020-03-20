package repos

import (
	"io/ioutil"
	"os"

	"github.com/gfleury/gobbs/common"

	"gopkg.in/check.v1"
)

func (s *S) TestCreateInvalidHost(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"repo", "create", "newRepo"}

	ctx := common.APIClientContext(&s.stashInfo)
	err := Create.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*connect: connection refused")
}

func (s *S) TestCreateInvalidHostTimeouted(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"repo", "create", "newRepo"}

	*s.stashInfo.Timeout() = 0

	ctx := common.APIClientContext(&s.stashInfo)
	err := Create.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*context deadline exceeded")
	*s.stashInfo.Timeout() = 2
}

func (s *S) TestCreateValidHost(c *check.C) {
	*s.host = "http://localhost:7991"

	os.Args = []string{"repo", "create", "newRepo"}

	s.mockStdout()

	ctx := common.APIClientContext(&s.stashInfo)
	err := Create.ExecuteContext(ctx)
	c.Assert(err, check.IsNil)

	s.closeMockStdout()

	got, err := ioutil.ReadAll(s.readStdout)
	c.Assert(err, check.IsNil)

	want, err := ioutil.ReadFile("mocked_responses/repos_create.output")
	c.Assert(err, check.IsNil)

	c.Assert(string(got), check.Equals, string(want))
}
