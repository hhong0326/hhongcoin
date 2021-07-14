package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/hhong0326/hhongcoin/explorer"
	"github.com/hhong0326/hhongcoin/rest"
)

// CLI practice with flag package
// Best Framework is cobra

// rest := flag.NewFlagSet("rest", flag.ExitOnError)

// portFlag := rest.Int("port", 4000, "Sets the port of the server")

// switch os.Args[1] {
// case "explorer":
// 	fmt.Println("Start Explorer")
// case "rest":
// 	rest.Parse(os.Args[2:])
// default:
// 	usage()
// }

// if rest.Parsed() {

// 	fmt.Println(*portFlag)
// 	fmt.Println("Start Server")
// }

func usage() {

	fmt.Printf("Welcome to 홍코인\n\n")
	fmt.Printf("Please use the following flags: \n\n")
	fmt.Printf("-rport=4000: 	Set port of the rest server\n")
	fmt.Printf("-hport=5000: 	Set port of the html server\n")
	fmt.Printf("-mode=rest: 	Choose between 'html', 'rest' and 'both' for both\n")

	runtime.Goexit() // exec defer function before terminating
}

func Start() {

	if len(os.Args) == 1 {
		usage()
	}

	// name, default, description
	rport := flag.Int("rport", 4000, "Set port of the rest server")
	hport := flag.Int("hport", 5000, "Set port of the html server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*rport)
	case "html":
		explorer.Start(*hport)
	case "both":
		go rest.Start(*rport)
		explorer.Start(*hport)
	default:
		usage()
	}
}
