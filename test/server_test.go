package test

import (
	"net/http/httptest"
	"testing"

	. "github.com/poy/onpar/expect"
	. "github.com/poy/onpar/matchers"
)

func TestServer(t *testing.T) {
	ts := httptest.NewServer(GetAPI())
	defer ts.Close()

	_, status, err := SendRequest("GET", ts.URL+"/testRoute?isError=false", nil, nil)
	Expect(t, err).To(Not(HaveOccurred()))

	Expect(t, status).To(Equal(200))

	_, status, err = SendRequest("GET", ts.URL+"/testRoute?isError=true", nil, nil)
	Expect(t, err).To(Not(HaveOccurred()))

	Expect(t, status).To(Equal(409))
}
