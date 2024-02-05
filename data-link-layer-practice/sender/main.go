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
	copy(frame[0:6], f.Destination[0:6])
	copy(frame[6:12], f.Source[0:6])
	frame[12] = byte(f.EtherType >> 8)
	frame[13] = byte(f.EtherType & 0xff)
	copy(frame[14:], f.Payload)
	return frame
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <src device name> <dst device name>")
	}
	ifiSrc, err := net.InterfaceByName(os.Args[1])
	if err != nil {
		log.Fatalf("could not find source interface %s: %v", ifiSrc, err)
	}
	ifiDst, err := net.InterfaceByName(os.Args[2])
	if err != nil {
		log.Fatalf("could not fidnd destination interface %s: %v", ifiDst, err)
	}

	conn, err := raw.ListenPacket(ifiSrc, 0x0800, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	payload := []byte("Hi!")

	frame := &EthernetFrame{
		Destination: ifiDst.HardwareAddr,
		Source:      ifiSrc.HardwareAddr,
		EtherType:   0x0800,
		Payload:     payload,
	}
	log.Printf("Sending frame from %s to %s with payload: %s\n", frame.Source, frame.Destination, hex.EncodeToString(frame.Payload))
	if _, err = conn.WriteTo(frame.Marshal(), &raw.Addr{HardwareAddr: ifiDst.HardwareAddr}); err != nil {
		log.Fatal(err)
	}
}
