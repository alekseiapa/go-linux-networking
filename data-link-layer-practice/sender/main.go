package main

import (
	"encoding/hex"
	"github.com/mdlayher/raw"
	"log"
	"net"
	"os"
)

type EthernetFrame struct {
	Destination net.HardwareAddr
	Source      net.HardwareAddr
	EtherType   uint16
	Payload     []byte
}

func (f *EthernetFrame) Marshal() []byte {
	// Ethernet frame format:
	// | Destination MAC (6 bytes) | Source MAC (6 bytes) | EtherType (2 bytes) | Payload (variable) |
	frame := make([]byte, 14+len(f.Payload))
	copy(frame[0:6], f.Destination)
	copy(frame[6:12], f.Source)
	frame[12] = byte(f.EtherType >> 8)
	frame[13] = byte(f.EtherType & 0xff)
	copy(frame[14:], f.Payload)
	return frame
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <src MAC address> <dst MAC address>")
	}
	srcMAC, err := net.ParseMAC(os.Args[1])
	if err != nil {
		log.Fatalf("could not parse source MAC address %s: %v", os.Args[1], err)
	}
	dstMAC, err := net.ParseMAC(os.Args[2])
	if err != nil {
		log.Fatalf("could not parse destination MAC address %s: %v", os.Args[2], err)
	}

	// Select a network interface to send the frame
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Available network interfaces: ", interfaces)
	var ifi *net.Interface
	for _, iface := range interfaces {
		log.Printf("Checking interface: %s with MAC address %s\n", iface.Name, iface.HardwareAddr.String())
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			if iface.HardwareAddr.String() == srcMAC.String() {
				ifi = &iface
				log.Printf("Selected network interface: %s\n", iface.Name)
				break
			}
		}
	}
	if ifi == nil {
		log.Fatal("no suitable interface found")
	}

	conn, err := raw.ListenPacket(ifi, 0x0800, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	payload := []byte("Hi!")

	frame := &EthernetFrame{
		Destination: dstMAC,
		Source:      srcMAC,
		EtherType:   0x0800,
		Payload:     payload,
	}

	log.Printf("Sending frame from %s to %s with payload: %s\n", frame.Source, frame.Destination, hex.EncodeToString(frame.Payload))
	if _, err = conn.WriteTo(frame.Marshal(), &raw.Addr{HardwareAddr: dstMAC}); err != nil {
		log.Fatal(err)
	}
}
