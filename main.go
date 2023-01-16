package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/types"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getRequest(url string, response any) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Response not OK: reponse %s", resp.Status))
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	return
}

func printHelp() {
	fmt.Println(`Usage: nextbus [transport] [line] [stop]
Options and arguments (and corresponding environment variables):
-h, --help : print help (this message) and exit (also print this message if no
             argument is provided)
transport  : type of transport (bus, car, tram...)
line       : line number
stop       : stop name`)
}

func getStop(line []types.LineStop, stopName string) (err error, stop types.LineStop) {
	for _, s := range line {
		if strings.Contains(strings.ToLower(s.Name), strings.ToLower(stopName)) {
			stop = s
			return
		}
	}
	err = errors.New("stop not found")
	return
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		printHelp()
		return
	}
	if len(args) < 3 {
		fmt.Println("Not enough argument provided, please refer to the help: nextbus -h")
		return
	}
	if len(args) > 3 {
		fmt.Println("Too many argument provided, please refer to the help: nextbus -h")
		return
	}
	const baseUrl = "https://ws.infotbm.com/ws/1.0"

	switch strings.ToLower(args[0]) {
	case "bus":
		var line types.Line
		err := getRequest(fmt.Sprintf("%v/network/line-informations/%v", baseUrl, args[1]), &line)
		if err != nil {
			panic(err)
		}

		// A stop can have multiple direction, so we loop over all of them
		for _, route := range line.Routes {
			// Find the bus stop, if the stop isn't found, stop here
			err, stop := getStop(route.StopPoints, args[2])
			if stop.Name == "" {
				fmt.Println("Stop not found")
				continue
			}

			// Find direction
			direction := route.StopPoints[len(route.StopPoints)-1]
			fmt.Printf("Bus %v, %v, direction %v\n", args[1], stop.Name, direction.Name)

			// Try to get realtime data, if it doesn't work try to get the
			// opposite direction, if it doesn't work stop here
			var realtimePass types.RealtimePass
			stopId, err := strconv.Atoi(strings.Split(stop.Id, ":")[3])
			url := fmt.Sprintf("%v/get-realtime-pass/%v/%v/route:TBC:%v", baseUrl, stopId, args[1], args[1])
			err = getRequest(url, &realtimePass)
			if err != nil {
				err = getRequest(url+"_R", &realtimePass)
				if err != nil {
					panic(err)
				}
			}

			// Print the next buses
			for _, e := range realtimePass.Destinations[strings.Split(direction.Id, ":")[3]] {
				fmt.Printf("- %v\n", e.WaitTimeText)
			}
		}
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
