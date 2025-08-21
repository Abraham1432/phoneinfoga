# Geolocation Plugin for PhoneInfoGA

## Overview

This plugin adds approximate geolocation capabilities to PhoneInfoGA by using area code databases and IP geolocation services to provide regional coordinates for phone numbers.

## Features

- **Area Code Database**: Built-in database of Mexican area codes with precise coordinates
- **Regional Coordinates**: Provides latitude/longitude for phone number regions
- **Coverage Radius**: Indicates search radius for each location
- **Google Maps Integration**: Generates direct links to Google Maps
- **Fallback Support**: Uses IP geolocation APIs for non-Mexican numbers

## Installation

1. Compile the plugin:
```bash
cd plugins
go build -buildmode=plugin -o ipapi_scanner.so ipapi_scanner.go
```

2. Use with PhoneInfoGA:
```bash
./bin/phoneinfoga scan -n "+526647639100" --plugin ./plugins/ipapi_scanner.so
```

## Configuration

### Optional Environment Variables

- `IPAPI_API_KEY`: API key for ipapi.com (for fallback geolocation)

### Supported Area Codes

Currently supports Mexican area codes:

| Area Code | Region | State | Coordinates |
|-----------|---------|-------|-------------|
| 664 | Tijuana | Baja California | 32.5027, -117.0083 |
| 663 | Ensenada | Baja California | 31.8590, -116.5969 |
| 665 | Mexicali | Baja California | 32.6519, -115.4683 |
| 646 | Tecate | Baja California | 32.5764, -116.6294 |
| 55 | Ciudad de México | CDMX | 19.4326, -99.1332 |
| 33 | Guadalajara | Jalisco | 20.6597, -103.3496 |
| 81 | Monterrey | Nuevo León | 25.6866, -100.3161 |

## Example Output

```
Results for geolocation
Method: Area Code Database + Regional Analysis
Area Code: 664
Approximate City: Tijuana
State: Baja California
Country: MX
Latitude: 32.502700
Longitude: -117.008300
Coverage Radius: ~15km
Google Maps: https://www.google.com/maps/@32.502700,-117.008300,12z
Search Radius: Búsqueda recomendada en un radio de ~15km desde Tijuana
Note: Coordenadas aproximadas basadas en código de área 664 para Tijuana, Baja California
```

## Limitations

- **Approximate Location**: Provides regional coordinates, not exact device location
- **Privacy Compliant**: Does not access real-time device data
- **Area Code Dependent**: Accuracy depends on area code coverage
- **Regional Coverage**: Currently focused on Mexican area codes

## Contributing

To add more area codes or regions:

1. Update the `mexicanAreaCodes` map in `ipapi_scanner.go`
2. Add coordinates for the new region
3. Test with numbers from that area code
4. Update documentation

## API Integration

The plugin can integrate with various IP geolocation services:

- **ipapi.com**: Primary fallback service
- **Extensible**: Can be modified to support other services

## Technical Details

- **Language**: Go
- **Plugin System**: Uses Go's plugin architecture
- **Interface**: Implements PhoneInfoGA Scanner interface
- **Dependencies**: Standard library + PhoneInfoGA libs

## License

This plugin follows the same license as PhoneInfoGA (GNU GPL v3.0).
