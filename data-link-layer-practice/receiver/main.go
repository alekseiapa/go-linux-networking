package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
)

func main() {

	if len(os.Args) != 3 {
		log.Fatalf("Usage: go run main.go <src device name> <dst device name>")
	}
	providedDeviceName, srcDeviceName := os.Args[1], os.Args[2]

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	var deviceFound bool = false
	for _, device := range devices {
		if device.Name == providedDeviceName {
			deviceFound = true
			break
		}
	}
	if !deviceFound {
		log.Fatalf("Device '%s' not found", providedDeviceName)
	}

	handle, err := pcap.OpenLive(providedDeviceName, 1500, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethernetLayer == nil {
			continue
		}
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)

		if ethernetPacket.SrcMAC.String() != srcDeviceName {
			continue
		}
		fmt.Println("Reading packet...")

		fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
		fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
		fmt.Println("Payload: ", string(ethernetPacket.Payload))
	}

}
