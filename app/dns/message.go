package dns

import "fmt"

type Message struct {
	Header   Header
	Question []Question
	Answer   []ResourceRecord
}

func NewMessage(data []byte) *Message {
	message := &Message{}
	message.parse(data)

	message.Header.QR = 1

	if message.Header.OPCODE != 0 {
		message.Header.RCODE = 4
	} else {
		message.Header.RCODE = 0
	}

	return message
}

func (m *Message) ToBytes() []byte {
	m.Header.ANCOUNT = uint16(len(m.Answer))

	bytes := m.Header.ToBytes()

	for i := 0; i < len(m.Question); i++ {
		bytes = append(bytes, m.Question[i].ToBytes()...)
	}

	for i := 0; i < len(m.Answer); i++ {
		bytes = append(bytes, m.Answer[i].ToBytes()...)
	}

	return bytes
}

func (m *Message) parse(data []byte) {
	if data == nil || len(data) < HeaderSize {
		return
	}

	m.Header = parseHeader(data)

	offset := HeaderSize
	m.Question = make([]Question, m.Header.QDCOUNT)
	for i := 0; i < int(m.Header.QDCOUNT); i++ {
		question, newOffset := parseQuestion(data, offset)
		m.Question[i] = question
		offset = newOffset
	}

	if m.Header.ANCOUNT > 0 {
		m.Answer = make([]ResourceRecord, m.Header.ANCOUNT)
		for i := 0; i < int(m.Header.ANCOUNT); i++ {
			answer, newOffset := parseResourceRecord(data, offset)
			m.Answer[i] = answer
			offset = newOffset
		}
	}
}

func (m *Message) String() string {
	var result string
	result = fmt.Sprintf(
		"Message{\n"+
			"  Header: %s\n",
		m.Header.String(),
	)

	if len(m.Question) > 0 {
		result += "  Questions: [\n"
		for i, q := range m.Question {
			result += fmt.Sprintf("    %d: %s\n", i, q.String())
		}
		result += "  ]\n"
	}

	result += "}"
	return result
}

func parseResourceRecord(data []byte, offset int) (ResourceRecord, int) {
	record := ResourceRecord{}

	name, newOffset := parseDomainName(data, offset)
	record.NAME = name
	offset = newOffset

	if offset+10 > len(data) {
		return record, offset
	}

	record.TYPE = uint16(data[offset])<<8 | uint16(data[offset+1])
	offset += 2

	record.CLASS = uint16(data[offset])<<8 | uint16(data[offset+1])
	offset += 2

	record.TTL = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
	offset += 4

	record.RDLENGTH = uint16(data[offset])<<8 | uint16(data[offset+1])
	offset += 2

	if offset+int(record.RDLENGTH) > len(data) {
		return record, offset
	}

	record.RDATA = make([]byte, record.RDLENGTH)
	copy(record.RDATA, data[offset:offset+int(record.RDLENGTH)])
	offset += int(record.RDLENGTH)

	return record, offset
}
