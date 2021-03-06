[![Actions Status](https://github.com/maliur/vlc-go/workflows/build/badge.svg)](https://github.com/maliur/vlc-go/actions) [![Actions Status](https://github.com/maliur/vlc-go/workflows/tests/badge.svg)](https://github.com/maliur/vlc-go/actions)

# vlcgo

A Library to interact with [VLC media player](https://www.videolan.org/)

## Install
```bash
go get -u github.com/maliur/vlc-go
```

## Usage
Get the current status of VLC
```go
func main() {
    vlcAddress := "http://localhost:8080"
    vlcPassword := "12345"

    vlcClient := vlc.New(vlcAddress, vlcPassword)

    status, err := vlcClient.Status()

    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(status)
}
```
See more in [examples](https://github.com/maliur/vlc-go/tree/master/examples)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
Released under [MIT license](https://raw.githubusercontent.com/maliur/vlc-go/master/LICENSE)
