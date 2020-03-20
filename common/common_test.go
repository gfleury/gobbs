package common

import "gopkg.in/check.v1"

func (s *S) TestPrintApiError(c *check.C) {
	values := map[string]interface{}{
		"errors": map[string]interface{}{
			"context":       nil,
			"message":       "You are not permitted to access this resource",
			"exceptionName": nil,
		},
	}

	PrintApiError(values)
}
