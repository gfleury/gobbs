package pullrequests

import (
	"io/ioutil"
	"os"

	"github.com/gfleury/gobbs/common"

	"gopkg.in/check.v1"
)

func (s *S) TestDiffInvalidHost(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"pr", "diff", "21"}

	ctx := common.APIClientContext(&s.stashInfo)
	err := Diff.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*connect: connection refused")
}

func (s *S) TestDiffInvalidHostTimeouted(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"pr", "diff", "21"}

	*s.stashInfo.Timeout() = 0

	ctx := common.APIClientContext(&s.stashInfo)
	err := Diff.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*context deadline exceeded")
	*s.stashInfo.Timeout() = 2
}

func (s *S) TestDiffValidHost(c *check.C) {
	*s.host = "http://localhost:7993"

	os.Args = []string{"pr", "diff", "21", "-N"}

	s.mockStdout()

	ctx := common.APIClientContext(&s.stashInfo)
	err := Diff.ExecuteContext(ctx)
	c.Assert(err, check.IsNil)

	s.closeMockStdout()

	got, err := ioutil.ReadAll(s.readStdout)
	c.Assert(err, check.IsNil)

	want, err := ioutil.ReadFile("mocked_responses/pullrequests_diff.output")
	c.Assert(err, check.IsNil)

	c.Assert(string(got), check.Equals, string(want))
}
