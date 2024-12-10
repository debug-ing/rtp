package rtp

import (
	"fmt"
	"net"
	"time"
)

type RTP struct {
	Port                 int
	handleClientFunction func(conn net.PacketConn, addr net.Addr, data []byte, rtp RTPPacket)
}

func Init(port int, handleClient func(conn net.PacketConn, addr net.Addr, data []byte, rtp RTPPacket)) RTP {
	return RTP{Port: port, handleClientFunction: handleClient}
}

func (r *RTP) Run() {
	addr := ":" + fmt.Sprint(r.Port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		panic(fmt.Sprintf("Failed to listen on %s: %v", addr, err))
	}
	defer conn.Close()

	fmt.Printf("Server is running on %s\n", addr)

	buffer := make([]byte, 1500)

	for {
		n, clientAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Printf("Error reading from client: %v\n", err)
			continue
		}
		go r.handleClient(conn, clientAddr, buffer[:n])
	}
}
func (r *RTP) handleClient(conn net.PacketConn, addr net.Addr, data []byte) {
	packet, err := r.decode(data)
	if err != nil {
		fmt.Printf("Failed to parse RTP packet: %v\n", err)
		return
	}
	r.handleClientFunction(conn, addr, data, packet)
}

func (r *RTP) decode(data []byte) (RTPPacket, error) {
	// packet := &server.RTPPacket{}
	packet, err := Unmarshal(data)
	if err != nil {
		fmt.Printf("Failed to parse RTP packet: %v\n", err)
		return RTPPacket{}, err
	}
	return packet, nil
}

func Send(conn net.PacketConn, addr net.Addr, packet RTPPacket, data []byte) {
	response := &RTPPacket{
		Version:        2,
		SequenceNumber: packet.SequenceNumber + 1,
		Timestamp:      uint32(time.Now().Unix()),
		SSRC:           packet.SSRC,
		Payload:        data,
	}
	responseData, _ := response.Marshal()
	conn.WriteTo(responseData, addr)
}
