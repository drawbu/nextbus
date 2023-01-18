package main

import (
	"fmt"
	"main/tbm"
	"main/types"
	"os"
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
		if args.Stop == "" {
			err, result := tbm.GetStopList(args.Line)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%v", result)
		}

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
