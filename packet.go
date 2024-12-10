package rtp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type RTPPacket struct {
	Version        uint8
	Padding        bool
	Extension      bool
	CSRCCount      uint8
	Marker         bool
	PayloadType    uint8
	SequenceNumber uint16
	Timestamp      uint32
	SSRC           uint32
	Payload        []byte
	OptionalField  uint32
}

func (p *RTPPacket) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)
	header := uint16(p.Version<<14 | boolToUint8(p.Padding)<<13 | boolToUint8(p.Extension)<<12 |
		p.CSRCCount<<8 | boolToUint8(p.Marker)<<7 | p.PayloadType)
	binary.Write(buf, binary.BigEndian, header)
	binary.Write(buf, binary.BigEndian, p.SequenceNumber)
	binary.Write(buf, binary.BigEndian, p.Timestamp)
	binary.Write(buf, binary.BigEndian, p.SSRC)
	if p.Extension {
		binary.Write(buf, binary.BigEndian, p.OptionalField)
	}
	buf.Write(p.Payload)

	return buf.Bytes(), nil
}

func (p *RTPPacket) unmarshal(data []byte) error {
	if len(data) < 12 {
		return fmt.Errorf("invalid RTP packet length")
	}

	header := binary.BigEndian.Uint16(data[0:2])
	p.Version = uint8(header >> 14)
	p.Padding = (header>>13)&1 == 1
	p.Extension = (header>>12)&1 == 1
	p.CSRCCount = uint8((header >> 8) & 0x0F)
	p.Marker = (header>>7)&1 == 1
	p.PayloadType = uint8(header & 0x7F)

	p.SequenceNumber = binary.BigEndian.Uint16(data[2:4])
	p.Timestamp = binary.BigEndian.Uint32(data[4:8])
	p.SSRC = binary.BigEndian.Uint32(data[8:12])

	offset := 12
	if p.Extension {
		if len(data) < offset+4 {
			return fmt.Errorf("invalid RTP packet length for extension")
		}
		p.OptionalField = binary.BigEndian.Uint32(data[offset : offset+4])
		offset += 4
	}

	p.Payload = data[offset:]

	return nil
}

func Unmarshal(data []byte) (RTPPacket, error) {
	packet := RTPPacket{}
	err := packet.unmarshal(data)
	if err != nil {
		fmt.Printf("Failed to parse RTP packet: %v\n", err)
		return RTPPacket{}, err
	}
	return packet, nil
}

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
