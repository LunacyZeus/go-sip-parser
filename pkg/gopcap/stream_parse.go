package gopcap

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sip-parser/pkg/siprocket"
	"time"
)

func streamParsePacket(src io.Reader, flipped bool, linkType Link) (msg siprocket.SipMsg, err error) {
	buffer := make([]byte, 16)
	read_count, err := src.Read(buffer)

	//sfmt.Println(buffer)

	if err != nil {
		log.Printf("Err->%v", err)
		return siprocket.SipMsg{}, err
	} else if read_count != 16 {
		log.Printf("InsufficientLength")
		return siprocket.SipMsg{}, InsufficientLength
	}

	// First is a pair of fields that build up the timestamp.
	//ts_seconds := getUint32(buffer[0:4], flipped)
	//ts_micros := getUint32(buffer[4:8], flipped)
	//Timestamp := time.Now()

	//log.Println(packet.Timestamp)
	// Next is the length of the data segment.
	IncludedLen := getUint32(buffer[8:12], flipped)

	// Then the original length of the packet.
	//ActualLen := getUint32(buffer[12:16], flipped)

	data := make([]byte, IncludedLen)
	//readlen, err := src.Read(data)

	/*
		if uint32(readlen) != pkt.IncludedLen {
			return UnexpectedEOF
		}

	*/

	d, err := parseLinkData(data, linkType)
	if d == nil {
		log.Println("err parseLinkData")
		return siprocket.SipMsg{}, errors.New("err parseLinkData")
	}

	td := d.LinkData().InternetData().TransportData()
	if td == nil {
		log.Println("unexpected transport data")
		return siprocket.SipMsg{}, errors.New("unexpected transport data")
	}
	fmt.Println(string(td))
	msg = siprocket.Parse(td, time.Now())
	fmt.Println(string(msg.To.Src))

	return msg, nil
}
