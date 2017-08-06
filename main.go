package main

import (
	"fmt"
	"os"
	"os/user"

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

	fmt.Println("WifiPlotter")
	fmt.Println("===========")

	w, err := wifi.New("wlan1")
	if err != nil {
		fmt.Println(err)
	}

	w.Scan()
	accesspoints := w.Results()

	for _, ap := range accesspoints {
		fmt.Println(ap.Address)
	}

}
