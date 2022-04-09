package types

import (
	"fmt"
)

type Error struct {
	Code        int      `json:"code"`        // either a specific error code in case of invalid request or http status code
	Message     string   `json:"message"`     // message describing error
	StackTraces []string `json:"stacktraces"` // optional stacktraces, sent when node is in debug mode
}

func (err Error) Error() string {
	s := fmt.Sprintf("BeaconError: message=%q code=%d", err.Message, err.Code)
	if len(err.StackTraces) != 0 {
		s = fmt.Sprintf("%v\nstacktraces=%q", s, err.StackTraces)
	}

	return s
}
