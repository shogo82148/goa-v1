package goaslog_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSlog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Slog Suite")
}
