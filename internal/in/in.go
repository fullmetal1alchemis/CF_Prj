package in

import (
	"fmt"
	"net"

	"github.com/songgao/water"
)

//从socket读入数据然后写到网卡

func IN(wi *water.Interface, as int) {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(192, 168, 85, 101),
		Port: 30000,
	})
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()

	for {
		data := make([]byte, 2048)
		n, addr, err := listen.ReadFromUDP(data)
		if err != nil {
			fmt.Println("read udp failed, err:", err)
			continue
		}
		fmt.Printf("data:%v addr:%v count:%v\n", string(data[:n]), addr, n)

		_, err1 := wi.Write(data[:n])
		if err1 != nil {
			fmt.Println("write to tun failed, err:", err)
			continue
		}
	}

}
