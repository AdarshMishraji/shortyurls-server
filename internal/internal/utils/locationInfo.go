package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LocationInfo struct {
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Timezone      string  `json:"timezone"`
	Offset        int     `json:"offset"`
	Currency      string  `json:"currency"`
	Ip            string  `json:"ip"`
}

func (info *LocationInfo) String() string {
	return fmt.Sprintf(
		`{"Continent":"%s","ContinentCode":"%s","Country":"%s","CountryCode":"%s","Region":"%s","RegionName":"%s","City":"%s","Zip":"%s","Lat":%f,"Lon":%f,"Timezone":"%s","Offset":%d,"Currency":"%s","Ip":"%s"}`,
		info.Continent,
		info.ContinentCode,
		info.Country,
		info.CountryCode,
		info.Region,
		info.RegionName,
		info.City,
		info.Zip,
		info.Lat,
		info.Lon,
		info.Timezone,
		info.Offset,
		info.Currency,
		info.Ip,
	)
}

func GetLocationInfo(ip string) LocationInfo {
	rootUrl := fmt.Sprintf("https://ip-api.com/json/%s?fields=continent,continentCode,country,countryCode,region,regionName,city,zip,lat,lon,timezone,offset,currency", ip)

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		return LocationInfo{
			Ip: ip,
		}
	}

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return LocationInfo{
			Ip: ip,
		}
	}

	if res.StatusCode != http.StatusOK {
		return LocationInfo{
			Ip: ip,
		}
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationInfo{
			Ip: ip,
		}
	}

	var LocationRes map[string]interface{}

	if err := json.Unmarshal(resBody, &LocationRes); err != nil {
		return LocationInfo{
			Ip: ip,
		}
	}

	return LocationInfo{
		Continent:     LocationRes["continent"].(string),
		ContinentCode: LocationRes["continentCode"].(string),
		Country:       LocationRes["country"].(string),
		CountryCode:   LocationRes["countryCode"].(string),
		Region:        LocationRes["region"].(string),
		RegionName:    LocationRes["regionName"].(string),
		City:          LocationRes["city"].(string),
		Zip:           LocationRes["zip"].(string),
		Lat:           LocationRes["lat"].(float64),
		Lon:           LocationRes["lon"].(float64),
		Timezone:      LocationRes["timezone"].(string),
		Offset:        int(LocationRes["offset"].(float64)),
		Currency:      LocationRes["currency"].(string),
		Ip:            LocationRes["query"].(string),
	}
}

func GetLocationInfoFromContext(ctx context.Context) LocationInfo {
	return ctx.Value("location").(LocationInfo)
}

func SetLocationInfoToContext(ctx *fiber.Ctx, ip *string, cachedLocationInfo *LocationInfo) *LocationInfo {
	var locationInfo LocationInfo

	if (cachedLocationInfo != nil) && (cachedLocationInfo.String() == locationInfo.String()) {
		locationInfo = *cachedLocationInfo
	} else {
		locationInfo = GetLocationInfo(*ip)
	}

	ctx.SetUserContext(context.WithValue(ctx.UserContext(), "location", locationInfo))
	return &locationInfo
}
