package vlc

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

type Status struct {
	State       string
	Information struct {
		Category []struct {
			Info []struct {
				Text string
				Name string
			}
		}
	}
}

func (s Status) Song() string {
	if len(s.Information.Category) > 0 {
		category := s.Information.Category[0]

		for _, info := range category.Info {
			if info.Name == "title" {
				return info.Text
			}
		}
	}

	return ""
}

func (s Status) Artist() string {
	if len(s.Information.Category) > 0 {
		category := s.Information.Category[0]
		for _, info := range category.Info {
			if info.Name == "artist" {
				return info.Text
			}
		}
	}

	return ""
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
	AddSong(uri string, playNow bool) error
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

func get(uri, password string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth("", url.QueryEscape(password))
	req.Header.Add("Accept", "text/xml")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (vlc *vlc) Status() (Status, error) {
	uri := fmt.Sprintf("%s/requests/status.xml", vlc.address)

	res, err := get(uri, vlc.password)
	if err != nil {
		return Status{}, err
	}
	defer res.Body.Close()

	type xmlStatus struct {
		State       string `xml:"state"`
		Information struct {
			Category []struct {
				Info []struct {
					Text string `xml:",chardata"`
					Name string `xml:"name,attr"`
				} `xml:"info"`
			} `xml:"category"`
		} `xml:"information"`
	}

	var status xmlStatus
	err = xml.NewDecoder(res.Body).Decode(&status)
	if err != nil {
		return Status{}, err
	}

	return Status(status), nil
}

func (vlc *vlc) Playlist() (Playlist, error) {
	uri := fmt.Sprintf("%s/requests/playlist.xml", vlc.address)

	res, err := get(uri, vlc.password)
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

func (vlc *vlc) AddSong(uri string, playNow bool) error {
	cmd := "in_enqueue"
	if playNow {
		cmd = "in_play"
	}

	statusUri := fmt.Sprintf("%s/requests/status.xml?command=%s&input=%s", vlc.address, cmd, uri)

	res, err := get(statusUri, vlc.password)
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
