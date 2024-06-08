package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func sendWakeOnLAN(mac string) error {
	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9")
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	macAddr, err := net.ParseMAC(mac)
	if err != nil {
		return err
	}
	packet := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	for i := 0; i < 16; i++ {
		packet = append(packet, macAddr...)
	}

	_, err = conn.Write(packet)
	return err
}

func main() {
	socketPath := "/tmp/wake-by-a-docker-container.sock"

	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		listener.Close()
		os.Remove(socketPath)
		os.Exit(1)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()

		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				break
			}

			data := strings.TrimSpace(string(buffer[:n]))

			_, err = net.ParseMAC(data)
			if err != nil {
				fmt.Printf("Invalid MAC address received : %s\n", data)
			} else {
				fmt.Printf("Valid MAC address received : %s\n", data)
				if err := sendWakeOnLAN(data); err != nil {
					fmt.Printf("Error sending WoL packet : %v\n", err)
				} else {
					fmt.Println("WoL package sent successfully!")
				}
			}
		}
	}
}
