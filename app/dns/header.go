package dns

import "fmt"

const HeaderSize = 12

type Header struct {
	ID      uint16 // Packet Identifier (16 bits)
	QR      uint8  // Query/Response Indicator (1 bit)
	OPCODE  uint8  // Operation Code (4 bits)
	AA      uint8  // Authoritative Answer (1 bit)
	TC      uint8  // Truncation (1 bit)
	RD      uint8  // Recursion Desired (1 bit)
	RA      uint8  // Recursion Available (1 bit)
	Z       uint8  // Reserved (3 bits)
	RCODE   uint8  // Response Code (4 bits)
	QDCOUNT uint16 // Question Count (16 bits)
	ANCOUNT uint16 // Answer Record Count (16 bits)
	NSCOUNT uint16 // Authority Record Count (16 bits)
	ARCOUNT uint16 // Additional Record Count (16 bits)
}

func (h *Header) String() string {
	return fmt.Sprintf(
		"Header{\n"+
			"  ID: %d\n"+
			"  QR: %d, OPCODE: %d, AA: %d, TC: %d, RD: %d\n"+
			"  RA: %d, Z: %d, RCODE: %d\n"+
			"  QDCOUNT: %d, ANCOUNT: %d, NSCOUNT: %d, ARCOUNT: %d\n"+
			"}",
		h.ID, h.QR, h.OPCODE, h.AA, h.TC, h.RD,
		h.RA, h.Z, h.RCODE,
		h.QDCOUNT, h.ANCOUNT, h.NSCOUNT, h.ARCOUNT,
	)
}

func parseHeader(data []byte) Header {
	if data == nil {
		return Header{}
	}

	h := Header{}
	h.ID = uint16(data[0])<<8 | uint16(data[1])
	h.QR = (data[2] >> 7) & 0x01
	h.OPCODE = data[2] >> 3 & 0x0F
	h.AA = (data[2] >> 2) & 0x01
	h.TC = (data[2] >> 1) & 0x01
	h.RD = data[2] & 0x01
	h.RA = (data[3] >> 7) & 0x01
	h.Z = (data[3] >> 4) & 0x07
	h.RCODE = data[3] & 0x0F
	h.QDCOUNT = uint16(data[4])<<8 | uint16(data[5])
	h.ANCOUNT = uint16(data[6])<<8 | uint16(data[7])
	h.NSCOUNT = uint16(data[8])<<8 | uint16(data[9])
	h.ARCOUNT = uint16(data[10])<<8 | uint16(data[11])

	return h
}

func (h *Header) ToBytes() []byte {
	bytes := make([]byte, 12)

	// ID (16 bitów)
	bytes[0] = byte(h.ID >> 8)
	bytes[1] = byte(h.ID & 0xFF)

	bytes[2] = byte(h.QR<<7 | h.OPCODE<<3 | h.AA<<2 | h.TC<<1 | h.RD)

	bytes[3] = byte(h.RA<<7 | h.Z<<4 | h.RCODE)

	bytes[4] = byte(h.QDCOUNT >> 8)
	bytes[5] = byte(h.QDCOUNT & 0xFF)
	bytes[6] = byte(h.ANCOUNT >> 8)
	bytes[7] = byte(h.ANCOUNT & 0xFF)
	bytes[8] = byte(h.NSCOUNT >> 8)
	bytes[9] = byte(h.NSCOUNT & 0xFF)
	bytes[10] = byte(h.ARCOUNT >> 8)
	bytes[11] = byte(h.ARCOUNT & 0xFF)

	return bytes
}
