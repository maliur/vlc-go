package main

import (
	"fmt"

	"github.com/maliur/vlc-go"
)

func main() {
	vlcAddress := "http://192.168.1.12:8080"
	vlcPasword := "12345"

	vlcClient := vlc.NewClient(vlcAddress, vlcPasword)

	status, err := vlcClient.Status()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
}
