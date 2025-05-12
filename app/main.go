package main

import (
	"encoding/hex"
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"net"
	"os"
)

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
	} else {
		fmt.Println("Current directory:", currentDir)
	}

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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

		// Create an empty response
		dnsMessage.Answer = make([]dns.ResourceRecord, len(dnsMessage.Question))

		for i, question := range dnsMessage.Question {
			dnsMessage.Answer[i] = dns.NewARecord(question.QNAME, net.ParseIP("127.0.0.1"), 3600)
		}

		response := dnsMessage.ToBytes()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
