package help

func HowToRun() string {
	content := `How to run
	simple-chatroom-in-GoLang-RPC --mode [server/client] --port [port_number] --server-port [server_port_number] --username [username] 
for example:
	simple-chatroom-in-GoLang-RPC --mode server --port 8080
	simple-chatroom-in-GoLang-RPC --mode client --port 8081 --server-ip 127.0.0.1 --server-port 8080 --username John
`
	return content
}
