package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/types"
	"net/http"
	"os"
	"strings"
)

const BaseUrl = "https://ws.infotbm.com/ws/1.0"

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

func getHelpMessage() string {
	return `Usage: nextbus [transport] [line] [stop]
Options and arguments (and corresponding environment variables):
-h, --help : print help (this message) and exit (also print this message if no
             argument is provided)
transport  : type of transport (bus, car, tram...)
line       : line number
stop       : stop name, optional, will print all stop in the line if missing`
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

func getRealTimeDataBuses(busName string, stop types.LineStop, directionId string) (err error, result []types.RealtimeStop) {
	var realtimePass types.RealtimePass
	stopId := strings.Split(stop.Id, ":")[3]
	url := fmt.Sprintf("%v/get-realtime-pass/%v/%v/route:TBC:%v", BaseUrl, stopId, busName, busName)
	err = getRequest(url, &realtimePass)
	if err == nil {
		result = realtimePass.Destinations[strings.Split(directionId, ":")[3]]
		return
	}

	// If the request or parsing failed, try again with the opposite direction
	err = getRequest(url+"_R", &realtimePass)
	if err == nil {
		result = realtimePass.Destinations[strings.Split(directionId, ":")[3]]
	}
	return
}

func getRealtimeBusArrival(route types.LineRoute, stopName string, line string) (err error) {
	err, stop := getStop(route.StopPoints, stopName)
	if err != nil {
		fmt.Println("stop not found")
		err = nil
		return
	}

	direction := route.StopPoints[len(route.StopPoints)-1]
	err, realTimeDataBuses := getRealTimeDataBuses(line, stop, direction.Id)
	if err != nil {
		return
	}
	fmt.Printf("Bus %v, %v, direction %v\n", line, stop.Name, direction.Name)
	for _, e := range realTimeDataBuses {
		fmt.Printf("- %v\n", e.WaitTimeText)
	}
	return
}

func main() {
	args := types.Args{}
	err := args.GetArgs(os.Args)
	if err != nil {
		panic(err)
	}
	if args.Help {
		fmt.Printf("%v", getHelpMessage())
		return
	}

	switch args.TransportType {
	case "bus":
		var line types.Line
		err := getRequest(fmt.Sprintf("%v/network/line-informations/%v", BaseUrl, args.Line), &line)
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
		for _, route := range line.Routes {
			err = getRealtimeBusArrival(route, args.Stop, args.Line)
			if err != nil {
				panic(err)
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
