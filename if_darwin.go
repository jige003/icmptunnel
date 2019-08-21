/*************************************************************************
   > File Name: if_darwin.go
   > Author: jige003
   > Created Time: Tue 05 Mar 2019 04:24:39 PM CST
************************************************************************/
// +build darwin
package main

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
)

// maximum devices supported by driver
const maxDevices = 16

var (
	ErrBusy        = errors.New("device is already in use")
	ErrNotReady    = errors.New("device is not ready")
	ErrExhausted   = errors.New("no devices are available")
	ErrUnsupported = errors.New("device is unsupported on this platform")
)

type TunDevice struct {
	Name string
	F    *os.File
}

func (dev *TunDevice) Up(src_ip, dst_ip string) (err error) {
	if err = exec.Command("ifconfig", dev.Name, src_ip, dst_ip, "up").Run(); err != nil {
		return
	}
	return
}

func (dev *TunDevice) Write(b []byte) (int, error) {
	return dev.F.Write(b)

}

func (dev *TunDevice) Read(b []byte) (int, error) {
	return dev.F.Read(b)

}

// return true if read error is result of device not being ready
func isNotReady(err error) bool {
	if perr, ok := err.(*os.PathError); ok {
		if code, ok := perr.Err.(syscall.Errno); ok {
			if code == 0x05 {
				return true
			}
		}
	}
	return false
}

// return true if file error is result of device already being used
func isBusy(err error) bool {
	if perr, ok := err.(*os.PathError); ok {
		if code, ok := perr.Err.(syscall.Errno); ok {
			if code == 0x10 || code == 0x11 { // device busy || exclusive lock
				return true
			}
		}
	}
	return false
}

func NewTun(name string) (tun *TunDevice, err error) {
	var file *os.File
	file, err = os.OpenFile("/dev/"+name, os.O_EXCL|os.O_RDWR, 0)
	if isBusy(err) {
		return nil, ErrBusy

	} else if err != nil {
		return

	}
	tun = &TunDevice{Name: name, F: file}
	return
}
