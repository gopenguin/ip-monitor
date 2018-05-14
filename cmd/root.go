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

package cmd

import (
	"fmt"
	"os"

	"log"
	"net"
	"text/template"
	"time"

	"github.com/gopenguin/ip-monitor/pkg"
	"github.com/spf13/cobra"
)

var cfgFile string

var defaultPublicNet = net.IPNet{IP: net.ParseIP("2000::"), Mask: net.IPMask(net.ParseIP("e000::"))}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ip-monitor",
	Short: "Monitor all available ip addresses and update the corresponding dyndns provider",
	Run:   rootCmd,
}

func rootCmd(cmd *cobra.Command, args []string) {
	publicNet, err := cmd.Flags().GetIPNet("publicNet")
	if err != nil {
		log.Printf("Error getting public net: %v", err)

		publicNet = defaultPublicNet
		log.Printf("Using default net %s", defaultPublicNet.String())
	} else {
		log.Printf("Using net %s", publicNet.String())
	}
	urlTemplateString, err := cmd.Flags().GetString("urlTemplate")
	if err != nil {
		log.Printf("Unable to get urlTemplate %v", err)
		return
	}

	urlTemplate := template.Must(template.New("urlTemplate").Parse(urlTemplateString))

	domainString, err := cmd.Flags().GetString("domain")
	if err != nil {
		log.Printf("Unable to get the domain: %v", err)
		return
	}

	tokenString, err := cmd.Flags().GetString("token")
	if err != nil {
		log.Printf("Unable to get the token: %v", err)
		return
	}

	interval, err := cmd.Flags().GetDuration("interval")
	if err != nil {
		log.Printf("Unable to get interval: %v", err)
	}

	var expectedResponse *string
	if cmd.Flags().Changed("expectedResponse") {
		expectedResponseString, _ := cmd.Flags().GetString("expectedResponse")
		expectedResponse = &expectedResponseString
	}

	updater := pkg.NewUpdater(urlTemplate, domainString, tokenString)

	for {
		func() {
			publicIP, err := pkg.GetPublicIPv6Address(publicNet)
			if err != nil {
				log.Printf("Unable to get public ip: %v", err)
				return
			}

			err = updater.UpdateIP(publicIP.String(), expectedResponse)
			if err != nil {
				log.Printf("Unable to update the ip: %v", err)
				return
			}

			log.Printf("Updated %s to %s", domainString, publicIP.String())
		}()

		time.Sleep(interval)
	}
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ip-monitor.yaml)")

	RootCmd.Flags().String("urlTemplate", "", "the address to call to update the ip address; the ipv6 address is available via '{.IPv6}'")
	RootCmd.Flags().String("token", "", "the token required to authenticate; in the url pattern available via '{.Token}'")
	RootCmd.Flags().String("domain", "", "the domain to update, available via '{.Domain}'")
	RootCmd.Flags().IPNet("publicNet", defaultPublicNet, "the subnet which contains public ipv6 addresses")
	RootCmd.Flags().Duration("interval", 10*time.Minute, "the time between two updates")
	RootCmd.Flags().String("expectedResponse", "", "the expected  response from the server; if empty, it will be ignored")
}
