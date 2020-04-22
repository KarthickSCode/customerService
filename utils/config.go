package utils

type Config struct {
	Port            int                `json:"port"`
	MongoAddressURI string             `json:"mongoURI"`
	DbName          string             `json:"dbName"`
	DbCollection    DbCollectionConfig `json:"dbCollection"`
	IdGenerationKey IdGenerationConfig `json:"idGenerationKey"`
}

func DefaultConfig() *Config {
	return &Config{
		Port:            8080,
		MongoAddressURI: "mongodb://localhost",
		DbName:          "erply",
		DbCollection:    DefaultDbCollectionConfig(),
		IdGenerationKey: DefaultIdGenerationConfig(),
	}
}

type DbCollectionConfig struct {
	Customer string `json:"customer"`
	IdGen    string `json:"idGen"`
}

func DefaultDbCollectionConfig() DbCollectionConfig {
	return DbCollectionConfig{
		Customer: "customer",
		IdGen:    "idGenerator",
	}
}

type IdGenerationConfig struct {
	CustomerKey string `json:"customerKey"`
}

func DefaultIdGenerationConfig() IdGenerationConfig {
	return IdGenerationConfig{
		CustomerKey: "customer",
	}
}
