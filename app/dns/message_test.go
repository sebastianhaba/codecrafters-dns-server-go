package dns

import (
	"encoding/hex"
	"testing"
)

func TestNewMessage(t *testing.T) {
	hexStr := "04d2010000010000000000000c636f6465637261667465727302696f0000010001"

	messageBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatalf("Błąd konwersji hex do bajtów: %v", err)
	}

	message := NewMessage(messageBytes)

	if message.Header.ID != 1234 { // 0x04d2 = 1234 w systemie dziesiętnym
		t.Errorf("want ID = 1234, got %d", message.Header.ID)
	}

	if message.Header.QR != 0 {
		t.Errorf("want QR = 0 (zapytanie), otrzymano %d", message.Header.QR)
	}

	if message.Header.QDCOUNT != 1 {
		t.Errorf("want QDCOUNT = 1, got %d", message.Header.QDCOUNT)
	}

	if message.Header.ANCOUNT != 0 {
		t.Errorf("want ANCOUNT = 0, got %d", message.Header.ANCOUNT)
	}

	if message.Header.NSCOUNT != 0 {
		t.Errorf("want NSCOUNT = 0, got %d", message.Header.NSCOUNT)
	}

	if message.Header.ARCOUNT != 0 {
		t.Errorf("want ARCOUNT = 0, got %d", message.Header.ARCOUNT)
	}

	// Weryfikacja pytania
	if len(message.Question) != 1 {
		t.Fatalf("want 1 pytanie, got %d", len(message.Question))
	}

	expectedQName := "codecrafters.io."
	if message.Question[0].QNAME != expectedQName {
		t.Errorf("want QNAME = %s, got %s", expectedQName, message.Question[0].QNAME)
	}

	if message.Question[0].QTYPE != 1 { // 1 = A record
		t.Errorf("want QTYPE = 1 (A), got %d", message.Question[0].QTYPE)
	}

	if message.Question[0].QCLASS != 1 { // 1 = IN (Internet)
		t.Errorf("want QCLASS = 1 (IN), got %d", message.Question[0].QCLASS)
	}
}
