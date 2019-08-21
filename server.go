/*************************************************************************
   > File Name: server.go
   > Author: jige003
   > Created Time: Fri 01 Mar 2019 05:56:44 PM CST
************************************************************************/
package main

import (
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var client_ip string

const (
	BUFFER_SIZE = 1600
)

func run_server() {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)

	}
	defer c.Close()

	tundevice, err := NewTun("tun0")
	if err != nil {
		log.Fatalln(err)
	}
	err = tundevice.Up("10.0.1.1", "")
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
			client_ip = addr.String()
			rm, err := icmp.ParseMessage(1, rb[:n])
			if err != nil {
				log.Println(err)
			}
			if rm.Type == ipv4.ICMPTypeEcho {
				echo_object := rm.Body.(*icmp.Echo)
				log.Printf("[ Server handle icmp ] addr:%v icmplen:%v tunlen:%v \n", client_ip, n, len(echo_object.Data))
				n, err := tundevice.Write(echo_object.Data)
				if err != nil {
					log.Println(n, err)
				}
			}
		}
	}()

	buffer := make([]byte, BUFFER_SIZE)
	for {
		n, e := tundevice.Read(buffer)
		tunlen := n
		if e != nil {
			log.Println(n, e)
		}

		wm := icmp.Message{
			Type: ipv4.ICMPTypeEchoReply, Code: 0,
			Body: &icmp.Echo{
				ID: os.Getpid() & 0xffff, Seq: 1,
				Data: buffer[:n],
			},
		}
		if client_ip == "" {
			continue
		}
		q := net.ParseIP(client_ip)
		addr := &net.IPAddr{q, ""}
		wb, err := wm.Marshal(nil)
		if err != nil {
			log.Fatal(err)
		}
		n, err = c.WriteTo(wb, addr)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[ Server reply ] addr:%v tunlen:%v icmplen:%v\n", addr, tunlen, len(wb))

	}

}
