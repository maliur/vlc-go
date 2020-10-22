package vlc

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Status struct {
	State string
}

type Playlist struct {
	Songs []PlaylistSong
}

type PlaylistSong struct {
	ID       int
	Name     string
	Current  bool
	Duration int
	Uri      string
}

type Vlc interface {
	// Get current playlist
	Playlist() (Playlist, error)
	// Get the current status from vlc
	Status() (Status, error)
	// Enqueue song in playlist
	AddSong(url string, playNow bool) error
	// Play next song in playlist
	NextSong() error
}

type vlc struct {
	address  string
	password string
}

// Create a new Vlc client
func NewClient(address, password string) Vlc {
	return &vlc{
		address,
		password,
	}
}

func get(url, password string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("", password)
	req.Header.Add("Accept", "text/xml")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (vlc *vlc) Status() (Status, error) {
	url := fmt.Sprintf("%s/requests/status.xml", vlc.address)

	res, err := get(url, vlc.password)
	if err != nil {
		return Status{}, err
	}
	defer res.Body.Close()

	type xmlStatus struct {
		State string `xml:"state"`
	}

	var status xmlStatus
	err = xml.NewDecoder(res.Body).Decode(&status)
	if err != nil {
		return Status{}, err
	}

	return Status(status), nil
}

func (vlc *vlc) Playlist() (Playlist, error) {
	url := fmt.Sprintf("%s/requests/playlist.xml", vlc.address)

	res, err := get(url, vlc.password)
	if err != nil {
		return Playlist{}, err
	}
	defer res.Body.Close()

	type xmlPlaylist struct {
		Songs []struct {
			ID       int     `xml:"id,attr"`
			Name     string  `xml:"name,attr"`
			Current  boolean `xml:"current,attr"`
			Duration int     `xml:"duration,attr"`
			Uri      string  `xml:"uri,attr"`
		} `xml:"node>leaf"`
	}

	var plist xmlPlaylist
	err = xml.NewDecoder(res.Body).Decode(&plist)
	if err != nil {
		return Playlist{}, err
	}

	playlist := Playlist{}
	for _, s := range plist.Songs {
		song := PlaylistSong{
			ID:       s.ID,
			Name:     s.Name,
			Current:  s.Current.toBool(),
			Duration: s.Duration,
			Uri:      s.Uri,
		}

		playlist.Songs = append(playlist.Songs, song)
	}

	return playlist, nil
}

func (vlc *vlc) AddSong(url string, playNow bool) error {
	cmd := "in_enqueue"
	if playNow {
		cmd = "in_play"
	}

	uri := fmt.Sprintf("%s/requests/status.xml?command=%s&input=%s", vlc.address, cmd, url)

	res, err := get(uri, vlc.password)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (vlc *vlc) NextSong() error {
	uri := fmt.Sprintf("%s/requests/status.xml?command=pl_next", vlc.address)

	res, err := get(uri, vlc.password)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
