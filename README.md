# RTP Server Go


## Introduction

RealTime Protocol implementation based on [RFC 8860](https://datatracker.ietf.org/doc/html/rfc8860) in Golang. 
Get Message and Send Message.


## Example

```go
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
```

