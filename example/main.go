package main

import (
	"fmt"
	"net"

	rtp "github.com/debug-ing/rtp"
)

func main() {
	server := rtp.Init(5004, func(conn net.PacketConn, addr net.Addr, data []byte, rtpModel rtp.RTPPacket) {
		fmt.Println("Received data from", rtpModel.Version)
		rtp.Send(conn, addr, rtpModel, []byte("Response from server"))
	})
	server.Run()
}
