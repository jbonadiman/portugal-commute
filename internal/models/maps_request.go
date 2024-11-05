package models

import (
	"encoding/json"
	"time"
)

type waypoint struct {
	Address string `json:"address"`
}

type matrix struct {
	Waypoint waypoint `json:"waypoint"`
}

type transitPrefs struct {
	AllowedTravelModes []string `json:"allowedTravelModes"`
}

type MapsRequest struct {
	Origins            []matrix     `json:"origins"`
	Destinations       []matrix     `json:"destinations"`
	ArrivalTime        time.Time    `json:"arrivalTime"`
	TravelMode         string       `json:"travelMode"`
	Units              string       `json:"units"`
	TransitPreferences transitPrefs `json:"transitPreferences"`
}

func (m *MapsRequest) MarshalJSON() ([]byte, error) {
	type Alias MapsRequest
	return json.Marshal(&struct {
		*Alias
		ArrivalTime string `json:"arrivalTime"`
	}{
		Alias:       (*Alias)(m),
		ArrivalTime: m.ArrivalTime.Format(time.RFC3339),
	})
}

func (m *MapsRequest) UnmarshalJSON(data []byte) error {
	type alias MapsRequest

	aux := &struct {
		*alias
		ArrivalTimeString string `json:"arrivalTime"`
	}{
		alias: (*alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	arrivalTime, err := time.Parse(time.RFC3339, aux.ArrivalTimeString)
	if err != nil {
		return err
	}

	m.ArrivalTime = arrivalTime

	return nil
}

func (m MapsRequest) String() string {
	s, err := m.MarshalJSON()
	if err != nil {
		return ""
	}
	return string(s)
}

func NewMapsRequest(from, to string) MapsRequest {
	return MapsRequest{
		ArrivalTime:        time.Now().Add(24 * time.Hour), // get next working day at 09h
		TravelMode:         "TRANSIT",
		Units:              "METRIC",
		TransitPreferences: transitPrefs{AllowedTravelModes: []string{"RAIL"}},

		Origins: []matrix{
			{
				Waypoint: waypoint{
					Address: from,
				},
			},
		},
		Destinations: []matrix{
			{
				Waypoint: waypoint{
					Address: to,
				},
			},
		},
	}
}

func NewMapsRequestWithCoordinates(from, to Coordinates) MapsRequest {
	return NewMapsRequest(from.String(), to.String())
}
