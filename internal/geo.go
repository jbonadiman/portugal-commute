package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/jbonadiman/portugal-commute/internal/models"
)

const portugalGeoHost = "https://json.geoapi.pt"
const googleMapsHost = "https://routes.googleapis.com/distanceMatrix/v2:computeRouteMatrix?key=%s&fields=duration,distanceMeters,originIndex"

// got from https://pt.wikipedia.org/wiki/Lista_de_esta%C3%A7%C3%B5es_do_Metropolitano_de_Lisboa
var metroStations = map[string]models.Coordinates{
	"Aeroporto":             {Latitude: 38.76861, Longitude: -9.12861},
	"Alameda":               {Latitude: 38.73713, Longitude: -9.13388},
	"Alfornelos":            {Latitude: 38.76038, Longitude: -9.20435},
	"Alto dos Moinhos":      {Latitude: 38.74994, Longitude: -9.18003},
	"Alvalade":              {Latitude: 38.75311, Longitude: -9.14396},
	"Amadora Este":          {Latitude: 38.75847, Longitude: -9.21803},
	"Ameixoeira":            {Latitude: 38.77937, Longitude: -9.15962},
	"Anjos":                 {Latitude: 38.72715, Longitude: -9.13485},
	"Areeiro":               {Latitude: 38.74233, Longitude: -9.13354},
	"Arroios":               {Latitude: 38.73266, Longitude: -9.13430},
	"Avenida":               {Latitude: 38.71948, Longitude: -9.14524},
	"Baixa-Chiado":          {Latitude: 38.71051, Longitude: -9.14006},
	"Bela Vista":            {Latitude: 38.74770, Longitude: -9.11763},
	"Cabo Ruivo":            {Latitude: 38.76293, Longitude: -9.10499},
	"Cais do Sodré":         {Latitude: 38.70609, Longitude: -9.14684},
	"Campo Grande":          {Latitude: 38.76012, Longitude: -9.15780},
	"Campo Pequeno":         {Latitude: 38.74076, Longitude: -9.14664},
	"Carnide":               {Latitude: 38.75919, Longitude: -9.19267},
	"Chelas":                {Latitude: 38.75479, Longitude: -9.11392},
	"Cidade Universitária":  {Latitude: 38.75156, Longitude: -9.15911},
	"Colégio Militar / Luz": {Latitude: 38.75336, Longitude: -9.18918},
	"Encarnação":            {Latitude: 38.77502, Longitude: -9.11555},
	"Entre Campos":          {Latitude: 38.74708, Longitude: -9.14823},
	"Intendente":            {Latitude: 38.72320, Longitude: -9.13523},
	"Jardim Zoológico":      {Latitude: 38.74140, Longitude: -9.16848},
	"Laranjeiras":           {Latitude: 38.74847, Longitude: -9.17253},
	"Lumiar":                {Latitude: 38.77339, Longitude: -9.15941},
	"Marquês de Pombal":     {Latitude: 38.72431, Longitude: -9.14931},
	"Martim Moniz":          {Latitude: 38.71753, Longitude: -9.13583},
	"Moscavide":             {Latitude: 38.77488, Longitude: -9.10292},
	"Odivelas":              {Latitude: 38.79334, Longitude: -9.17301},
	"Olaias":                {Latitude: 38.74003, Longitude: -9.12333},
	"Olivais":               {Latitude: 38.76101, Longitude: -9.11198},
	"Oriente":               {Latitude: 38.76787, Longitude: -9.09959},
	"Parque":                {Latitude: 38.72892, Longitude: -9.15051},
	"Picoas":                {Latitude: 38.73037, Longitude: -9.14692},
	"Pontinha":              {Latitude: 38.76222, Longitude: -9.19684},
	"Praça de Espanha":      {Latitude: 38.73775, Longitude: -9.15927},
	"Quinta das Conchas":    {Latitude: 38.76737, Longitude: -9.15557},
	"Rato":                  {Latitude: 38.72015, Longitude: -9.15479},
	"Reboleira":             {Latitude: 38.75227, Longitude: -9.22408},
	"Restauradores":         {Latitude: 38.71590, Longitude: -9.14212},
	"Roma":                  {Latitude: 38.74784, Longitude: -9.14115},
	"Rossio":                {Latitude: 38.71402, Longitude: -9.13794},
	"Saldanha":              {Latitude: 38.73471, Longitude: -9.14521},
	"Santa Apolónia":        {Latitude: 38.71369, Longitude: -9.12247},
	"São Sebastião":         {Latitude: 38.73392, Longitude: -9.15358},
	"Senhor Roubado":        {Latitude: 38.78564, Longitude: -9.17182},
	"Telheiras":             {Latitude: 38.75997, Longitude: -9.16643},
	"Terreiro do Paço":      {Latitude: 38.70703, Longitude: -9.13441},
}

// func GetRealDistance(from, to models.Coordinates) float64 {

// }

func sortStationsByDistance(distancesMap map[string]float64) *[]string {
	// sort stations by distance
	stationsSlice := make([]string, 0, len(distancesMap))

	for station := range distancesMap {
		stationsSlice = append(stationsSlice, station)
	}

	sort.SliceStable(stationsSlice, func(i, j int) bool {
		return distancesMap[stationsSlice[i]] < distancesMap[stationsSlice[j]]
	})

	return &stationsSlice
}

func GetClosestMetroStations(origin models.Coordinates, maxDuration time.Duration) ([]string, map[string]float64) {
	distances := make(map[string]float64)
	for station, coordinates := range metroStations {
		distance := origin.KmDistanceTo(coordinates)
		distances[station] = distance
	}

	return *sortStationsByDistance(distances), distances

	// var farthestStation string
	// for _, station := range stationsByDistance {

	// }
}

func GetGeoData() ([]models.Municipality, error) {
	response, err := http.Get(fmt.Sprintf("%s/municipios/freguesias", portugalGeoHost))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var municipalities []models.Municipality
	if err := json.NewDecoder(response.Body).Decode(&municipalities); err != nil {
		return nil, err
	}

	return municipalities, nil
}
