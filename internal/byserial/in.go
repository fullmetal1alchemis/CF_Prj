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
		packetBuffer = append(packetBuffer, serialdata[:n]...)

		ok, pktSize, err := isPkt(&packetBuffer)

		if err != nil {
			fmt.Println
		}

		if ok {

			//fmt.Printf("%x\n", packetBuffer[:pktSize])
			_, err = tun.Write(packetBuffer[:pktSize])
			if err != nil {
				fmt.Println("Tun write error: ", err)
			}

			// delete writed pkt
			packetBuffer = packetBuffer[pktSize:]
		}

	}
}

func isPkt(buffer *[]byte) (ok bool, pktSize int, err error) {
	if len(packetBuffer) <= 40 {
		return false, 0, nil
	}

	ipVersion := int(packetBuffer[0] / 16)

	if ipVersion == 4 {
		packetSize := int(packetBuffer[3]) + 256*int(packetBuffer[2])
	} else if ipVersion == 6 {
		packetSize := int(packetBuffer[5]) + 256*int(packetBuffer[4]) + 40
	} else {
		return false, 0, errors.New("wrong IP version")
	}

	if len(packetBuffer) >= packetSize {
		fmt.Println("IP packet is ready")

		return true, packetSize, nil
	}
}
