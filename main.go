package main

import (
	"fmt"
	"main/tbm"
	"main/types"
	"os"
	"strings"
)

func main() {
	args := types.Args{}
	err := args.GetArgs(os.Args)
	if err != nil {
		panic(err)
	}
	if args.Help {
		fmt.Printf("%v", tbm.GetHelpMessage())
		return
	}

	switch args.TransportType {
	case "bus":
		var line types.Line
		err := tbm.GetRequest(fmt.Sprintf("%v/network/line-informations/%v", tbm.BaseUrl, args.Line), &line)
		if err != nil {
			panic(err)
		}

		// List all the stops
		if args.Stop == "" {
			fmt.Printf("Bus %v\n", line.Name)
			for _, e := range line.Routes {
				fmt.Printf("\n%v\n", strings.ToUpper(e.Name))
				for _, s := range e.StopPoints {
					fmt.Printf("- %v\n", s.Name)
				}
			}
			return
		}

		// Get the arrival time for a specific stop
		err, result := tbm.GetRealtimeBusArrival(args.Stop, args.Line)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v", result)
		break
	case "tram":
		fmt.Println("Not implemented yet")
		break
	case "car":
		fmt.Println("Not implemented yet")
		break
	default:
		fmt.Println("Unknown transport type, please refer to the help: nextbus -h")
		return
	}
}
