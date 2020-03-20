package pullrequests

import (
	"io/ioutil"
	"os"

	"github.com/gfleury/gobbs/common"

	"gopkg.in/check.v1"
)

func (s *S) TestInfoInvalidHost(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"pr", "info", "21"}

	ctx := common.APIClientContext(&s.stashInfo)
	err := Info.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*connect: connection refused")
}

func (s *S) TestInfoInvalidHostTimeouted(c *check.C) {
	*s.host = "http://localhost:7992"

	os.Args = []string{"pr", "info", "21"}

	*s.stashInfo.Timeout() = 0

	ctx := common.APIClientContext(&s.stashInfo)
	err := Info.ExecuteContext(ctx)
	c.Assert(err, check.ErrorMatches, ".*context deadline exceeded")
	*s.stashInfo.Timeout() = 2
}

func (s *S) TestInfoValidHost(c *check.C) {
	*s.host = "http://localhost:7993"

	os.Args = []string{"pr", "info", "21"}

	s.mockStdout()

	ctx := common.APIClientContext(&s.stashInfo)
	err := Info.ExecuteContext(ctx)
	c.Assert(err, check.IsNil)

	s.closeMockStdout()

	got, err := ioutil.ReadAll(s.readStdout)
	c.Assert(err, check.IsNil)

	want, err := ioutil.ReadFile("mocked_responses/pullrequests_info.output")
	c.Assert(err, check.IsNil)

	c.Assert(string(got), check.Equals, string(want))
}
