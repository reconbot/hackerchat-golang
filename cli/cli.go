package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type SendingMessage struct {
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

type ReceivedMessage struct {
	Text      string
	Timestamp int64
	Ip        string
}

type UDPMessage struct {
	data []byte
	ip   string
}

func Ux(rx chan ReceivedMessage, tx chan SendingMessage, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			if text != "" {
				fmt.Println("Ux sending ", text)
				tx <- SendingMessage{
					Text:      text,
					Timestamp: time.Now().Unix(),
				}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		} else {
			fmt.Fprintln(os.Stderr, "reading standard input: had some unknown reason it stopped")
		}
		close(tx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for message := range rx {
			fmt.Printf("UX Received: %v %v %v\n", message.Timestamp, message.Text, message.Ip)
		}
		wg.Done()
	}()
}

func Decoder(rx chan UDPMessage, tx chan ReceivedMessage, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for data := range rx {
			var decoded SendingMessage
			err := json.Unmarshal(data.data, &decoded)
			if err != nil {
				fmt.Println("Error parsing message from", data.ip, err)
				continue
			}
			tx <- ReceivedMessage{
				Ip:        data.ip,
				Text:      decoded.Text,
				Timestamp: decoded.Timestamp,
			}
		}
		close(tx)
		wg.Done()
	}()
}

func Encoder(rx chan SendingMessage, tx chan []byte, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for message := range rx {
			b, err := json.Marshal(message)
			if err != nil {
				panic(err)
			}
			fmt.Println("Encoder encoded ", string(b))
			tx <- b
		}
		close(tx)
		wg.Done()
	}()
}

func Network(rx chan []byte, tx chan UDPMessage, wg *sync.WaitGroup) {
	broadcast_address, err := net.ResolveUDPAddr("udp", "255.255.255.255:31337")
	if err != nil {
		panic(fmt.Errorf("cannot resolve udp address or something %v", err))
	}

	receive_address, err := net.ResolveUDPAddr("udp", ":31337")
	if err != nil {
		panic(fmt.Errorf("cannot resolve udp address or something %v", err))
	}

	udp, err := net.ListenUDP("udp", receive_address)
	if err != nil {
		panic(fmt.Errorf("cannot open udp port %v", err))
	}

	wg.Add(1)
	go func() {
		for {
			buf := make([]byte, 1500)
			read, address, err := udp.ReadFromUDP(buf)
			if err != nil {
				fmt.Println(fmt.Errorf("cannot read udp %v", err))
				break
			}
			fmt.Println("got a message from the network!", string(buf[:read]))
			tx <- UDPMessage{data: buf[:read], ip: address.IP.To16().String()}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for packet := range rx {
			fmt.Println("Network sending", string(packet), broadcast_address)
			count, err := udp.WriteTo(packet, broadcast_address)
			if err != nil {
				panic(fmt.Errorf("cannot write udp %v", err))
			}
			if count != len(packet) {
				panic(fmt.Errorf("cannot write all of a udp packet %v", err))
			}
		}
		wg.Done()
	}()
}

func main() {
	var wg sync.WaitGroup
	uxTx := make(chan SendingMessage)
	decoderTx := make(chan ReceivedMessage)
	encoderTx := make(chan []byte)
	networkTx := make(chan UDPMessage)

	Ux(decoderTx, uxTx, &wg)
	Decoder(networkTx, decoderTx, &wg)
	Encoder(uxTx, encoderTx, &wg)
	Network(encoderTx, networkTx, &wg)

	wg.Wait()
}
