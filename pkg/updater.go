// Copyright © 2018 gopenguin
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
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/certifi/gocertifi"
)

// Updater updates the dns record at a dyndns provider using a rest api
type Updater struct {
	netClient   *http.Client
	urlTemplate *template.Template
	token       string
	domain      string
}

// NewUpdater creates a new client for updating the ip address at a dynamic domain hoster.
func NewUpdater(urlTemplate *template.Template, domain string, token string) *Updater {
	certPool, err := gocertifi.CACerts()
	if err != nil {
		panic(err)
	}

	return &Updater{
		netClient: &http.Client{
			Timeout: time.Second * 30,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: certPool},
			},
		},
		urlTemplate: urlTemplate,
		token:       token,
		domain:      domain,
	}
}

// UpdateIP updates the public ip of a dns record
func (u *Updater) UpdateIP(ip string, expectedResponse *string) error {
	url, err := formatURL(u.urlTemplate, u.token, u.domain, ip)
	if err != nil {
		return nil
	}

	res, err := u.netClient.Get(url)
	if err != nil {
		return fmt.Errorf("unable to update the record: %v", err)
	}

	defer res.Body.Close()

	buf, _ := ioutil.ReadAll(res.Body)
	serverResponse := string(buf)

	if res.StatusCode < 200 && res.StatusCode >= 300 {
		return fmt.Errorf("unable to update the record, status %d: %v", res.StatusCode, serverResponse)
	}

	if expectedResponse != nil && serverResponse != *expectedResponse {
		return fmt.Errorf("Server response didn't match the expected: '%s' != '%s'", serverResponse, *expectedResponse)
	}

	return nil
}

func formatURL(urlTemplate *template.Template, token string, domain string, ip string) (urlString string, err error) {
	templateData := struct {
		Token  string
		Domain string
		IPv6   string
	}{
		Token:  url.QueryEscape(token),
		Domain: url.QueryEscape(domain),
		IPv6:   url.QueryEscape(ip),
	}

	buffer := new(bytes.Buffer)
	err = urlTemplate.Execute(buffer, templateData)
	if err != nil {
		return "", fmt.Errorf("unable to execute url template: %v", err)
	}

	return buffer.String(), nil
}
