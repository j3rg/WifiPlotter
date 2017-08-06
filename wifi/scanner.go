// Copyright (c) 2017 Jorgen Ordonez
//
// Redistribution and use in source and binary forms, with or without
// modification, is permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

// Package wifi provides an interface to parse the
// output of the iwlist command.
package wifi

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

// WInterface is the wifi card that will be use to
// scan access points
type WInterface interface {
	Scan() error
	// Results() []AccessPoint
}

// AccessPoint defines a wifi access point
type AccessPoint struct {
	// MAC address
	Address   string
	SSID      string
	Channel   int
	Frequency float32
	Quality   quality
	// Currently this is just a true/flase flag
	// I shall find a way to get the encryption the
	// access point is using
	Encryption bool
}

type quality struct {
	Percent float32
	// The signal is in units of decibles dB
	Signal int
}

// Scanner implements WInterface
type Scanner struct {
	iface   *net.Interface
	results []AccessPoint
}

// regexp constants used to parse output
const (
	macAddrExp = "Cell [0-9]{2} - Address: "
)

// New creates a new scanner
func New(iface string) (*Scanner, error) {
	var s Scanner
	networkCard, err := net.InterfaceByName(iface)

	s.iface = networkCard

	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Scan simply scans for wireless AP
func (s *Scanner) Scan() error {
	if s.iface == nil {
		return errors.New("Invalid network card")
	}

	iwlistCmd := exec.Command("iwlist", s.iface.Name, "scan")
	iwlistOut, err := iwlistCmd.Output()

	if err != nil {
		fmt.Println(err)
	}

	s.results = parseOutput(&iwlistOut)
	fmt.Println(s.results)
	return nil
}

// Results return the resuts of scan
func (s *Scanner) Results() []AccessPoint {
	return s.results
}

func parseOutput(output *[]byte) []AccessPoint {

	accesspoints := make([]AccessPoint, 1, 50)

	results := strings.Split(string(*output), "\n")

	macAddrRegExp := regexp.MustCompile(macAddrExp)

	for _, result := range results {

		if macAddrRegExp.MatchString(result) {

			var ap AccessPoint

			idx := macAddrRegExp.FindStringIndex(result)
			ap.Address = result[idx[1]:]

			accesspoints = append(accesspoints, ap)
		}
	}

	return accesspoints
}
