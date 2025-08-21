package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
)

// Base de datos de códigos de área mexicanos con coordenadas aproximadas
var mexicanAreaCodes = map[string]AreaCodeInfo{
	"664": {Region: "Tijuana", State: "Baja California", Lat: 32.5027, Lng: -117.0083, Radius: "~15km"},
	"663": {Region: "Ensenada", State: "Baja California", Lat: 31.8590, Lng: -116.5969, Radius: "~20km"},
	"665": {Region: "Mexicali", State: "Baja California", Lat: 32.6519, Lng: -115.4683, Radius: "~25km"},
	"646": {Region: "Tecate", State: "Baja California", Lat: 32.5764, Lng: -116.6294, Radius: "~10km"},
	"55":  {Region: "Ciudad de México", State: "CDMX", Lat: 19.4326, Lng: -99.1332, Radius: "~50km"},
	"33":  {Region: "Guadalajara", State: "Jalisco", Lat: 20.6597, Lng: -103.3496, Radius: "~40km"},
	"81":  {Region: "Monterrey", State: "Nuevo León", Lat: 25.6866, Lng: -100.3161, Radius: "~35km"},
}

type AreaCodeInfo struct {
	Region string  `json:"region"`
	State  string  `json:"state"`
	Lat    float64 `json:"latitude"`
	Lng    float64 `json:"longitude"`
	Radius string  `json:"coverage_radius"`
}

// Estructura para respuesta de ipapi.com
type IPApiResponse struct {
	IP        string  `json:"ip"`
	City      string  `json:"city"`
	Region    string  `json:"region_name"`
	Country   string  `json:"country_name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	ISP       string  `json:"org"`
	Timezone  string  `json:"timezone"`
}

// Estructura para nuestro escáner mejorado
type IpapiScanner struct{}

type IpapiScannerResponse struct {
	Method           string  `json:"method" console:"Method"`
	AreaCode         string  `json:"area_code,omitempty" console:"Area Code,omitempty"`
	ApproximateCity  string  `json:"approximate_city,omitempty" console:"Approximate City,omitempty"`
	State            string  `json:"state,omitempty" console:"State,omitempty"`
	Country          string  `json:"country,omitempty" console:"Country,omitempty"`
	Latitude         float64 `json:"latitude,omitempty" console:"Latitude,omitempty"`
	Longitude        float64 `json:"longitude,omitempty" console:"Longitude,omitempty"`
	CoverageRadius   string  `json:"coverage_radius,omitempty" console:"Coverage Radius,omitempty"`
	GoogleMapsLink   string  `json:"google_maps_link,omitempty" console:"Google Maps,omitempty"`
	SearchRadius     string  `json:"search_radius,omitempty" console:"Search Radius,omitempty"`
	Note             string  `json:"note" console:"Note"`
}

func (s *IpapiScanner) Name() string {
	return "geolocation"
}

func (s *IpapiScanner) Description() string {
	return "Approximate geolocation based on area code and carrier information"
}

func (s *IpapiScanner) DryRun(_ number.Number, opts remote.ScannerOptions) error {
	// Este escáner no requiere API key para códigos de área mexicanos
	return nil
}

func (s *IpapiScanner) Run(n number.Number, opts remote.ScannerOptions) (interface{}, error) {
	result := IpapiScannerResponse{
		Method:  "Area Code Database + Regional Analysis",
		Country: n.Country,
	}

	// Extraer código de área del número mexicano
	if n.Country == "MX" {
		areaCode := extractMexicanAreaCode(n.RawLocal)
		if areaCode != "" {
			if info, exists := mexicanAreaCodes[areaCode]; exists {
				result.AreaCode = areaCode
				result.ApproximateCity = info.Region
				result.State = info.State
				result.Latitude = info.Lat
				result.Longitude = info.Lng
				result.CoverageRadius = info.Radius
				result.GoogleMapsLink = fmt.Sprintf("https://www.google.com/maps/@%f,%f,12z", info.Lat, info.Lng)
				result.SearchRadius = fmt.Sprintf("Búsqueda recomendada en un radio de %s desde %s", info.Radius, info.Region)
				result.Note = fmt.Sprintf("Coordenadas aproximadas basadas en código de área %s para %s, %s", areaCode, info.Region, info.State)
				return result, nil
			}
		}
	}

	// Si no es México o no encontramos el código de área, usar ipapi como fallback
	apiKey := opts.GetStringEnv("IPAPI_API_KEY")
	if apiKey != "" {
		ipData, err := getIPLocationData(apiKey)
		if err == nil {
			result.Method = "IP Geolocation (Fallback)"
			result.ApproximateCity = ipData.City
			result.Country = ipData.Country
			result.Latitude = ipData.Latitude
			result.Longitude = ipData.Longitude
			result.GoogleMapsLink = fmt.Sprintf("https://www.google.com/maps/@%f,%f,10z", ipData.Latitude, ipData.Longitude)
			result.Note = "Ubicación aproximada basada en geolocalización IP regional"
		}
	}

	if result.ApproximateCity == "" {
		result.Note = "No se pudo determinar ubicación aproximada. Considere agregar IPAPI_API_KEY para más datos."
	}

	return result, nil
}

func extractMexicanAreaCode(rawLocal string) string {
	// Para números mexicanos, extraer código de área
	if len(rawLocal) >= 10 {
		// Intentar códigos de 3 dígitos primero
		if len(rawLocal) >= 3 {
			code3 := rawLocal[:3]
			if _, exists := mexicanAreaCodes[code3]; exists {
				return code3
			}
		}
		// Luego códigos de 2 dígitos
		if len(rawLocal) >= 2 {
			code2 := rawLocal[:2]
			if _, exists := mexicanAreaCodes[code2]; exists {
				return code2
			}
		}
	}
	return ""
}

func getIPLocationData(apiKey string) (*IPApiResponse, error) {
	url := fmt.Sprintf("http://api.ipapi.com/check?access_key=%s&format=1", apiKey)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var ipApiData IPApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipApiData); err != nil {
		return nil, err
	}
	
	return &ipApiData, nil
}

// Función requerida para plugins de PhoneInfoGA
func init() {
	remote.RegisterPlugin(&IpapiScanner{})
}
