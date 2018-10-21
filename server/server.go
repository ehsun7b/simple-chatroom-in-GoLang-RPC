package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"simple-chatroom-in-GoLang-RPC/client"
	"strconv"
	"strings"
)

type Server struct {
	clients map[string]client.Client
}

func StartServer(port int) {
	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Printf("listening on port %v\n", port)
	http.Serve(l, nil)
}

func sendToClient(ip string, port int, from string, message string) {
	client, err := rpc.DialHTTP("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		log.Printf("dialing: %v\n", err.Error())
	}

	var response int

	error := client.Call("Client.Deliver", fmt.Sprintf("%v,%v", from, message), &response)
	if error != nil {
		log.Printf("calling... %v", error.Error())
	}
}

func (s *Server) SignIn(request string, response *string) error {
	if request == "" {
		*response = "empty request. sign in failed"
	} else {
		parts := strings.Split(request, ",")

		if len(parts) != 3 {
			*response = "not enough information. sign in failed"
			log.Printf("failed sign in: %v", request)
			return nil
		} else {
			newClient := client.Client{Username: parts[0], Ip: parts[1]}
			if port, e := strconv.Atoi(parts[2]); e == nil {
				newClient.Port = port

				if s.clients == nil {
					s.clients = make(map[string]client.Client)
				}

				s.clients[parts[0]] = newClient
				*response = "signed in successfully"
				log.Printf("%+v signed in successfully. Clients: %+v\n", newClient, s.clients)
			} else {
				*response = "wrong port. sign in failed"
				log.Printf("failed sign in: %v", request)
				return nil
			}
		}
	}

	return nil
}

func (s *Server) SignOut(request string, response *string) error {
	if request == "" {
		*response = "empty request. sign out failed"
	} else {
		if s.clients == nil {
			*response = "you are not signed in"
		} else if _, ok := s.clients[request]; !ok {
			*response = "you are not signed in"
		} else {
			delete(s.clients, request)
			*response = "signed out successfully"
			log.Printf("%v signed out. Clients: %v", request, s.clients)
		}
	}

	return nil
}

func (s *Server) Send(request string, response *string) error {
	log.Printf("incoming message: %v", request)
	parts := strings.Split(request, ",")

	for username, c := range s.clients {
		if username != parts[0] {
			sendToClient(c.Ip, c.Port, parts[0], parts[1])
		}
	}

	return nil
}
