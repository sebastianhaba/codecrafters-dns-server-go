package dns

import "fmt"

type Message struct {
	Header   Header
	Question []Question
}

func NewMessage(data []byte) *Message {
	message := &Message{}
	message.parse(data)
	return message
}

func (m *Message) ToBytes() []byte {
	bytes := m.Header.ToBytes()

	for i := 0; i < len(m.Question); i++ {
		bytes = append(bytes, m.Question[i].ToBytes()...)
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
