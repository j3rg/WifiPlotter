package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/j3rg/WifiPlotter/wifi"
)

func main() {

	currentUser, err := user.Current()

	if err != nil {
		fmt.Println(err)
	}

	if currentUser.Uid > "0" {
		fmt.Println("Program must be run as root user")
		os.Exit(1)
	}

	fmt.Println("~WifiPlotter~")
	fmt.Printf("=============\n\n")

	w, err := wifi.New("wlan1")
	if err != nil {
		fmt.Println(err)
	}

	w.Scan()
	accesspoints := w.Results()

	for _, ap := range accesspoints {
		fmt.Println(ap.SSID)
		fmt.Println(strings.Repeat("=", len(ap.SSID)))
		fmt.Println("MAC Address: ", ap.Address)
		fmt.Println("Channel:     ", ap.Channel)
		fmt.Println("Frequency:   ", ap.Frequency, "GHz")
		fmt.Println("Quality:")
		fmt.Printf("  Percent:    %v%%\n", ap.Quality.Percent)
		fmt.Printf("  Signal:     %v dBm\n\n\n", ap.Quality.Signal)
	}

}
