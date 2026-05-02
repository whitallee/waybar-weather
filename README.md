# waybar-weather

A small Go binary that fetches current weather from [wttr.in](https://wttr.in) and outputs it in Waybar's custom module JSON format — with a Nerd Font icon, temperature, and a tooltip showing feels-like, humidity, wind, and daily high/low.

## Usage

```bash
# Build
go build -o waybar-weather

# Run (auto-detects location by IP)
./waybar-weather

# Fahrenheit
./waybar-weather -f

# Specific location
./waybar-weather "Austin, TX"
./waybar-weather -f "Austin, TX"
```

## Waybar Config

```json
"custom/weather": {
  "exec": "/path/to/waybar-weather -f",
  "interval": 900,
  "return-type": "json",
  "tooltip": true
}
```

Requires a [Nerd Font](https://www.nerdfonts.com/) for the weather icons.
