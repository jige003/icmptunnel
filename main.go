/*************************************************************************
   > File Name: main.go
   > Author: jige003
   > Created Time: Fri 01 Mar 2019 05:12:55 PM CST
************************************************************************/
package main

import (
	"flag"
	"os"
)

var (
	is_server   bool   = false
	server_addr string = ""
)

func init() {
	flag.BoolVar(&is_server, "s", false, "server mode")
	flag.StringVar(&server_addr, "c", "10.0.1.1", "server addr")
	flag.Parse()

}

func main() {
	if len(os.Args) < 2 {
		flag.Usage()
	}
	if is_server {
		run_server()
	} else {
		run_client(server_addr)
	}
}
