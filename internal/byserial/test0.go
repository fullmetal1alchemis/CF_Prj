package byserial

import (
	"fmt"
	"log"

	"go.bug.st/serial"
)

func CountBytes() {

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
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}

		fmt.Println(n, "bytes received.")
	}
}
