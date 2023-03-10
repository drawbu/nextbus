package tbm

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/types"
	"net/http"
	"strings"
)

const BaseUrl = "https://ws.infotbm.com/ws/1.0"

func GetRequest(url string, response any) (err error) {
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

func GetHelpMessage() string {
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

func GetRealTimeDataBuses(busName string, stop types.LineStop, directionId string) (err error, result []types.RealtimeStop) {
	var realtimePass types.RealtimePass
	stopId := strings.Split(stop.Id, ":")[3]
	url := fmt.Sprintf("%v/get-realtime-pass/%v/%v/route:TBC:%v", BaseUrl, stopId, busName, busName)
	err = GetRequest(url, &realtimePass)
	if err == nil {
		result = realtimePass.Destinations[strings.Split(directionId, ":")[3]]
		return
	}

	// If the request or parsing failed, try again with the opposite direction
	err = GetRequest(url+"_R", &realtimePass)
	if err == nil {
		result = realtimePass.Destinations[strings.Split(directionId, ":")[3]]
	}
	return
}

func GetBusLine(line string) (err error, result types.Line) {
	err = GetRequest(fmt.Sprintf("%v/network/line-informations/%v", BaseUrl, line), &result)
	if err != nil {
		return
	}
	return
}

func GetRealtimeBusArrival(stopName string, line string) (err error, result string) {
	err, route := GetBusLine(line)
	if err != nil {
		return
	}
	for _, route := range route.Routes {
		err, stop := getStop(route.StopPoints, stopName)
		if err != nil {
			result += "Stop not found\n"
			err = nil
			continue
		}

		direction := route.StopPoints[len(route.StopPoints)-1]
		err, realTimeDataBuses := GetRealTimeDataBuses(line, stop, direction.Id)
		if err != nil {
			result += "No realtime data available\n"
			continue
		}
		result += fmt.Sprintf("Bus %v, %v, direction %v\n", line, stop.Name, direction.Name)
		for _, e := range realTimeDataBuses {
			result += fmt.Sprintf("- %v\n", e.WaitTimeText)
		}
	}
	return
}

func GetStopList(lineName string) (err error, result string) {
	err, line := GetBusLine(lineName)
	if err != nil {
		return
	}

	result += fmt.Sprintf("Bus %v\n", line.Name)
	for _, e := range line.Routes {
		result += fmt.Sprintf("\n%v\n", strings.ToUpper(e.Name))
		for _, s := range e.StopPoints {
			result += fmt.Sprintf("- %v\n", s.Name)
		}
	}
	return
}
