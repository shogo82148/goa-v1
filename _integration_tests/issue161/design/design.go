package design

import (
	. "github.com/shogo82148/goa-v1/design"
	. "github.com/shogo82148/goa-v1/design/apidsl"
)

var _ = API("media", func() {
	Title("An API exercising the DefaultMedia definition")
	Host("localhost:8080")
	Scheme("http")
})

var Author = Type("Author", func() {
	Attribute("name", String)
	Attribute("country", String)
	Required("name")
})

var Book = Type("Book", func() {
	Attribute("title", String)
	Attribute("author", Author)
})
