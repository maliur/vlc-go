package vlc

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Vlc interface {
	Status() (Status, error)
	Playlist() (Playlist, error)
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
