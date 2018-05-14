// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package pkg

import (
	"fmt"
	"log"
	"net"
)

// GetPublicIPv6Address returns the first local ip address matching the public subnet
func GetPublicIPv6Address(publicIPNet net.IPNet) (ipNet net.IP, err error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("unable to get interfaces: %v", err)
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // iface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback device
		}

		addrs, err := iface.Addrs()
		if err != nil {
			log.Printf("WARN unable to get addresses of %s: %v", iface.Name, err)
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if publicIPNet.Contains(v.IP) {
					return v.IP, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no ip found which is in %s", publicIPNet.String())
}
