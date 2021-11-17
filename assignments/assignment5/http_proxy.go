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
	// if err != nil {
	// 	newResponse := []byte("HTTP 500 Internal Error")
	// 	conn.Write(newResponse)
	// 	return
	// }
	// request.RequestURI = ""
	// resp, err := client.Do(request)
	// if err != nil {
	// 	if !strings.Contains(err.Error(), "net/http: use last response") {
	// 		fmt.Println(err)
	// 		newResponse := []byte("HTTP 500 Internal Error")
	// 		conn.Write(newResponse)
	// 		return
	// 	}
	// }

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
