package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9998,
	})

	if err != nil {
		log.Printf("Listen failed %+v", err)
		return
	}
	log.Println("udp server starting ... ")
	for {
		var data [1024]byte
		n, addr, err := udpConn.ReadFromUDP(data[:])
		if err != nil {
			log.Printf("Read from udp server:%s failed,err:%s", addr, err)
			break
		}
		go func() {
			// 返回数据
			fmt.Printf("Addr:%s,data:%v count:%d \n", addr, string(data[:n]), n)

			//time.Sleep(time.Second * time.Duration(2))
			_, err := udpConn.WriteToUDP([]byte(fmt.Sprintf("ok -- recv:%s,response:%s",string(data[:n]),string(data[:n]))), addr)
			if err != nil {
				log.Printf("write to udp msg  failed! %+v", err)
			}
		}()
	}
}
