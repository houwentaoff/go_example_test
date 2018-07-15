package main

/*
* 上位机运行
 */
import (
	"fmt"
	"net"
	"time"
)

func Discover() {
	conn, err := net.Dial("udp", "239.255.255.250:3701")
	if err != nil {
		fmt.Println(err)
	}
	conn.Write([]byte("probe discover: wake up! every device!"))
	defer conn.Close()
}
func RecvProbe() {
	data := make([]byte, 1024)
	mcaddr, _ := net.ResolveUDPAddr("udp4", "239.255.255.250:3701")
	conn, _ := net.ListenMulticastUDP("udp4", nil, mcaddr)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
	}
	conn.Close()
}
func main() {
	/* 3s */
	go func() {
		for {
			time.Sleep(3 * time.Second)
			Discover()
		}
	}()
	RecvProbe()
}
