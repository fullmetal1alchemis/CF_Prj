package byserial

import (
	"fmt"
	"log"

	"github.com/songgao/water"
	"go.bug.st/serial"
)

func IN(wi *water.Interface) {

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	device, err := serial.Open("/dev/ttyS1", mode)
	if err != nil {
		log.Fatal(err)
	}

	serialdata := make([]byte, 1024)
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

		_, err1 := wi.Write(serialdata[:n])
		if err1 != nil {
			fmt.Println("Write to tun failed, err:", err1)
			continue
		}
	}

}
