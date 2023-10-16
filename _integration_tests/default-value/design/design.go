package design

import (
	"time"

	. "github.com/shogo82148/goa-v1/design"
	. "github.com/shogo82148/goa-v1/design/apidsl"
)

var _ = API("default-time", func() {
	Title("This API has time.Time default values")
	Host("localhost:8080")
	Scheme("https")
})

var _ = Resource("timetest", func() {
	Action("check", func() {
		Routing(GET("/"))
		Params(func() {
			Param("times", DateTime, func() {
				Default(time.Time{}.Format(time.RFC3339))
			})
		})
		Response(OK)
	})
})
