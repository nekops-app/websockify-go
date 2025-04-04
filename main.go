package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	vncAddr *string
)

func init() {
	vncAddr = flag.String("t", "127.0.0.1:5900", "vnc service address")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func forwardTcp(wsConn *websocket.Conn, conn net.Conn) {
	var tcpBuffer [1024]byte
	defer func() {
		if conn != nil {
			conn.Close()
		}
		if wsConn != nil {
			wsConn.Close()
		}
	}()
	for {
		if (conn == nil) || (wsConn == nil) {
			return
		}
		n, err := conn.Read(tcpBuffer[0:])
		if err != nil {
			LogError(fmt.Sprintf("reading from TCP failed: %s", err))
			return
		} else {
			if err := wsConn.WriteMessage(websocket.BinaryMessage, tcpBuffer[0:n]); err != nil {
				LogError(fmt.Sprintf("writing to WS failed: %s", err))
			}
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		LogError(fmt.Sprintf("failed to upgrade to WS: %s", err))
		return
	}

	vnc, err := net.Dial("tcp", *vncAddr)
	if err != nil {
		LogError(fmt.Sprintf("failed to dial the VNC server: %s", err))
	}

	go forwardTcp(ws, vnc)
}

func main() {
	flag.Parse()

	// Create random listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		LogError(fmt.Sprintf("failed to listen: %s", err))
	}

	LogStart(StartJSON{l.Addr().(*net.TCPAddr).IP.String(), l.Addr().(*net.TCPAddr).Port})

	http.HandleFunc("/websockify", serveWs)
	if err := http.Serve(l, nil); err != nil {
		LogError(fmt.Sprintf("failed to start http server: %s", err))
	}
}
