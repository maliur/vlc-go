package main

import (
	"fmt"

	"github.com/maliur/vlc-go"
)

func main() {
	vlcAddress := "http://192.168.1.16:8080"
	vlcPasword := "12345"

	vlc := vlc.New(vlcAddress, vlcPasword)

	status, err := vlc.Status()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)

	playlist, err := vlc.Playlist()
	if err != nil {
		fmt.Println(err)
	}

	for _, p := range playlist.Songs {
		fmt.Println(p.ID)
		fmt.Println(p.Name)
		fmt.Println(p.Duration)
		fmt.Print("\n\n")
	}

	fmt.Println(status)
}
