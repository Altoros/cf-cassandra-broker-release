package lifecycle_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"fmt"
	. "github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	. "github.com/cloudfoundry-incubator/cf-test-helpers/runner"
	"time"
)

var _ = Describe("Cassandra Lifecycle", func() {
	BeforeEach(func() {
		AppName = RandomName()

		Eventually(Cf("push", AppName, "-m", "128M", "-p", testAppPath, "-no-start"), 60*time.Second).Should(Exit(0))
	})

	AfterEach(func() {
		Eventually(Cf("delete", AppName, "-f"), 60*time.Second).Should(Exit(0))
	})

	It("Allows users to create, bind, write to, read from, unbind, and destroy the service instance", func() {
		ServiceName := ServiceName()
		PlanName := PlanName()
		ServiceInstanceName := RandomName()

		Eventually(Cf("create-service", ServiceName, PlanName, ServiceInstanceName), 60*time.Second).Should(Exit(0))
		// Bind & start
		Eventually(Cf("bind-service", AppName, ServiceInstanceName), 60*time.Second).Should(Exit(0))
		Eventually(Cf("start", AppName), 5*60*time.Second).Should(Exit(0))

		testUri := AppUri(AppName) + "/test"

		fmt.Println("Curling url: ", testUri)
		Eventually(Curl(testUri), 10.0, 1.0).Should(Say("works"))
		fmt.Println("\n")

		Eventually(Cf("unbind-service", AppName, ServiceInstanceName), 60*time.Second).Should(Exit(0))

		Eventually(Cf("delete-service", "-f", ServiceInstanceName), 60*time.Second).Should(Exit(0))
	})
})
