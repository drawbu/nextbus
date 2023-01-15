package types

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
