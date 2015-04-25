package lifecycle_test

import (
	. "github.com/cloudfoundry-incubator/cf-test-helpers/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"../../helpers"
	context_setup "github.com/cloudfoundry-incubator/cf-test-helpers/services/context_setup"
	"testing"
)

func TestService(t *testing.T) {
	context_setup.TimeoutScale = CassandraIntegrationConfig.TimeoutScale

	context_setup.SetupEnvironment(context_setup.NewContext(CassandraIntegrationConfig.IntegrationConfig, "CassandraATS"))

	RegisterFailHandler(Fail)
	RunSpecs(t, "Lifecycle tests")
}

func AppUri(appname string) string {
	return "http://" + appname + "." + CassandraIntegrationConfig.AppsDomain
}

func Curling(args ...string) func() *gexec.Session {
	return func() *gexec.Session {
		return Curl(args...)
	}
}

func ServiceName() string {
	return CassandraIntegrationConfig.ServiceName
}

func PlanName() string {
	return CassandraIntegrationConfig.PlanName
}

var CassandraIntegrationConfig = helpers.LoadConfig()

var AppName = ""

var testAppPath = "../../assets/app_sinatra_service"
