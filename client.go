/*************************************************************************
   > File Name: client.go
   > Author: jige003
   > Created Time: Fri 01 Mar 2019 06:02:22 PM CST
************************************************************************/
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func sha256hex(b []byte) (hexstr string) {
	h := sha256.New()
	h.Write(b)
	hexstr = hex.EncodeToString(h.Sum(nil))
	return
}

func run_client(ip string) {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	tundevice, err := NewTun("tun0")
	if err != nil {
		log.Fatalln(err)
	}
	err = tundevice.Up("10.0.1.2", "10.0.1.1")
	if err != nil {
		log.Fatalln(err)
	}
	rb := make([]byte, BUFFER_SIZE)

	go func() {
		for {

			n, addr, err := c.ReadFrom(rb)
			if err != nil {
				log.Println(err)
			}

			if addr.String() == "" {
				continue
			}
			rm, err := icmp.ParseMessage(1, rb[:n])
			if err != nil {
				log.Println(err)
			}
			if rm.Type == ipv4.ICMPTypeEchoReply {
				echo_object := rm.Body.(*icmp.Echo)
				log.Printf("[ Client handle icmp] addr:%v type:%v icmplen:%v, tunlen:%v\n", addr, rb[0], n, len(echo_object.Data))
				n, err := tundevice.Write(echo_object.Data)
				if err != nil {
					log.Println(n, err)
				}
			}
		}
	}()

	buffer := make([]byte, BUFFER_SIZE)
	for {
		n, err := tundevice.Read(buffer)
		tunlen := n
		if err != nil {
			log.Fatal(err)
		}
		wm := icmp.Message{
			Type: ipv4.ICMPTypeEcho, Code: 0,
			Body: &icmp.Echo{
				ID: os.Getpid() & 0xffff, Seq: 1,
				Data: buffer[:n],
			},
		}
		q := net.ParseIP(ip)
		addr := &net.IPAddr{q, ""}
		wb, err := wm.Marshal(nil)
		if err != nil {
			log.Fatal(err)
		}
		n, err = c.WriteTo(wb, addr)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[ Client handle tun ] addr:%v icmplen:%v tunlen:%v\n", addr, len(wb), tunlen)
	}
}
