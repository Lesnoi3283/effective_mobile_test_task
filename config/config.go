package config

import "os"

type Config struct {
	ServerAddress       string
	DBConnectionString  string
	ExtraDataAPIAddress string
	LogLevel            string
}

// Configure reads environmental variables and returns them in the "Config" struct.
func Configure() Config {
	conf := Config{}

	conf.ServerAddress = os.Getenv("SERVER_ADDRESS")
	conf.DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
	conf.ExtraDataAPIAddress = os.Getenv("EXTRA_DATA_API_ADDRESS")
	conf.LogLevel = os.Getenv("LOG_LEVEL")

	return conf
}
