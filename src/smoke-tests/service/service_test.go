package service_test

import (
	"fmt"
	"time"

	"github.com/pborman/uuid"

	"github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	"github.com/pivotal-cf-experimental/cf-test-helpers/runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("CF Cassandra Service Broker", func() {
	var timeout = time.Second * 60
	var retryInterval = time.Second * 1
	var appPath = "../assets/app_sinatra_service"

	var appName string

	randomName := func() string {
		return uuid.NewRandom().String()
	}

	appUri := func(appName string) string {
		return "https://" + appName + "." + config.AppsDomain
	}

	BeforeEach(func() {
		appName = randomName()
		Eventually(cf.Cf("push", appName, "-m", "256M", "-p", appPath, "-s", "cflinuxfs2", "-no-start"), config.ScaledTimeout(timeout)).Should(Exit(0))
	})

	AfterEach(func() {
		Eventually(cf.Cf("delete", appName, "-f"), config.ScaledTimeout(timeout)).Should(Exit(0))
	})

	AssertLifeCycleBehavior := func(planName string) {
		It("can create, bind to, write to, read from, unbind, and destroy a service instance using the "+planName+" plan", func() {
			serviceInstanceName := randomName()

			createServiceSession := cf.Cf("create-service", config.ServiceName, planName, serviceInstanceName)
			createServiceSession.Wait(config.ScaledTimeout(timeout))

			createServiceStdout := createServiceSession.Out

			select {
			case <-createServiceStdout.Detect("FAILED"):
				Eventually(createServiceSession, config.ScaledTimeout(timeout)).Should(Say("instance limit for this service has been reached"))
				Eventually(createServiceSession, config.ScaledTimeout(timeout)).Should(Exit(1))
				fmt.Println("No Plan Instances available for testing plan:", planName)
			case <-createServiceStdout.Detect("OK"):
				Eventually(createServiceSession, config.ScaledTimeout(timeout)).Should(Exit(0))
				Eventually(cf.Cf("bind-service", appName, serviceInstanceName), config.ScaledTimeout(timeout)).Should(Exit(0))
				Eventually(cf.Cf("start", appName), config.ScaledTimeout(5*time.Minute)).Should(Exit(0))

				uri := appUri(appName) + "/test"
				fmt.Println("Getting from url: ", uri)
				Eventually(runner.Curl(uri, "-k"), config.ScaledTimeout(timeout), retryInterval).Should(Say("works"))
				fmt.Println("\n")

				Eventually(cf.Cf("unbind-service", appName, serviceInstanceName), config.ScaledTimeout(timeout)).Should(Exit(0))
				Eventually(cf.Cf("delete-service", "-f", serviceInstanceName), config.ScaledTimeout(timeout)).Should(Exit(0))
			}
			createServiceStdout.CancelDetects()

		})
	}

	Context("for each plan", func() {
		for _, planName := range config.PlanNames {
			AssertLifeCycleBehavior(planName)
		}
	})
})
