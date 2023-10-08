package urls

import (
	"context"
	"shorty-urls-server/internal/database"
	"shorty-urls-server/internal/routes/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func getMeta(userId string, ctx context.Context) (*utils.UrlMetaData, error) {
	tx := database.DB.Table("shorten_urls su").
		Select("JSON_AGG(JSONB_BUILD_OBJECT('device', suv.device, 'location', suv.\"location\", 'visited_at', suv.created_at)) as visits").
		Joins("LEFT JOIN shorten_url_visits suv ON suv.shorten_url_id = su.id")

	if userId != "" {
		tx = tx.Where("su.user_id = ?", userId)
	}
	tx = tx.Group("su.id")

	rows, err := tx.Rows()
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	defer rows.Close()

	urlDetails := []utils.UrlMetaData{}
	for rows.Next() {
		urlDetail := utils.UrlDetails{}
		if err := tx.ScanRows(rows, &urlDetail); err != nil {
			return nil, fiber.ErrInternalServerError
		}

		response, err := utils.GetUrlDetails(urlDetail, true, ctx)
		if err != nil {
			return nil, fiber.ErrInternalServerError
		}

		urlDetails = append(urlDetails, *response.Meta)
	}

	yearMonthDayClicks := make(map[int]utils.YearData)
	browserClicks := map[string]int{}
	osClicks := map[string]int{}
	deviceClicks := map[string]int{}
	countryClicks := make(map[string]utils.CountryData)

	for _, urlDetail := range urlDetails {
		for year, yearData := range urlDetail.YearMonthDayClicks {
			if yearElement, ok := yearMonthDayClicks[year]; ok {
				yearElement.Count += yearElement.Count
				yearMonthDayClicks[year] = yearElement
			} else {
				yearMonthDayClicks[year] = utils.YearData{
					Count:  yearData.Count,
					Months: map[int]utils.MonthData{},
				}
			}

			for month, monthData := range yearData.Months {
				if monthElement, ok := yearMonthDayClicks[year].Months[month]; ok {
					monthElement.Count += monthElement.Count
					yearMonthDayClicks[year].Months[month] = monthElement
				} else {
					yearMonthDayClicks[year].Months[month] = utils.MonthData{
						Count: monthData.Count,
						Dates: map[int]utils.DateData{},
					}
				}

				for day, dayData := range monthData.Dates {
					if dayElement, ok := yearMonthDayClicks[year].Months[month].Dates[day]; ok {
						dayElement.Count += dayElement.Count
						yearMonthDayClicks[year].Months[month].Dates[day] = dayElement
					} else {
						yearMonthDayClicks[year].Months[month].Dates[day] = utils.DateData{
							Count: dayData.Count,
						}
					}
				}
			}
		}

		for browser, count := range urlDetail.BrowserClicks {
			if _, ok := browserClicks[browser]; ok {
				browserClicks[browser] += count
			} else {
				browserClicks[browser] = count
			}
		}

		for os, count := range urlDetail.OsClicks {
			if _, ok := osClicks[os]; ok {
				osClicks[os] += count
			} else {
				osClicks[os] = count
			}
		}

		for device, count := range urlDetail.DeviceClicks {
			if _, ok := deviceClicks[device]; ok {
				deviceClicks[device] += count
			} else {
				deviceClicks[device] = count
			}
		}

		for country, countryData := range urlDetail.CountryClicks {
			if countryElement, ok := countryClicks[country]; ok {
				countryElement.Count += countryData.Count
				countryClicks[country] = countryElement
			} else {
				countryClicks[country] = utils.CountryData{
					Count: countryData.Count,
					City:  map[string]int{},
				}
			}

			for city, cityData := range countryData.City {
				if _, ok := countryClicks[country].City[city]; ok {
					countryClicks[country].City[city] += cityData
				} else {
					countryClicks[country].City[city] = cityData
				}
			}
		}
	}

	return &utils.UrlMetaData{
		YearMonthDayClicks: yearMonthDayClicks,
		BrowserClicks:      browserClicks,
		OsClicks:           osClicks,
		DeviceClicks:       deviceClicks,
		CountryClicks:      countryClicks,
	}, nil

}
