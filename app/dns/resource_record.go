package dns

import (
	"encoding/binary"
	"fmt"
	"net"
)

type ResourceRecord struct {
	NAME     string
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

func NewARecord(name string, ip net.IP, ttl uint32) ResourceRecord {
	ipBytes := ip.To4()

	rr := ResourceRecord{
		NAME:     name,
		TYPE:     1,
		CLASS:    1,
		TTL:      ttl,
		RDLENGTH: 4,
		RDATA:    ipBytes,
	}

	return rr
}

func (rr *ResourceRecord) ToBytes() []byte {
	var bytes []byte

	bytes = append(bytes, domainNameToBytes(rr.NAME)...)
	bytes = append(bytes, byte(rr.TYPE>>8), byte(rr.TYPE&0xFF))
	bytes = append(bytes, byte(rr.CLASS>>8), byte(rr.CLASS&0xFF))

	ttlBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ttlBytes, rr.TTL)
	bytes = append(bytes, ttlBytes...)
	bytes = append(bytes, byte(rr.RDLENGTH>>8), byte(rr.RDLENGTH&0xFF))
	bytes = append(bytes, rr.RDATA...)

	return bytes
}

func (rr *ResourceRecord) String() string {
	var rdataStr string
	if rr.TYPE == 1 && len(rr.RDATA) == 4 {
		rdataStr = fmt.Sprintf("%d.%d.%d.%d", rr.RDATA[0], rr.RDATA[1], rr.RDATA[2], rr.RDATA[3])
	} else {
		rdataStr = fmt.Sprintf("%v", rr.RDATA)
	}
	return fmt.Sprintf("ResourceRecord{NAME: %s, TYPE: %d, CLASS: %d, TTL: %d, RDLENGTH: %d, RDATA: %s}",
		rr.NAME, rr.TYPE, rr.CLASS, rr.TTL, rr.RDLENGTH, rdataStr)
}
