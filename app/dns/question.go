package dns

import (
	"fmt"
	"strings"
)

type Question struct {
	QNAME  string
	QTYPE  uint16
	QCLASS uint16
}

func (q *Question) String() string {
	return fmt.Sprintf("Question{QNAME: %s, QTYPE: %d, QCLASS: %d}", q.QNAME, q.QTYPE, q.QCLASS)
}

func (q *Question) ToBytes() []byte {
	var bytes []byte
	bytes = append(bytes, domainNameToBytes(q.QNAME)...)
	bytes = append(bytes, byte(q.QTYPE>>8), byte(q.QTYPE&0xFF))
	bytes = append(bytes, byte(q.QCLASS>>8), byte(q.QCLASS&0xFF))

	return bytes
}

func parseQuestion(data []byte, offset int) (Question, int) {
	question := Question{}
	endOffset := offset

	name, newOffset := parseDomainName(data, offset)
	question.QNAME = name
	endOffset = newOffset

	if len(data) < endOffset+4 {
		return question, endOffset
	}

	question.QTYPE = uint16(data[endOffset])<<8 | uint16(data[endOffset+1])
	endOffset += 2

	question.QCLASS = uint16(data[endOffset])<<8 | uint16(data[endOffset+1])
	endOffset += 2

	return question, endOffset

}

func parseDomainName(data []byte, offset int) (string, int) {
	if offset >= len(data) {
		return "", offset
	}

	var name string
	position := offset
	length := int(data[position])
	// Śledzenie czy nastąpiło przekierowanie, aby uniknąć nieskończonych pętli
	jumped := false
	// Zapisz oryginalną pozycję dla zwrócenia poprawnego offsetu
	originalPosition := position
	// Ustawiamy maksymalną liczbę skoków, aby zapobiec atakom
	maxJumps := 10
	jumps := 0

	for length > 0 {
		// Sprawdź czy to wskaźnik kompresji (dwa najstarsze bity ustawione na 1)
		if (length & 0xC0) == 0xC0 {
			if jumps >= maxJumps {
				return name, position + 2
			}

			// To jest wskaźnik, oblicz offset
			if position+1 >= len(data) {
				return name, position + 1
			}

			pointer := ((length & 0x3F) << 8) | int(data[position+1])

			// Jeśli jest to pierwszy skok, zapisz następną pozycję
			if !jumped {
				position += 2
				originalPosition = position
			}

			// Ustaw nową pozycję na podstawie wskaźnika
			position = pointer
			jumped = true
			jumps++

			// Pobierz długość etykiety w nowej pozycji
			if position >= len(data) {
				return name, originalPosition
			}
			length = int(data[position])
			continue
		}

		position++

		if position+length > len(data) {
			break
		}

		label := string(data[position : position+length])
		name += label + "."

		position += length

		if position >= len(data) {
			break
		}

		length = int(data[position])
	}

	if length == 0 {
		position++
	}

	// Jeśli nastąpiło przekierowanie, zwróć oryginalną pozycję
	if jumped {
		return name, originalPosition
	}

	return name, position

}

func domainNameToBytes(domain string) []byte {
	var bytes []byte

	if len(domain) > 0 && domain[len(domain)-1] == '.' {
		domain = domain[:len(domain)-1]
	}

	if domain == "" {
		return []byte{0}
	}

	labels := strings.Split(domain, ".")
	for _, label := range labels {
		bytes = append(bytes, byte(len(label)))
		bytes = append(bytes, []byte(label)...)
	}

	bytes = append(bytes, 0)
	return bytes
}
