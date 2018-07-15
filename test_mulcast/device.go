package main

import (
	"fmt"
	"net"
)

func SendHello() {
	conn, err := net.Dial("udp", "239.255.255.250:3701")
	if err != nil {
		fmt.Println(err)
	}
	conn.Write([]byte("hello every one, I waked up!"))
	defer conn.Close()
}
func SendBye() {
	conn, err := net.Dial("udp", "239.255.255.250:3701")
	if err != nil {
		fmt.Println(err)
	}
	conn.Write([]byte("hello every one, I was asleep!"))
	defer conn.Close()

}
func main() {
	//SendHello()
	addr, err := net.ResolveUDPAddr("udp", "239.255.255.250:3701")
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("send hello over!")
	data := make([]byte, 1024)
	for {
		//如果第二参数为nil,它会使用系统指定多播接口
		listener, err := net.ListenMulticastUDP("udp", nil, addr)
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		listener.Close()

		fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
		conn, err := net.Dial("udp", "239.255.255.250:3701")
		_, err = conn.Write([]byte("I got it! My ID is 007!"))
		if err != nil {
			fmt.Println(err)
			return
		}
		conn.Close()
		//time.Sleep(1 * time.Second)
	}
	defer SendBye()
}
