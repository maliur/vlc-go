package main

import (
	"fmt"

	"github.com/maliur/vlc-go"
)

func main() {
	vlcAddress := "http://127.0.0.1"
	vlcPort := "8080"
	vlcPasword := "12345"

	vlc := vlc.New(vlcAddress, vlcPort, vlcPasword)

	status, err := vlc.Status()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
}
