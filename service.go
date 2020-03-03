package main

import (
	"log"
	"sync"
	"io/ioutil"
	"github.com/nats-io/nats.go"
	"encoding/binary"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	// [begin subscribe_queue]
	nc, err := nats.Connect("127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Use a WaitGroup to wait for 1st message to arrive
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Create a queue subscription on "video" with queue name "service"
	if _, err := nc.QueueSubscribe("video", "service", func(m *nats.Msg) {
		log.Printf("RECIEVED FILE PATH %v", string(m.Data))

		dat, err := ioutil.ReadFile(string(m.Data))
		if err != nil {
			m.Respond([]byte("READ FILE ERROR "))
			return
		}

		var ftypBox []byte
		var moovBox []byte
		i := uint32(0)
		fileSize := uint32(len(dat))

		for {
			size := binary.BigEndian.Uint32(dat[i:i+4])
			boxType := string(dat[i+4:i+8])


			if boxType == "ftyp" {
				ftypBox = dat[i:i+size]
			}
			if boxType == "moov" {
				moovBox = dat[i:i+size]
			}

			if len(boxType) > 0 && len(moovBox) > 0 {
				
				initSegPath := "./tmp/initseg"
				err := ioutil.WriteFile(initSegPath, append(ftypBox[:], moovBox[:]...), 0644)
				if err != nil {
					log.Fatal(err)
					m.Respond([]byte("WRITE FILE FAILED"))
				} else {
					m.Respond([]byte(initSegPath))
					wg.Done()
				}

				break
			}


			i += size
			if i >=  fileSize{
				m.Respond([]byte("INIT SEGMENT NOT FOUND"))
				break
			}
		}
	}); err != nil {
		log.Fatal(err)
	}

	// Wait for messages to come in
	wg.Wait()
}