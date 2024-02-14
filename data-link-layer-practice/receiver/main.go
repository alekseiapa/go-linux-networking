package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: go run main.go <receiver MAC address>")
	}
	providedMACAddress := os.Args[1]

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	var interfaceName string
	for _, iface := range interfaces {
		log.Printf("Checking interface: %s with MAC address %s\n", iface.Name, iface.HardwareAddr.String())
		if iface.HardwareAddr.String() == providedMACAddress {
			interfaceName = iface.Name
			break
		}
	}

	if interfaceName == "" {
		log.Fatalf("No network interface found with MAC address '%s'", providedMACAddress)
	}

	handle, err := pcap.OpenLive(interfaceName, 1500, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	log.Printf("Listening for frames on interface %s ...\n", interfaceName)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {

		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethernetLayer == nil {
			continue
		}
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)

		if ethernetPacket.DstMAC.String() != providedMACAddress {
			continue
		}
		log.Println("Received ethernet packet...")
		log.Println("Reading packet...")

		log.Println("Source MAC: ", ethernetPacket.SrcMAC.String())
		log.Println("Destination MAC: ", ethernetPacket.DstMAC.String())
		log.Println("Payload: ", string(ethernetPacket.Payload))
	}
}
