package helpers

import (
	"encoding/json"
	"os"

	. "github.com/cloudfoundry-incubator/cf-test-helpers/services/context_setup"
)

type CassandraIntegrationConfig struct {
	IntegrationConfig

	ServiceName string `json:"service_name"`
	PlanName    string `json:"plan_name"`
	BrokerHost  string `json:"broker_host"`
}

func LoadConfig() (config CassandraIntegrationConfig) {
	path := os.Getenv("CONFIG")
	if path == "" {
		panic("Must set $CONFIG to point to an integration config .json file.")
	}

	return LoadPath(path)
}

func LoadPath(path string) (config CassandraIntegrationConfig) {
	configFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	if config.ApiEndpoint == "" {
		panic("missing configuration 'api'")
	}

	if config.AdminUser == "" {
		panic("missing configuration 'admin_user'")
	}

	if config.ApiEndpoint == "" {
		panic("missing configuration 'admin_password'")
	}

	if config.ServiceName == "" {
		panic("missing configuration 'service_name'")
	}

	if config.PlanName == "" {
		panic("missing configuration 'plan_name'")
	}

	if config.BrokerHost == "" {
		panic("missing configuration 'broker_host'")
	}

	if config.TimeoutScale <= 0 {
		config.TimeoutScale = 1
	}

	return
}
