//go:build ignore

package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func tcpScanByGoroutine(ip string, portStart int, portEnd int) {
	start := time.Now()

	isok := verifyParam(ip, portStart, portEnd)
	if isok == false {
		fmt.Printf("[Exit]\n")
		return
	}

	//并发控制：最大同时 N 个拨测
	const maxConcurrency = 200
	sem := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup
	var mu sync.Mutex
	openPorts := make([]int, 0, 64)

	timeout := 300 * time.Millisecond

	for i := portStart; i <= portEnd; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func (port int)  {
			defer wg.Done()

			defer func(){ <-sem}()

			address := fmt.Sprintf("%s:%d", ip, port)
			conn, err := net.DialTimeout("tcp", address, timeout)

			if err != nil {
				//fmt.Printf("[info] %s Closed\n",address)
				return
			}
			_ = conn.Close()
			mu.Lock()
			openPorts = append(openPorts, port)
			mu.Unlock()
			fmt.Printf("[info] %s Open\n", address)
		}(i)
	}
	wg.Wait()

	sort.Ints(openPorts)
	cost := time.Since(start)
	fmt.Printf("[tcpScan] finished, open ports: %v\n", openPorts)
	fmt.Printf("[tcpScan] cost %v\n", cost)
}

func verifyParam(ip string, portStart int, portEnd int) bool {
	netip := net.ParseIP(ip)
	if netip == nil {
		fmt.Println("[Error] invalid ip, type is must net.ip")
		return false
	}
	fmt.Printf("[info] ip=%s | ip type is: %T | ip length: %d bytes\n", netip, netip, len(netip))

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

func main (){
	tcpScanByGoroutine("127.0.0.1", 1, 65535)
}