//go:build e2e
// +build e2e

package e2e

import (
	"testing"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/openshift-psap/special-resource-operator/test/framework"
)

const (
	NamespaceSRO = "openshift-special-resource-operator"
	pollInterval = 1 * time.Second
	waitDuration = 10 * time.Minute
)

func TestSRO(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Special Resource Operator e2e tests: basic")
}

var _ = ginkgo.BeforeSuite(func() {
	cs := framework.NewClientSet()
	cl := framework.NewControllerRuntimeClient()

	ginkgo.By("[pre] Creating kube client set...")
	clientSet, err := GetKubeClientSet()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	ginkgo.By("[pre] Checking SRO status...")
	err = WaitSRORunning(clientSet, "openshift-special-resource-operator")
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	ginkgo.By("[pre] Creating preamble...")
	err = CreatePreamble(cl)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	ginkgo.By("[pre] Checking ClusterOperator conditions...")
	err = WaitClusterOperatorConditions(cs)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	ginkgo.By("[pre] Checking ClusterOperator related objects...")
	err = WaitClusterOperatorNamespace(cs)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
})
