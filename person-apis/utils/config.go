package utils

import (
	"github.com/spf13/viper"
)

const (
	// GrpcRouterPort is...
	GrpcRouterPort = "grpc_router_port"
	// BirthdayServiceAddress is...
	PersonServiceAddress = "person_service_address"
	// ConfigMongoConnectionString IS ..
	ConfigMongoConnectionString = "mongo_host"
)

// Config is used to declare every variable in
// app.env file in order to use it in go files
type Config struct {
	GrpcRouterPort              string `mapstructure:"GRPC_ROUTER_PORT"`
	PersonServiceAddress        string `mapstructure:"PERSON_SERVICE_ADDRESS"`
	ConfigMongoConnectionString string `mapstructure:"MONGO_HOST"`
}

// LoadConfig loads all variables from app.env file
func LoadConfig() (err error) {
	viper.SetDefault(GrpcRouterPort, "9000")
	viper.SetDefault(PersonServiceAddress, "localhost:50051")
	viper.SetDefault(ConfigMongoConnectionString, "mongodb://root:example@0.0.0.0:27017")
	viper.AutomaticEnv()

	return nil
}