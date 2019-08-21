/*************************************************************************
   > File Name: if.go
   > Author: jige003
   > Created Time: Thu 28 Feb 2019 01:53:21 PM CST
************************************************************************/
package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

type interfaces string

/* TUNSETIFF ifr flags */
const (
	IFF_TUN   = 0x0001
	IFF_TAP   = 0x0002
	IFF_NO_PI = 0x1000
)

type IfReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

type TunDevice struct {
	Name string
	F    *os.File
}

func (dev *TunDevice) Up(src_ip, dst_ip string) (err error) {
	//if err = exec.Command("ip", "link", "set", dev.Name, "up").Run(); err != nil {
	//	return
	//}
	//if err = exec.Command("ip", "addr", "add", src_ip, "dev", dev.Name).Run(); err != nil {
	//	return
	//}
	//in := dev.GetCIDR(src_ip)
	//if err = exec.Command("ip", "route", "add", in.String(), "via", src_ip, "dev", dev.Name).Run(); err != nil {
	//	return
	//}
	if err = exec.Command("ifconfig", dev.Name, src_ip, "netmask", "255.255.255.0", "up").Run(); err != nil {
		return
	}
	return
}

func (dev *TunDevice) GetCIDR(ipstr string) (in *net.IPNet) {
	mask := net.IPv4Mask(byte(255), byte(255), byte(255), byte(0))
	ip := net.ParseIP(ipstr).Mask(mask)
	in = &net.IPNet{ip, mask}
	return
}

func (dev *TunDevice) Write(b []byte) (int, error) {
	return dev.F.Write(b)
}

func (dev *TunDevice) Read(b []byte) (int, error) {
	return dev.F.Read(b)
}

//func (dev *TunDevice) WriteICMPPacket(b []byte) (int, error) {
//	wm := icmp.Message{
//		Type: 0, Code: 0,
//		Body: &icmp.Echo{
//			ID: os.Getpid() & 0xffff, Seq: 1,
//			Data: []byte("hello reply"),
//		},
//	}
//	dev.Write(wm.Marshal(nil))
//}

func NewInterface(name string, flags uint16) (device *TunDevice, err error) {
	var file *os.File
	file, err = os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return

	}
	ifr := IfReq{}
	ifr.Flags = flags
	copy(ifr.Name[:], name)
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&ifr))); errno != 0 {
		err = os.NewSyscallError("ioctl", errno)
		return
	}
	ifName := strings.Trim(string(ifr.Name[:]), "\x00")
	log.Println("ifName: ", ifName)
	device = &TunDevice{Name: ifName, F: file}
	return
}

func NewTun(name string) (tundevice *TunDevice, err error) {
	tundevice, err = NewInterface(name, IFF_TUN|IFF_NO_PI)
	if err != nil {
		return
	}
	return
}
