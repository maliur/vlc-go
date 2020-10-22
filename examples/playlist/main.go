package main

import (
	"fmt"

	"github.com/maliur/vlc-go"
)

func main() {
	vlcClient := vlc.NewClient("http://localhost:8080", "12345")
	playlist, err := vlcClient.Playlist()
	if err != nil {
		fmt.Println(err)
	}

	// print playlist
	for _, p := range playlist.Songs {
		fmt.Println(p.ID)
		fmt.Println(p.Name)
		fmt.Println(p.Duration)
		fmt.Println(p.Current)
		fmt.Println(p.Uri)
		fmt.Print("\n\n")
	}

	// add song to playlist
	err = vlcClient.AddSong("https://youtu.be/dQw4w9WgXcQ", false)
	if err != nil {
		fmt.Println(err)
	}
}
