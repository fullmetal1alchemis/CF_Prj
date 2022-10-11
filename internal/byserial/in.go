package byserial

import (
	"fmt"
	"log"

	"github.com/songgao/water"
	"go.bug.st/serial"
)

func IN(tun *water.Interface) {

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	device, err := serial.Open("/dev/ttyS1", mode)
	if err != nil {
		log.Fatal(err)
	}

	//packetBuffer := make([]byte, 2048)
	var packetBuffer []byte
	var testSize chan struct{}

	go getSerial(&packetBuffer, device, testSize)

	for {
		if len(packetBuffer) > 40 {
			//fmt.Println("BUFFER:")
			//fmt.Printf("% x\n", packetBuffer)
			ipVersion := int(packetBuffer[0] / 16)
			if ipVersion == 4 {
				packetSize := int(packetBuffer[3]) + 256*int(packetBuffer[2])
				if len(packetBuffer) > packetSize {
					packet := packetBuffer[:packetSize]
					fmt.Println("IPV4 Packet Recived:")
					fmt.Printf("% x\n", packet)
					packetBuffer = packetBuffer[packetSize:]
					_, err := tun.Write(packet)
					if err != nil {
						fmt.Println("Write to tun failed, err:", err)
						continue
					}
				}
			} else if ipVersion == 6 {
				packetSize := int(packetBuffer[5]) + 256*int(packetBuffer[4]) + 40
				if len(packetBuffer) > packetSize {
					packet := packetBuffer[:packetSize]
					fmt.Println("IPV6 Packet Recived:")
					fmt.Printf("% x\n", packet)
					packetBuffer = packetBuffer[packetSize:]
					_, err := tun.Write(packet)
					if err != nil {
						fmt.Println("Write to tun failed, err:", err)
						continue
					}
				}
			} else {
				packetBuffer = packetBuffer[:0]
			}
		}
	}

}

func getSerial(buffer *[]byte, device serial.Port, ch chan struct{}) {
	serialdata := make([]byte, 512)
	for {
		n, err := device.Read(serialdata)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n <= 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Println("Remote serial data received:")
		fmt.Printf("% x\n", serialdata[:n])
		*buffer = append(*buffer, serialdata[:n]...)
	}
}
