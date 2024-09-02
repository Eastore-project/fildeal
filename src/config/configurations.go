package configurations

type Configurations struct {
	Port               int
	GenerateCarPath    string
	AggregateCarPath   string
}

func LoadConfigurations() Configurations {

	return Configurations{
		Port:                8000,
		GenerateCarPath:     "generated_car/", //store car of file from download folder
		AggregateCarPath:    "aggregate_car_file/",
	}
}

