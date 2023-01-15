package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Line struct {
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	Code   string      `json:"code"`
	Type   string      `json:"type"`
	Routes []LineRoute `json:"routes"`
}

type LineRoute struct {
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Start      string     `json:"start"`
	End        string     `json:"end"`
	StopPoints []LineStop `json:"stopPoints"`
}

type LineStop struct {
	Id                    string `json:"id"`
	Name                  string `json:"name"`
	FullLabel             string `json:"fullLabel"`
	Latitude              string `json:"latitude"`
	Longitude             string `json:"longitude"`
	ExternalCode          string `json:"externalCode"`
	City                  string `json:"city"`
	HasWheelchairBoarding bool   `json:"hasWheelchairBoarding"`
	StopAreaId            string `json:"stopAreaId"`
	PartialStop           bool   `json:"partialStop"`
}

type RealtimeStop struct {
	VehicleLattitude         float64 `json:"vehicle_lattitude"`
	VehicleLongitude         float64 `json:"vehicle_longitude"`
	WaitTimeText             string  `json:"waittime_text"`
	TripId                   string  `json:"trip_id"`
	ScheduleId               string  `json:"schedule_id"`
	DestinationId            string  `json:"destination_id"`
	DestinationName          string  `json:"destination_name"`
	Departure                string  `json:"departure"`
	DepartureCommande        string  `json:"departure_commande"`
	DepartureTheorique       string  `json:"departure_theorique"`
	Arrival                  string  `json:"arrival"`
	ArrivalCommande          string  `json:"arrival_commande"`
	ArrivalTheorique         string  `json:"arrival_theorique"`
	Comment                  string  `json:"comment"`
	Realtime                 string  `json:"realtime"`
	WaitTime                 string  `json:"waittime"`
	UpdatedAt                string  `json:"updated_at"`
	VehicleId                string  `json:"vehicle_id"`
	VehiclePositionUpdatedAt string  `json:"vehicle_position_updated_at"`
	Origin                   string  `json:"origin"`
}

type RealtimePass struct {
	Destinations map[string][]RealtimeStop `json:"destinations"`
}

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

func getStop(line []LineStop, stopName string) (err error, stop LineStop) {
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
		var line Line
		err := getRequest(fmt.Sprintf("%v/network/line-informations/%v", baseUrl, args[1]), &line)
		if err != nil {
			panic(err)
		}

		// A stop can have multiple direction, so we print all of them
		for _, route := range line.Routes {
			// Find the bus stop
			err, stop := getStop(route.StopPoints, args[2])
			// If stop not found, stop here
			if stop.Name == "" {
				fmt.Println("Stop not found")
				continue
			}

			// Find direction
			direction := route.StopPoints[len(route.StopPoints)-1]
			fmt.Printf("Bus %v, %v, direction %v\n", args[1], stop.Name, direction.Name)

			// Try to get realtime data, if it doesn't work, trying to get the
			// opposite direction, if it doesn't work, stop here
			var realtimePass RealtimePass
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
