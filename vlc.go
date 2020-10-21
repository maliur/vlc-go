package vlc

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Vlc interface {
	Playlist() (Playlist, error)
	Status() (Status, error)
}

type vlc struct {
	address  string
	password string
}

// Create a new Vlc client
func New(address, password string) *vlc {
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

// Status
type Status struct {
	State string `xml:"state"`
}

// Get the current status from vlc
func (vlc *vlc) Status() (Status, error) {
	url := fmt.Sprintf("%s/requests/status.xml", vlc.address)

	res, err := get(url, vlc.password)
	if err != nil {
		return Status{}, err
	}
	defer res.Body.Close()

	var status Status
	err = xml.NewDecoder(res.Body).Decode(&status)
	if err != nil {
		return Status{}, err
	}

	return status, nil
}

// Playlist
type Playlist struct {
	Songs []PlaylistSong `xml:"node>leaf"`
}

// Playlist song
type PlaylistSong struct {
	ID       int    `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	Current  string `xml:"current,attr"`
	Duration int    `xml:"duration,attr"`
}

func (vlc *vlc) Playlist() (Playlist, error) {
	url := fmt.Sprintf("%s/requests/playlist.xml", vlc.address)

	res, err := get(url, vlc.password)
	if err != nil {
		return Playlist{}, err
	}
	defer res.Body.Close()

	var playlist Playlist
	err = xml.NewDecoder(res.Body).Decode(&playlist)
	if err != nil {
		return Playlist{}, err
	}

	return playlist, nil
}
