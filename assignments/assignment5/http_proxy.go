/*****************************************************************************
 * http_proxy.go
 * Names: Mohammad Alqudah, Jonathan Salama
 * NetIds: malqudah, jjsalama
 *****************************************************************************/

// TODO: implement an HTTP proxy

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

func handleConnection(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewReader(conn)
	request, err := http.ReadRequest(reader)

	// come back here
	if err != nil {
		log.Fatal(err)
		return
	}

	if request.Method != "GET" {
		newResponse := []byte("HTTP 500 Internal Error")
		conn.Write(newResponse)
		return
	}

	client := &http.Client{}
	// newRequest, err := http.NewRequest(request.Method, request.URL.String(), request.Body)

	// newRequest.Header = request.Header

	newURL, err := url.Parse(request.RequestURI)
	request.URL = newURL
	fmt.Println(newURL.Path)
	request.RequestURI = ""
	request.Header.Add("Host", request.Host)
	request.Header.Add("Connection", "close")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(request)
	fmt.Println()
	resp.Write(os.Stdout)
	resp.Write(conn)
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
