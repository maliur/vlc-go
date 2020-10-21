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
}
