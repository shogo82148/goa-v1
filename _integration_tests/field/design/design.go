package design

import (
	. "github.com/shogo82148/goa-v1/design"
	. "github.com/shogo82148/goa-v1/design/apidsl"
)

var _ = API("field", func() {
	Title("An API exercising the custom parameter name definition")
	Host("localhost:8080")
	Scheme("http")
})

var MultimediaListMedia = MediaType("application/vnd.multimedialist+json", func() {
	Description("multimedia list")
	Attributes(func() {
		Attribute("media", ArrayOf(Multimedia), func() {
			Metadata("struct:field:name", "MediaList")
			Description("A required array field in the parent media type")
		})
		Required("media")
	})

	View("default", func() {
		Attribute("media")
	})
})

var Multimedia = MediaType("application/vnd.multimedia+json", func() {
	Attributes(func() {
		Attribute("id", Integer, func() {
			Metadata("struct:field:name", "MediaID")
			Description("Media ID")
		})
		Attribute("url", String, func() {
			Metadata("struct:field:name", "MediaURL")
			Description("Media URL")
		})
		Attribute("optional_note", String, func() {
			Metadata("struct:field:name", "Note")
			Description("An optional string field in the Multimedia")
		})
		Required("id", "url")
	})

	View("default", func() {
		Attribute("id")
		Attribute("url")
		Attribute("optional_note")
	})
})

var UploadPayload = Type("UploadPayload", func() {
	Attribute("id", Integer, func() {
		Description("A required int field in the parent type.")
	})
	Attribute("file1", File, func() {
		Metadata("struct:field:name", "FilePrimary")
		Description("A required file field in the parent type.")
	})
	Attribute("file2", File, func() {
		Metadata("struct:field:name", "FileSecondary")
		Description("An optional file field in the parent type.")
	})
	Required("id", "file1")
})

var _ = Resource("Multimedia", func() {
	Action("list", func() {
		Routing(GET("/"))
		Response(OK, MultimediaListMedia)
	})

	Action("get", func() {
		Routing(GET("/:id"))
		Params(func() {
			Param("id", Integer)
			Required("id")
		})
		Response(OK, Multimedia)
		Response(NotFound)
	})

	Action("upload", func() {
		Routing(POST("/upload"))
		MultipartForm()
		Payload(UploadPayload)
		Response(OK)
	})
})
