package types

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
