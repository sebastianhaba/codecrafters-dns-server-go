package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
	} else {
		fmt.Println("Current directory:", currentDir)
	}

	resolverFlag := flag.String("resolver", "", "DNS resolver address in the format IP:port")
	flag.Parse()

	if *resolverFlag == "" {
		fmt.Println("Error: --resolver flag is required")
		os.Exit(1)
	}

	fmt.Println("DNS server starting, forwarding to resolver:", *resolverFlag)

	resolverAddr, err := net.ResolveUDPAddr("udp", *resolverFlag)
	if err != nil {
		fmt.Println("Failed to resolve resolver address:", err)
		os.Exit(1)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		fmt.Println(hex.EncodeToString(buf[:size]))

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		dnsMessage := dns.NewMessage(buf[:size])
		fmt.Printf("Parsed message %s\n", dnsMessage)

		if len(dnsMessage.Question) > 1 {
			fmt.Printf("Multiple questions detected (%d), splitting requests\n", len(dnsMessage.Question))

			// Przygotuj odpowiedź dla klienta
			responseMessage := &dns.Message{
				Header:   dnsMessage.Header,
				Question: dnsMessage.Question,
				Answer:   []dns.ResourceRecord{},
			}
			responseMessage.Header.QR = 1 // To jest odpowiedź

			// Dla każdego pytania wyślij oddzielne zapytanie do resolvera
			for _, question := range dnsMessage.Question {
				// Utwórz nowe zapytanie z jednym pytaniem
				singleQuestion := &dns.Message{
					Header: dns.Header{
						ID:      dnsMessage.Header.ID,
						QR:      0, // To jest zapytanie
						OPCODE:  dnsMessage.Header.OPCODE,
						AA:      0,
						TC:      0,
						RD:      dnsMessage.Header.RD,
						RA:      0,
						Z:       0,
						RCODE:   0,
						QDCOUNT: 1,
						ANCOUNT: 0,
						NSCOUNT: 0,
						ARCOUNT: 0,
					},
					Question: []dns.Question{question},
				}

				// Wysyłanie zapytania do resolvera
				resolverResponse := forwardDNSQuery(singleQuestion.ToBytes(), resolverAddr)
				if resolverResponse != nil {
					// Parsuj odpowiedź
					resolverMessage := dns.NewMessage(resolverResponse)

					// Dodaj odpowiedzi do naszej odpowiedzi zbiorczej
					responseMessage.Answer = append(responseMessage.Answer, resolverMessage.Answer...)
				}
			}

			// Aktualizuj licznik odpowiedzi
			responseMessage.Header.ANCOUNT = uint16(len(responseMessage.Answer))

			// Wyślij zbiorczą odpowiedź z powrotem do klienta
			response := responseMessage.ToBytes()
			_, err = udpConn.WriteToUDP(response, source)
			if err != nil {
				fmt.Println("Failed to send response:", err)
			}
		} else {
			// Standardowy przypadek z jednym pytaniem
			// Przekaż zapytanie do resolvera
			resolverResponse := forwardDNSQuery(buf[:size], resolverAddr)

			if resolverResponse != nil {
				// Wyślij odpowiedź z powrotem do klienta
				_, err = udpConn.WriteToUDP(resolverResponse, source)
				if err != nil {
					fmt.Println("Failed to send response:", err)
				}
			} else {
				fmt.Println("No response from resolver")
			}
		}
	}
}

func forwardDNSQuery(query []byte, resolverAddr *net.UDPAddr) []byte {
	conn, err := net.DialUDP("udp", nil, resolverAddr)
	if err != nil {
		fmt.Println("Failed to connect to resolver:", err)
		return nil
	}
	defer conn.Close()

	_, err = conn.Write(query)
	if err != nil {
		fmt.Println("Failed to forward query to resolver:", err)
		return nil
	}

	buffer := make([]byte, 512)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	size, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			fmt.Println("Timeout waiting for response from resolver")
		} else {
			fmt.Println("Error receiving response from resolver:", err)
		}
		return nil
	}

	fmt.Printf("Received %d bytes response from resolver\n", size)
	return buffer[:size]
}
