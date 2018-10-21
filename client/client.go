package client

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
)

type Client struct {
	Username string
	Ip       string
	Port     int
}

func SignInAtServer(serverIp string, serverPort int, username string, ip string, port int) bool {
	log.Printf("signing in...\n")

	client, err := rpc.DialHTTP("tcp", serverIp+":"+strconv.Itoa(serverPort))
	if err != nil {
		log.Printf("dialing: %v\n", err.Error())
		return false
	}

	request := fmt.Sprintf("%v,%v,%v", username, ip, port)
	var response string

	error := client.Call("Server.SignIn", request, &response)
	if error != nil {
		log.Printf("calling... %v", error.Error())
		return false
	}

	log.Printf("response from server: %v\n", response)
	return true
}

func SignOutAtServer(serverIp string, serverPort int, username string) bool {
	log.Printf("signing out...\n")

	client, err := rpc.DialHTTP("tcp", serverIp+":"+strconv.Itoa(serverPort))
	if err != nil {
		log.Printf("dialing: %v\n", err.Error())
		return false
	}

	var response string

	error := client.Call("Server.SignOut", username, &response)
	if error != nil {
		log.Printf("calling... %v", error.Error())
		return false
	}

	log.Printf("response from server: %v\n", response)
	return true
}

func StartClient(port int) {
	client := new(Client)
	rpc.Register(client)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Printf("listening on port %v\n", port)
	http.Serve(l, nil)
}

func SendToServer(serverIp string, serverPort int, username string, message string) {
	client, err := rpc.DialHTTP("tcp", serverIp+":"+strconv.Itoa(serverPort))
	if err != nil {
		log.Printf("dialing: %v\n", err.Error())
	}

	var response string

	error := client.Call("Server.Send", fmt.Sprintf("%v,%v", username, message), &response)
	if error != nil {
		log.Printf("calling... %v", error.Error())
	}
}

func (c *Client) Deliver(message string, response *int) error {
	if message != "" {
		parts := strings.Split(message, ",")
		if len(parts) == 2 {
			fmt.Printf("\n%v: %v\n>", parts[0], parts[1])
		}
	}

	return nil
}
