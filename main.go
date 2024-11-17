package main

import (
	"flag"
	"fmt"
	"os"
	"sip-parser/pkg/sip"
	"strings"
)

var (
	file   = flag.String("file", "", "The SIP pcap file that will be parsed")
	to     = flag.String("to", "", "SIP To: Header")
	from   = flag.String("from", "", "SIP From: Header")
	callid = flag.String("callid", "", "SIP Call-ID header")
	help   = flag.Bool("help", false, "Display usage help")
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s \n", os.Args[0])
		fmt.Println("If no `--to` and `--from` are specified then the program will output `To:` and `From:` from all SIP dialogs.")
		fmt.Println("Parameters:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *file == "" {
		errorOut("you need to specify the file with --file")
	}

	*from = strings.ToLower(*from)
	*to = strings.ToLower(*to)
	*callid = strings.ToLower(*callid)

	trace, err := sip.ParsePcapFile(*file)
	if err != nil {
		errorOut(fmt.Sprintf("cannot parse SIP trace: %s", err))
	}

	// Parse the the SIP data
	fp, err := sip.ParseSIPTrace(trace)
	if err != nil {
		errorOut(err.Error())
	}

	// Search the SIP packets for the filters
	sip.HandleSipPackets(fp)
}
