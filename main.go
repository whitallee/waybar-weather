package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// --- wttr.in response types ---

type WttrResponse struct {
	CurrentCondition []CurrentCondition `json:"current_condition"`
	Weather          []Weather          `json:"weather"`
}

type CurrentCondition struct {
	TempC         string        `json:"temp_C"`
	TempF         string        `json:"temp_F"`
	FeelsLikeC    string        `json:"FeelsLikeC"`
	FeelsLikeF    string        `json:"FeelsLikeF"`
	Humidity      string        `json:"humidity"`
	WeatherCode   string        `json:"weatherCode"`
	WeatherDesc   []WeatherDesc `json:"weatherDesc"`
	WindspeedKmph string        `json:"windspeedKmph"`
}

type WeatherDesc struct {
	Value string `json:"value"`
}

type Weather struct {
	MaxTempC string `json:"maxtempC"`
	MinTempC string `json:"mintempC"`
	MaxTempF string `json:"maxtempF"`
	MinTempF string `json:"mintempF"`
}

// --- Waybar output type ---

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
}

// weatherIcon maps wttr.in weather codes to Nerd Font icons (weather-icons set).
func weatherIcon(code string) string {
	c, _ := strconv.Atoi(code)
	switch {
	case c == 113: // Sunny / Clear
		return "\ue30d"
	case c == 116: // Partly cloudy
		return "\ue302"
	case c == 119, c == 122: // Cloudy / Overcast
		return "\ue33d"
	case c >= 143 && c <= 260: // Fog / Mist / Blizzard
		return "\ue313"
	case c >= 200 && c <= 232: // Thunderstorm
		return "\ue31d"
	case c >= 263 && c <= 296: // Drizzle / Light rain
		return "\ue319"
	case c >= 299 && c <= 314: // Rain
		return "\ue318"
	case c >= 317 && c <= 338: // Snow
		return "\ue31a"
	case c >= 350 && c <= 377: // Sleet / Ice
		return "\ue3ac"
	default:
		return "\ue331" // Thermometer fallback
	}
}

func main() {
	fahrenheit := flag.Bool("f", false, "display temperature in Fahrenheit")
	flag.Parse()

	location := "auto" // auto-detect by IP
	if flag.NArg() > 0 {
		location = flag.Arg(0)
	}

	url := fmt.Sprintf("https://wttr.in/%s?format=j1", location)

	const maxAttempts = 5
	var wttr WttrResponse
	var fetchErr string
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err := http.Get(url)
		if err != nil {
			fetchErr = "network error"
		} else if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fetchErr = "bad response"
		} else if err := json.NewDecoder(resp.Body).Decode(&wttr); err != nil {
			resp.Body.Close()
			fetchErr = "parse error"
		} else {
			resp.Body.Close()
			if len(wttr.CurrentCondition) == 0 {
				fetchErr = "no data"
			} else {
				fetchErr = ""
				break
			}
		}
		if attempt < maxAttempts {
			time.Sleep(1 * time.Second)
		}
	}
	if fetchErr != "" {
		printError(fetchErr)
		return
	}

	cur := wttr.CurrentCondition[0]
	icon := weatherIcon(cur.WeatherCode)
	desc := ""
	if len(cur.WeatherDesc) > 0 {
		desc = cur.WeatherDesc[0].Value
	}

	var text, tooltip string
	if *fahrenheit {
		text = fmt.Sprintf("%s %s°F", icon, cur.TempF)
		tooltip = fmt.Sprintf("%s\nFeels like: %s°F\nHumidity: %s%%\nWind: %s km/h",
			desc, cur.FeelsLikeF, cur.Humidity, cur.WindspeedKmph,
		)
		if len(wttr.Weather) > 0 {
			w := wttr.Weather[0]
			tooltip += fmt.Sprintf("\nHigh: %s°F  Low: %s°F", w.MaxTempF, w.MinTempF)
		}
	} else {
		text = fmt.Sprintf("%s %s°C", icon, cur.TempC)
		tooltip = fmt.Sprintf("%s\nFeels like: %s°C\nHumidity: %s%%\nWind: %s km/h",
			desc, cur.FeelsLikeC, cur.Humidity, cur.WindspeedKmph,
		)
		if len(wttr.Weather) > 0 {
			w := wttr.Weather[0]
			tooltip += fmt.Sprintf("\nHigh: %s°C  Low: %s°C", w.MaxTempC, w.MinTempC)
		}
	}

	out := WaybarOutput{
		Text:    text,
		Tooltip: tooltip,
		Class:   "weather",
	}

	b, _ := json.Marshal(out)
	fmt.Println(string(b))
}

func printError(msg string) {
	out := WaybarOutput{Text: "⚠️ " + msg, Class: "weather-error"}
	b, _ := json.Marshal(out)
	fmt.Println(string(b))
}
