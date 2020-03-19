package pullrequests

import (
	"io/ioutil"
	"os"

	"github.com/gfleury/gobbs/common"

	"gopkg.in/check.v1"
)

func (s *S) TestCreateInvalidHost(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"pr", "create", "featureBranch", "master"}

	ctx := common.APIClientContext(&s.stashInfo)
	err := Create.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*connect: connection refused")
}

func (s *S) TestCreateInvalidHostTimeouted(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"pr", "create", "featureBranch", "master"}

	*s.stashInfo.Timeout() = 0

	ctx := common.APIClientContext(&s.stashInfo)
	err := Create.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*context deadline exceeded")
	*s.stashInfo.Timeout() = 2
}

func (s *S) TestCreateValidHost(c *check.C) {
	*s.host = "http://localhost:7993"

	os.Args = []string{"pr", "create", "featureBranch", "master"}

	s.mockStdout()

	ctx := common.APIClientContext(&s.stashInfo)
	err := Create.ExecuteContext(ctx)
	c.Assert(err, check.IsNil)

	s.closeMockStdout()

	got, err := ioutil.ReadAll(s.readStdout)
	c.Assert(err, check.IsNil)

	want, err := ioutil.ReadFile("mocked_responses/pullrequests_create.output")
	c.Assert(err, check.IsNil)

	c.Assert(string(got), check.Equals, string(want))
}
