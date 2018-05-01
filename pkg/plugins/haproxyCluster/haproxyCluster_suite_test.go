package haproxyCluster_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHaproxyCluster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HaproxyCluster Suite")
}
