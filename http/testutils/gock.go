package httptestutils

import (
	http "net/http"

	gomock "github.com/golang/mock/gomock"
	"gopkg.in/h2non/gock.v1"
)

type GockMatcher struct {
	gock gock.Mock
}

func NewGockRequest() *gock.Request {
	req := gock.NewRequest()
	req.Response = gock.NewResponse()
	return req
}

func NewGockMatcher(req *gock.Request) *GockMatcher {
	if req.Response == nil {
		req.Response = gock.NewResponse()
	}

	return &GockMatcher{
		gock: gock.NewMock(req, req.Response),
	}
}

// Matches returns whether x is a match.
func (m *GockMatcher) Matches(x interface{}) bool {
	req, ok := x.(*http.Request)
	if !ok {
		return false
	}

	match, err := m.gock.Match(req)
	if err != nil {
		return false
	}

	return match
}

// String describes what the matcher matches.
func (m *GockMatcher) String() string {
	return "HTTP request matching right method"
}

// Gock indicates an expected call of Do.
func (mr *MockSenderMockRecorder) Gock(req *gock.Request) *gomock.Call {
	return mr.Do(NewGockMatcher(req)).DoAndReturn(func(r *http.Request) (*http.Response, error) { return gock.Responder(r, req.Response, nil) })
}
