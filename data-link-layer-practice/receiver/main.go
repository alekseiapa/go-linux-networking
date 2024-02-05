package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("Usage: go run main.go <receiver device name>")
	}
	providedDeviceName := os.Args[1]

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
	listenInterface, _ := net.InterfaceByName(providedDeviceName)
	listenInterfaceMACAddress := listenInterface.HardwareAddr.String()
	handle, err := pcap.OpenLive(providedDeviceName, 1500, false, pcap.BlockForever)
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

		if ethernetPacket.DstMAC.String() != listenInterfaceMACAddress {
			continue
		}
		fmt.Println("Reading packet...")

		fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
		fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
		fmt.Println("Payload: ", string(ethernetPacket.Payload))
	}

}
