package main

import (
	"fmt"
	"linkontun/internal/byudp"
	"log"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

func main() {
	cmd := exec.Command("bash", "scripts/tun.sh")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	waterConfig := water.Config{
		DeviceType: water.TUN,
	}
	waterConfig.Name = "chain0"
	ifce, err := water.New(waterConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	go byudp.IN(ifce)
	go byudp.OUT(ifce)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	sig := <-sigs
	fmt.Println()
	fmt.Println(sig)
	done <- true

	<-done

	fmt.Println("exiting")
}
