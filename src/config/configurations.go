package configurations

import (
	"os"
	"strconv"
)

type Configurations struct {
	Port               int
	GenerateCarPath    string
	AggregateCarPath   string
	LighthouseAPIKey   string
	LighthouseDownloadURL string
}

func LoadConfigurations() Configurations {
	port, err := strconv.Atoi(getEnv("PORT", "8000"))
	if err != nil {
		port = 8000 // Default port if parsing fails
	}

	return Configurations{
		Port:               port,
		GenerateCarPath:    getEnv("GENERATE_CAR_PATH", "generated_car/"),
		AggregateCarPath:   getEnv("AGGREGATE_CAR_PATH", "aggregate_car_file/"),
		LighthouseAPIKey:   getEnv("LIGHTHOUSE_API_KEY", ""),
		LighthouseDownloadURL: "https://gateway.lighthouse.storage/ipfs/",
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}