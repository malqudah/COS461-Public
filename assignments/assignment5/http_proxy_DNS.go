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
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func handleConnection(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewReader(conn)
	request, err := http.ReadRequest(reader)

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

	relativeURL, err := url.Parse(request.URL.Path)
	if err != nil {
		log.Fatal(err)
	}
	request.URL = relativeURL
	request.Header.Set("Connection", "close")

	hostport := request.Host + ":http"
	sconn, err := net.Dial("tcp", hostport)
	if err != nil {
		log.Fatal(err)
	}
	request.Write(sconn)

	response := new(bytes.Buffer)
	io.Copy(response, sconn)
	dnsReader := strings.NewReader(response.String())
	go DNS(dnsReader)
	io.Copy(conn, response)
	sconn.Close()
	return
}

func DNS(r io.Reader) {

	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return
		} else {
			name, hasAttr := z.TagName()
			if string(name) == "a" {
				if hasAttr {
					var key, val []byte
					key, val, hasAttr = z.TagAttr()
					if string(key) == "href" && (string(val[:4]) == "http") {
						theURL, err := url.Parse(string(val))
						if err != nil {
							return
						}
						net.LookupHost(theURL.Host)
					}
				}
			}
		}
	}
	return

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
