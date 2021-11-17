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
	"fmt"
	"io"
	"log"
	"net"
	"net/html"
	"net/http"
	"net/url"
	"os"
	"strings"
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

	// client := &http.Client{
	// 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	// 		return errors.New("net/http: use last response")
	// 	},
	// }

	// newURL, err := url.Parse(request.RequestURI)
	// if err != nil {
	// 	newResponse := []byte("HTTP 500 Internal Error")
	// 	conn.Write(newResponse)
	// 	return
	// }
	// request.URL = newURL
	// request.RequestURI = ""
	// request.Header.Add("Host", request.Host)
	// request.Header.Add("Connection", "close")

	// resp, err := client.Do(request)
	// if err != nil {
	// 	if !strings.Contains(err.Error(), "net/http: use last response") {
	// 		newResponse := []byte("HTTP 500 Internal Error")
	// 		conn.Write(newResponse)
	// 		return
	// 	}
	// }
	// resp.Write(conn)

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
	// fmt.Println(request)
	sreader := bufio.NewReader(sconn)
	sresponse, err := http.ReadResponse(sreader, request)
	if err != nil {
		log.Fatal(err)
	}
	sresponse.Write(conn)
	sconn.Close()
	go DNS(sresponse.Body)
	return
}

func DNS(r io.Reader) {

	z := html.NewTokenizer(r)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			name, hasAttr := z.TagName()
			if string(name) == "a" {
				if hasAttr {
					key, val, hasAttr := z.TagAttr()
					for hasAttr {
					// need to check start with http
						if key == "href" && strings.HasPrefix(string(val), "http") {
							go net.LookupHost(string(val))
						}
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
