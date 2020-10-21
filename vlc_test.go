package vlc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	testHandler := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "GET")
		assert.NotNil(t, req.Header.Get("Authorization"))

		status := `
		<root>
		  <state>stopped</state>
		</root>`

		w.Header().Set("Content-Type", "text/xml")

		fmt.Fprint(w, status)
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()

	vlc := New(testServer.URL, "1234")

	status, err := vlc.Status()
	if err != nil {
		t.Fatalf("%v", err)
	}

	assert.Equal(t, "stopped", status.State)
}
