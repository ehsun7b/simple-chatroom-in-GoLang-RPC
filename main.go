package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	c "simple-chatroom-in-GoLang-RPC/client"
	h "simple-chatroom-in-GoLang-RPC/help"
	"simple-chatroom-in-GoLang-RPC/server"
	"strings"
	"time"
)

var (
	mode       string
	serverIp   string
	username   string
	port       int
	serverPort int
	showHelp   bool
)

func main() {
	flag.StringVar(&mode, "mode", "client", "determines the mode [server/client]")
	flag.StringVar(&serverIp, "server-ip", "127.0.0.1", "the ip address of the server")
	flag.StringVar(&username, "username", "user", "the user identity (client mode only)")
	flag.IntVar(&port, "port", 8080, "the port number to listen")
	flag.IntVar(&serverPort, "server-port", 8080, "the port number of the server")
	flag.BoolVar(&showHelp, "showHelp", false, "display showHelp")
	flag.Parse()

	defer cleanup()

	if showHelp {
		log.Println(h.HowToRun())
	} else {
		setCleanupOnCtrlC()

		switch mode {
		case "server":
			log.Println("server mode")
			server.StartServer(port)
			break
		case "client":
			log.Printf("client mode username: %v\n", username)
			for !c.SignInAtServer(serverIp, serverPort, username, "127.0.0.1", port) {
				time.Sleep(time.Second * 3)
			}

			go c.StartClient(port)
			scanner := bufio.NewScanner(os.Stdin)
			for {
				fmt.Print(">")
				scanner.Scan()
				message := scanner.Text()
				if strings.TrimSpace(message) != "" {
					c.SendToServer(serverIp, serverPort, username, message)
				}
			}
			break
		}
	}
}

func setCleanupOnCtrlC() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				cleanup()
			}
		}
	}()
}

func cleanup() {
	switch mode {
	case "client":
		c.SignOutAtServer(serverIp, serverPort, username)
	}

	log.Println("bye")
	os.Exit(0)
}
