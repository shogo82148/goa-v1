package goakit_test

import (
	"bytes"

	"github.com/go-kit/kit/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/goa-v1"
	goakit "github.com/shogo82148/goa-v1/logging/kit"
)

var _ = Describe("New", func() {
	var buf bytes.Buffer
	var logger log.Logger
	var adapter goa.LogAdapter

	BeforeEach(func() {
		buf.Reset()
		logger = log.NewLogfmtLogger(&buf)
		adapter = goakit.New(logger)
	})

	It("creates an adapter that logs", func() {
		msg := "msg"
		adapter.Info(msg)
		Ω(buf.String()).Should(Equal("lvl=info msg=" + msg + "\n"))
	})

	It("creates an adapter that logs", func() {
		adapter := adapter.(goa.WarningLogAdapter)
		msg := "msg"
		adapter.Warn(msg)
		Ω(buf.String()).Should(Equal("lvl=warn msg=" + msg + "\n"))
	})

	It("creates an adapter that logs", func() {
		msg := "msg"
		adapter.Error(msg)
		Ω(buf.String()).Should(Equal("lvl=error msg=" + msg + "\n"))
	})
})
