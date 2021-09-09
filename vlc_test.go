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
          <information>
              <category name="meta">
                  <info name='album'>F-Zero</info>
                  <info name='filename'>04 Fire Field.mp3</info>
                  <info name='copyright'>1990 Nintendo</info>
                  <info name='dumper'>Datschge</info>
                  <info name='artist'>Yumiko Kanki</info>
                  <info name='title'>Fire Field</info>
              </category>
              <category name='Stream 0'>
                  <info name='Bitrate'>320 kb/s</info>
                  <info name='Codec'>MPEG Audio layer 1/2 (mpga)</info>
                  <info name='Channels'>Stereo</info>
                  <info name='Bits per sample'>32</info>
                  <info name='Sample rate'>44100 Hz</info>
                  <info name='Type'>Audio</info>
              </category>
          </information>
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
	assert.Equal(t, "Yumiko Kanki", status.Artist())
	assert.Equal(t, "Fire Field", status.Song())
}

func TestPlaylist(t *testing.T) {
	testHandler := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "GET")
		assert.NotNil(t, req.Header.Get("Authorization"))

		playlist := `
		<node ro="rw" name="" id="0">
  			<node ro="ro" name="Playlist" id="1">
    			<leaf ro="rw" name="Foo song" id="4" duration="226" uri="foo.mp3" current="current"/>
    			<leaf ro="rw" name="Bar song" id="5" duration="6" uri="bar.mp3"/>
			</node>
		</node>`

		w.Header().Set("Content-Type", "text/xml")

		fmt.Fprint(w, playlist)
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
