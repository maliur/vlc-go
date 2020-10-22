package vlc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

	vlcClient := NewClient(testServer.URL, "1234")

	status, err := vlcClient.Status()
	if err != nil {
		t.Fatalf("%v", err)
	}

	assert.Equal(t, "stopped", status.State)
}

func TestPlaylist(t *testing.T) {
	testHandler := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "GET")
		assert.NotNil(t, req.Header.Get("Authorization"))

		status := `
		<node ro="rw" name="" id="0">
  			<node ro="ro" name="Playlist" id="1">
    			<leaf ro="rw" name="Foo song" id="4" duration="226" uri="foo.mp3" current="current"/>
    			<leaf ro="rw" name="Bar song" id="5" duration="6" uri="bar.mp3"/>
			</node>
		</node>`

		w.Header().Set("Content-Type", "text/xml")

		fmt.Fprint(w, status)
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()

	vlcClient := NewClient(testServer.URL, "1234")

	playlist, err := vlcClient.Playlist()
	if err != nil {
		t.Fatalf("%v", err)
	}

	assert.Equal(t, 2, len(playlist.Songs))
}

func TestAddSong(t *testing.T) {
	songUrl := "https://youtu.be/dQw4w9WgXcQ"

	testHandler := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "GET")
		assert.NotNil(t, req.Header.Get("Authorization"))
		assert.True(t, strings.Contains(req.URL.RawQuery, songUrl))
		assert.True(t, strings.Contains(req.URL.RawQuery, "in_enqueue"))

		fmt.Fprintf(w, "OK")
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()

	vlcClient := NewClient(testServer.URL, "1234")

	err := vlcClient.AddSong(songUrl, false)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestNextSong(t *testing.T) {
	testHandler := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "GET")
		assert.NotNil(t, req.Header.Get("Authorization"))
		assert.True(t, strings.Contains(req.URL.RawQuery, "command=pl_next"))

		fmt.Fprintf(w, "OK")
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()

	vlcClient := NewClient(testServer.URL, "1234")

	err := vlcClient.NextSong()
	if err != nil {
		t.Fatalf("%v", err)
	}
}
