package main

import (
	"fmt"
	"net"
	"time"
	"unsafe"
)

func tcpScan(ip string, portStart int, portEnd int) {
	start := time.Now()

	isok := verifyParam(ip, portStart, portEnd)
	if isok == false {
		fmt.Printf("[Exit]\n")
		return
	}

	for i := portStart; i <= portEnd; i++ {
		address := fmt.Sprintf("%s:%d", ip, i)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("[info] %s Close \n", address)
			continue
		}
		conn.Close()
		fmt.Printf("[info] %s Open \n", address)
	}

	cost := time.Since(start)
	fmt.Printf("[tcpScan] cost %s \n", cost)
	
}

func verifyParam(ip string, portStart int, portEnd int) bool {
	netip := net.ParseIP(ip)
	if netip == nil {
		fmt.Println("[Error] invalid ip, type is must net.ip")
		return false
	}
	fmt.Printf("[info] ip=%s | ip type is: %T | ip size is: %d \n", netip, netip, unsafe.Sizeof(netip))

	if portStart < 1 || portEnd > 65535 || portStart > 65535 || portEnd	< 1{
		fmt.Println("[Error] port is must in the range of 1~65535")
		return false
	}
	if portStart > portEnd {
		fmt.Println("[Error] portStart must be <= portEnd")
	}
	fmt.Printf("[info] port start:%d end:%d \n",portStart, portEnd)

	return true
}

func main() {
	tcpScan("127.0.0.1", 1, 65535)
}