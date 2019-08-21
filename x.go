/*************************************************************************
   > File Name: main.go
   > Author: jige003
   > Created Time: Thu 28 Feb 2019 03:51:53 PM CST
************************************************************************/
package main

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func tttt() {
	ifce, err := NewTun("tun0")
	if err != nil {
		log.Fatal(err)
	}
	b := make([]byte, 1500)
	for {
		n, err := ifce.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		if n <= 0 {
			continue
		}
		log.Println(n)
		packet := gopacket.NewPacket(b[:n], layers.LinkTypeRaw, gopacket.Default)
		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		log.Println(ethernetLayer)
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			log.Println(ip.SrcIP.String(), ip.DstIP.String())
		}

	}
}
