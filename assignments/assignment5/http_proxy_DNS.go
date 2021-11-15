/*****************************************************************************
 * http_proxy_DNS.go                                                                 
 * Names: Mohammad Alqudah, Jonathan Salama
 * NetIds: malqudah, jjsalama
 *****************************************************************************/

// TODO: implement an HTTP proxy with DNS Prefetching

// Note: it is highly recommended to complete http_proxy.go first, then copy it
// with the name http_proxy_DNS.go, thus overwriting this file, then edit it
// to add DNS prefetching (don't forget to change the filename in the header
// to http_proxy_DNS.go in the copy of http_proxy.go)

package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"net/html"
)

func handleConnection(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewReader(conn)
	request, err := http.ReadRequest(reader)

	// come back here
	if err != nil {
		newResponse := []byte("HTTP 500 Internal Error")
		conn.Write(newResponse)
		return
	}

	if request.Method != "GET" {
		newResponse := []byte("HTTP 500 Internal Error")
		conn.Write(newResponse)
		return
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("net/http: use last response")
		},
	}

	newURL, err := url.Parse(request.RequestURI)
	if err != nil {
		newResponse := []byte("HTTP 500 Internal Error")
		conn.Write(newResponse)
		return
	}
	request.URL = newURL
	request.RequestURI = ""
	request.Header.Add("Host", request.Host)
	request.Header.Add("Connection", "close")

	resp, err := client.Do(request)
	if err != nil {
		if !strings.Contains(err.Error(), "net/http: use last response") {
			newResponse := []byte("HTTP 500 Internal Error")
			conn.Write(newResponse)
			return
		}
	}
	resp.Write(conn)
	return
}

func DNS(r io.Reader) {

	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return z.Err()
		case html.StartTagToken:
			name, hasAttr := z.TagName()
			if name == 'a' {
				if hasAttr {
					key, val, hasAttr := z.TagAttr()
					// need to check start with http
					if key == "href" {
						go net.LookupHost(val)
					}
				}
			}
		}
	}
	
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server-go [server port]")
	}
	server_port := os.Args[1]

	ln, err := net.Listen("tcp", ":"+server_port)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on port: " + server_port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		go handleConnection(conn)
	}

}
