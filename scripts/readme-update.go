package main

import (
	"main/tbm"
	"os"
)

func main() {
	file := `# nextbus
	
This project aims to provide a terminal-based realtime dashboard for monitoring
buses and tram in the city of Bordeaux, France, from the TBM transport company.
	
/!\ This project is not an official TBM project /!\
	
## How to install

` + "```" + `bash
$ git clone https://github/drawbu/nextbus.git
$ cd nextbus
$ go build -o nextbus main.go
` + "```" + `
	
## How to use
	
` + "```" + `bash
$ ./nextbus
` + tbm.GetHelpMessage() + "\n```" + `

## Examples

` + "```" + `bash
$ ./nextbus bus 10 Peixotto
`
	err, exampleOne := tbm.GetRealtimeBusArrival("Peixotto", "10")
	if err != nil {
		panic(err)
	}
	file += exampleOne + "```" + `

` + "```" + `bash
$ ./nextbus bus 21
`
	err, exampleTwo := tbm.GetStopList("21")
	if err != nil {
		panic(err)
	}
	file += exampleTwo + "```"
	err = os.WriteFile("README.md", []byte(file), 0644)
	if err != nil {
		panic(err)
	}
}
