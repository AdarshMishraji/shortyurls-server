package utils

import (
	"context"
	"encoding/json"
	"shorty-urls-server/internal/database"
	"shorty-urls-server/internal/internal/utils"
	"time"
)

type UrlDetails struct {
	Visits json.RawMessage `json:"visits" gorm:"column:visits"`
	database.ShortenURL
}

type UrlMetaData struct {
	YearMonthDayClicks map[int]YearData       `json:"yearMonthDayClicks"`
	BrowserClicks      map[string]int         `json:"browserClicks"`
	OsClicks           map[string]int         `json:"osClicks"`
	DeviceClicks       map[string]int         `json:"deviceClicks"`
	CountryClicks      map[string]CountryData `json:"countryClicks"`
}

type UrlDetailsResponse struct {
	Meta *UrlMetaData `json:"meta,omitempty"`
	database.ShortenURL
}

type UrlVisitData struct {
	Device    utils.DeviceInfo   `json:"device"`
	Location  utils.LocationInfo `json:"location"`
	VisitedAt string             `json:"visited_at"`
}

type DateData struct {
	Count int
}

type MonthData struct {
	Count int
	Dates map[int]DateData
}

type YearData struct {
	Count  int
	Months map[int]MonthData
}

type CountryData struct {
	Count int
	City  map[string]int
}

func GetUrlDetails(urlDetails UrlDetails, withMeta bool, ctx context.Context) (*UrlDetailsResponse, error) {
	response := &UrlDetailsResponse{
		ShortenURL: database.ShortenURL{
			OriginalURL: urlDetails.OriginalURL,
			Alias:       urlDetails.Alias,
			IsActive:    urlDetails.IsActive,
			IsDeleted:   urlDetails.IsDeleted,
		},
	}

	if urlDetails.Visits == nil || !withMeta {
		return response, nil
	}

	visits := []UrlVisitData{}

	if err := json.Unmarshal(urlDetails.Visits, &visits); err != nil {
		return nil, err
	}

	timeZone := utils.GetLocationInfoFromContext(ctx).Timezone
	time.Local, _ = time.LoadLocation(timeZone)

	yearMonthDayClicks := make(map[int]YearData)
	browserClicks := map[string]int{}
	osClicks := map[string]int{}
	deviceClicks := map[string]int{}
	countryClicks := make(map[string]CountryData)
	for i := 0; i < len(visits); i++ {

		if visits[i].VisitedAt == "" {
			continue
		}

		date, err := time.Parse(time.RFC3339, visits[i].VisitedAt)
		if err != nil {
			return nil, err
		}

		date = date.In(time.Local)
		year, _month, day := date.Date()
		month := int(_month)

		deviceInfo := visits[i].Device
		locationInfo := visits[i].Location

		osName := deviceInfo.OsName
		deviceType := deviceInfo.DeviceType
		browserName := deviceInfo.BrowserName
		country := locationInfo.Country
		city := locationInfo.City

		if _yearData, ok := yearMonthDayClicks[year]; ok {
			_yearData.Count++
			yearMonthDayClicks[year] = _yearData
		} else {
			yearMonthDayClicks[year] = YearData{
				Count:  1,
				Months: make(map[int]MonthData),
			}
		}

		if _monthData, ok := yearMonthDayClicks[year].Months[month]; ok {
			_monthData.Count++
			yearMonthDayClicks[year].Months[month] = _monthData
		} else {
			yearMonthDayClicks[year].Months[month] = MonthData{
				Count: 1,
				Dates: make(map[int]DateData),
			}
		}

		if _dateData, ok := yearMonthDayClicks[year].Months[month].Dates[day]; ok {
			_dateData.Count++
			yearMonthDayClicks[year].Months[month].Dates[day] = _dateData
		} else {
			yearMonthDayClicks[year].Months[month].Dates[day] = DateData{
				Count: 1,
			}
		}

		if browserName != "" {
			if _, ok := browserClicks[browserName]; ok {
				browserClicks[browserName]++
			} else {
				browserClicks[browserName] = 1
			}
		}

		if osName != "" {
			if _, ok := osClicks[osName]; ok {
				osClicks[osName]++
			} else {
				osClicks[osName] = 1
			}
		}

		if deviceType != "" {
			if _, ok := deviceClicks[deviceType]; ok {
				deviceClicks[deviceType]++
			} else {
				deviceClicks[deviceType] = 1
			}
		}

		if country != "" && city != "" {
			if _countryData, ok := countryClicks[country]; ok {
				_countryData.Count++
				if _, ok := _countryData.City[city]; ok {
					_countryData.City[city]++
				} else {
					_countryData.City[city] = 1
				}
				countryClicks[country] = _countryData
			} else {
				countryClicks[country] = CountryData{
					Count: 1,
					City:  make(map[string]int),
				}
			}
		}
	}

	response.Meta = &UrlMetaData{
		YearMonthDayClicks: yearMonthDayClicks,
		BrowserClicks:      browserClicks,
		OsClicks:           osClicks,
		DeviceClicks:       deviceClicks,
		CountryClicks:      countryClicks,
	}

	return response, nil
}
