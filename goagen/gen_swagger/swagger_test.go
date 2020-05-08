package genswagger_test

import (
	"bytes"
	"encoding/json"

	"github.com/go-openapi/loads"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/goa-v1/design"
	"github.com/shogo82148/goa-v1/design/apidsl"
	"github.com/shogo82148/goa-v1/dslengine"
	genschema "github.com/shogo82148/goa-v1/goagen/gen_schema"
	genswagger "github.com/shogo82148/goa-v1/goagen/gen_swagger"
	_ "github.com/shogo82148/goa-v1/goagen/gen_swagger/internal/design"
)

// validateSwagger validates that the given swagger object represents a valid Swagger spec.
func validateSwagger(swagger *genswagger.Swagger) {
	b, err := json.Marshal(swagger)
	Ω(err).ShouldNot(HaveOccurred())
	doc, err := loads.Analyzed(json.RawMessage(b), "")
	Ω(err).ShouldNot(HaveOccurred())
	Ω(doc).ShouldNot(BeNil())
}

// validateSwaggerWithFragments validates that the given swagger object represents a valid Swagger spec
// and contains fragments
func validateSwaggerWithFragments(swagger *genswagger.Swagger, fragments [][]byte) {
	b, err := json.Marshal(swagger)
	Ω(err).ShouldNot(HaveOccurred())
	doc, err := loads.Analyzed(json.RawMessage(b), "")
	Ω(err).ShouldNot(HaveOccurred())
	Ω(doc).ShouldNot(BeNil())
	for _, sub := range fragments {
		Ω(bytes.Contains(b, sub)).Should(BeTrue())
	}
}

var _ = Describe("New", func() {
	var swagger *genswagger.Swagger
	var newErr error

	BeforeEach(func() {
		swagger = nil
		newErr = nil
		dslengine.Reset()
		genschema.Definitions = make(map[string]*genschema.JSONSchema)
	})

	JustBeforeEach(func() {
		err := dslengine.Run()
		Ω(err).ShouldNot(HaveOccurred())
		swagger, newErr = genswagger.New(design.Design)
	})

	Context("with a valid API definition", func() {
		const (
			title        = "title"
			description  = "description"
			terms        = "terms"
			contactEmail = "contactEmail@goa.design"
			contactName  = "contactName"
			contactURL   = "http://contactURL.com"
			license      = "license"
			licenseURL   = "http://licenseURL.com"
			host         = "host"
			scheme       = "https"
			basePath     = "/base"
			tag          = "tag"
			docDesc      = "doc description"
			docURL       = "http://docURL.com"
		)

		BeforeEach(func() {
			apidsl.API("test", func() {
				apidsl.Title(title)
				apidsl.Metadata("swagger:tag:" + tag)
				apidsl.Metadata("swagger:tag:"+tag+":desc", "Tag desc.")
				apidsl.Metadata("swagger:tag:"+tag+":url", "http://example.com/tag")
				apidsl.Metadata("swagger:tag:"+tag+":url:desc", "Huge docs")
				apidsl.Description(description)
				apidsl.TermsOfService(terms)
				apidsl.Contact(func() {
					apidsl.Email(contactEmail)
					apidsl.Name(contactName)
					apidsl.URL(contactURL)
				})
				apidsl.License(func() {
					apidsl.Name(license)
					apidsl.URL(licenseURL)
				})
				apidsl.Docs(func() {
					apidsl.Description(docDesc)
					apidsl.URL(docURL)
				})
				apidsl.Host(host)
				apidsl.Scheme(scheme)
				apidsl.BasePath(basePath)
			})
		})

		It("sets all the basic fields", func() {
			Ω(newErr).ShouldNot(HaveOccurred())
			Ω(swagger).Should(Equal(&genswagger.Swagger{
				Swagger: "2.0",
				Info: &genswagger.Info{
					Title:          title,
					Description:    description,
					TermsOfService: terms,
					Contact: &design.ContactDefinition{
						Name:  contactName,
						Email: contactEmail,
						URL:   contactURL,
					},
					License: &design.LicenseDefinition{
						Name: license,
						URL:  licenseURL,
					},
					Version: "",
				},
				Host:     host,
				BasePath: basePath,
				Schemes:  []string{"https"},
				Paths:    make(map[string]interface{}),
				Consumes: []string{"application/json", "application/xml", "application/gob", "application/x-gob"},
				Produces: []string{"application/json", "application/xml", "application/gob", "application/x-gob"},
				Tags: []*genswagger.Tag{{Name: tag, Description: "Tag desc.", ExternalDocs: &genswagger.ExternalDocs{
					URL: "http://example.com/tag", Description: "Huge docs",
				}}},
				ExternalDocs: &genswagger.ExternalDocs{
					Description: docDesc,
					URL:         docURL,
				},
			}))
		})

		It("serializes into valid swagger JSON", func() { validateSwagger(swagger) })

		Context("with base params", func() {
			const (
				basePath    = "/s/:strParam/i/:intParam/n/:numParam/b/:boolParam"
				strParam    = "strParam"
				intParam    = "intParam"
				numParam    = "numParam"
				boolParam   = "boolParam"
				queryParam  = "queryParam"
				description = "description"
				intMin      = 1.0
				floatMax    = 2.4
				enum1       = "enum1"
				enum2       = "enum2"
			)

			BeforeEach(func() {
				base := design.Design.DSLFunc
				design.Design.DSLFunc = func() {
					base()
					apidsl.BasePath(basePath)
					apidsl.Params(func() {
						apidsl.Param(strParam, design.String, func() {
							apidsl.Description(description)
							apidsl.Format("email")
						})
						apidsl.Param(intParam, design.Integer, func() {
							apidsl.Minimum(intMin)
						})
						apidsl.Param(numParam, design.Number, func() {
							apidsl.Maximum(floatMax)
						})
						apidsl.Param(boolParam, design.Boolean)
						apidsl.Param(queryParam, func() {
							apidsl.Enum(enum1, enum2)
						})
					})
				}
			})

			It("sets the BasePath and Parameters fields", func() {
				Ω(newErr).ShouldNot(HaveOccurred())
				Ω(swagger.BasePath).Should(Equal(basePath))
				Ω(swagger.Parameters).Should(HaveLen(5))
				Ω(swagger.Parameters[strParam]).ShouldNot(BeNil())
				Ω(swagger.Parameters[strParam].Name).Should(Equal(strParam))
				Ω(swagger.Parameters[strParam].In).Should(Equal("path"))
				Ω(swagger.Parameters[strParam].Description).Should(Equal("description"))
				Ω(swagger.Parameters[strParam].Required).Should(BeTrue())
				Ω(swagger.Parameters[strParam].Type).Should(Equal("string"))
				Ω(swagger.Parameters[strParam].Format).Should(Equal("email"))
				Ω(swagger.Parameters[intParam]).ShouldNot(BeNil())
				Ω(swagger.Parameters[intParam].Name).Should(Equal(intParam))
				Ω(swagger.Parameters[intParam].In).Should(Equal("path"))
				Ω(swagger.Parameters[intParam].Required).Should(BeTrue())
				Ω(swagger.Parameters[intParam].Type).Should(Equal("integer"))
				Ω(*swagger.Parameters[intParam].Minimum).Should(Equal(intMin))
				Ω(swagger.Parameters[numParam]).ShouldNot(BeNil())
				Ω(swagger.Parameters[numParam].Name).Should(Equal(numParam))
				Ω(swagger.Parameters[numParam].In).Should(Equal("path"))
				Ω(swagger.Parameters[numParam].Required).Should(BeTrue())
				Ω(swagger.Parameters[numParam].Type).Should(Equal("number"))
				Ω(*swagger.Parameters[numParam].Maximum).Should(Equal(floatMax))
				Ω(swagger.Parameters[boolParam]).ShouldNot(BeNil())
				Ω(swagger.Parameters[boolParam].Name).Should(Equal(boolParam))
				Ω(swagger.Parameters[boolParam].In).Should(Equal("path"))
				Ω(swagger.Parameters[boolParam].Required).Should(BeTrue())
				Ω(swagger.Parameters[boolParam].Type).Should(Equal("boolean"))
				Ω(swagger.Parameters[queryParam]).ShouldNot(BeNil())
				Ω(swagger.Parameters[queryParam].Name).Should(Equal(queryParam))
				Ω(swagger.Parameters[queryParam].In).Should(Equal("query"))
				Ω(swagger.Parameters[queryParam].Type).Should(Equal("string"))
				Ω(swagger.Parameters[queryParam].Enum).Should(Equal([]interface{}{enum1, enum2}))
			})

			It("serializes into valid swagger JSON", func() { validateSwagger(swagger) })
		})

		Context("with required payload", func() {
			BeforeEach(func() {
				p := apidsl.Type("RequiredPayload", func() {
					apidsl.Member("m1", design.String)
				})
				apidsl.Resource("res", func() {
					apidsl.Action("act", func() {
						apidsl.Routing(
							apidsl.PUT("/"),
						)
						apidsl.Payload(p)
					})
				})
			})

			It("serializes into valid swagger JSON", func() {
				validateSwaggerWithFragments(swagger, [][]byte{
					[]byte(`"required":true`),
				})
			})
		})

		Context("with a payload of type Any", func() {
			BeforeEach(func() {
				apidsl.Resource("res", func() {
					apidsl.Action("act", func() {
						apidsl.Routing(
							apidsl.PUT("/"),
						)
						apidsl.Payload(design.Any, func() {
							apidsl.Example("example")
						})
					})
				})
			})

			It("serializes into valid swagger JSON", func() {
				validateSwaggerWithFragments(swagger, [][]byte{
					[]byte(`"ActResPayload":{"title":"ActResPayload","example":"example"}`),
				})
			})

		})

		Context("with optional payload", func() {
			BeforeEach(func() {
				p := apidsl.Type("OptionalPayload", func() {
					apidsl.Member("m1", design.String)
				})
				apidsl.Resource("res", func() {
					apidsl.Action("act", func() {
						apidsl.Routing(
							apidsl.PUT("/"),
						)
						apidsl.OptionalPayload(p)
					})
				})
			})

			It("serializes into valid swagger JSON", func() {
				validateSwaggerWithFragments(swagger, [][]byte{
					[]byte(`"required":false`),
				})
			})

		})

		Context("with multipart/form-data payload", func() {
			BeforeEach(func() {
				f := apidsl.Type("MultipartPayload", func() {
					apidsl.Attribute("image", design.File, "Binary image data")
				})
				apidsl.Resource("res", func() {
					apidsl.Action("act", func() {
						apidsl.Routing(
							apidsl.PUT("/"),
						)
						apidsl.MultipartForm()
						apidsl.Payload(f)
					})
				})
			})

			It("does not modify the API level consumes", func() {
				Ω(newErr).ShouldNot(HaveOccurred())
				Ω(swagger.Consumes).Should(HaveLen(4))
				Ω(swagger.Consumes).Should(ConsistOf("application/json", "application/xml", "application/gob", "application/x-gob"))
			})

			It("adds an Action level consumes for multipart/form-data", func() {
				Ω(newErr).ShouldNot(HaveOccurred())
				Ω(swagger.Paths).Should(HaveLen(1))
				Ω(swagger.Paths["/"]).ShouldNot(BeNil())

				a := swagger.Paths["/"].(*genswagger.Path)
				Ω(a.Put).ShouldNot(BeNil())
				cs := a.Put.Consumes
				Ω(cs).Should(HaveLen(1))
				Ω(cs[0]).Should(Equal("multipart/form-data"))
			})

			It("adds an File parameter", func() {
				Ω(newErr).ShouldNot(HaveOccurred())
				Ω(swagger.Paths).Should(HaveLen(1))
				Ω(swagger.Paths["/"]).ShouldNot(BeNil())

				a := swagger.Paths["/"].(*genswagger.Path)
				Ω(a.Put).ShouldNot(BeNil())
				ps := a.Put.Parameters
				Ω(ps).Should(HaveLen(1))
				Ω(ps[0]).Should(Equal(&genswagger.Parameter{In: "formData", Name: "image", Type: "file", Description: "Binary image data", Required: false}))
			})

			It("serializes into valid swagger JSON", func() { validateSwagger(swagger) })
		})

		Context("with recursive payload", func() {
			BeforeEach(func() {
				p := apidsl.Type("RecursivePayload", func() {
					apidsl.Member("m1", "RecursivePayload")
					apidsl.Member("m2", apidsl.ArrayOf("RecursivePayload"))
					apidsl.Member("m3", apidsl.HashOf(design.String, "RecursivePayload"))
					apidsl.Member("m4", func() {
						apidsl.Member("m5", design.String)
						apidsl.Member("m6", "RecursivePayload")
					})
				})
				apidsl.Resource("res", func() {
					apidsl.Action("act", func() {
						apidsl.Routing(
							apidsl.PUT("/"),
						)
						apidsl.Payload(p)
					})
				})
			})

			It("serializes into valid swagger JSON", func() { validateSwagger(swagger) })
		})

		Context("with zero value validations", func() {
			const (
				intParam = "intParam"
				numParam = "numParam"
				strParam = "strParam"
				intMin   = 0.0
				floatMax = 0.0
			)

			BeforeEach(func() {
				PayloadWithZeroValueValidations := apidsl.Type("PayloadWithZeroValueValidations", func() {
					apidsl.Attribute(strParam, design.String, func() {
						apidsl.MinLength(0)
						apidsl.MaxLength(0)
					})
				})
				apidsl.Resource("res", func() {
					apidsl.Action("act", func() {
						apidsl.Routing(
							apidsl.PUT("/"),
						)
						apidsl.Params(func() {
							apidsl.Param(intParam, design.Integer, func() {
								apidsl.Minimum(intMin)
							})
							apidsl.Param(numParam, design.Number, func() {
								apidsl.Maximum(floatMax)
							})
						})
						apidsl.Payload(PayloadWithZeroValueValidations)
					})
				})
			})

			It("serializes into valid swagger JSON", func() {
				validateSwaggerWithFragments(swagger, [][]byte{
					// payload
					[]byte(`"minLength":0`),
					[]byte(`"maxLength":0`),
					// param
					[]byte(`"minimum":0`),
					[]byte(`"maximum":0`),
				})
			})
		})

		Context("with minItems and maxItems validations in payload's attribute", func() {
			const (
				arrParam = "arrParam"
				minVal   = 0
				maxVal   = 42
			)

			BeforeEach(func() {
				PayloadWithValidations := apidsl.Type("Payload", func() {
					apidsl.Attribute(arrParam, apidsl.ArrayOf(design.String), func() {
						apidsl.MinLength(minVal)
						apidsl.MaxLength(maxVal)
					})
				})
				apidsl.Resource("res", func() {
					apidsl.Action("act", func() {
						apidsl.Routing(
							apidsl.PUT("/"),
						)
						apidsl.Payload(PayloadWithValidations)
					})
				})
			})

			It("serializes into valid swagger JSON", func() {
				validateSwaggerWithFragments(swagger, [][]byte{
					// payload
					[]byte(`"minItems":0`),
					[]byte(`"maxItems":42`),
				})
			})
		})

		Context("with minItems and maxItems validations in payload", func() {
			const (
				strParam = "strParam"
				minVal   = 0
				maxVal   = 42
			)

			BeforeEach(func() {
				PayloadWithValidations := apidsl.Type("Payload", func() {
					apidsl.Attribute(strParam, design.String)
				})
				apidsl.Resource("res", func() {
					apidsl.Action("act", func() {
						apidsl.Routing(
							apidsl.PUT("/"),
						)
						apidsl.Payload(apidsl.ArrayOf(PayloadWithValidations), func() {
							apidsl.MinLength(minVal)
							apidsl.MaxLength(maxVal)
						})
					})
				})
			})

			It("serializes into valid swagger JSON", func() {
				validateSwaggerWithFragments(swagger, [][]byte{
					// payload
					[]byte(`"minItems":0`),
					[]byte(`"maxItems":42`),
				})
			})
		})

		Context("with response templates", func() {
			const okName = "OK"
			const okDesc = "OK description"
			const notFoundName = "NotFound"
			const notFoundDesc = "NotFound description"
			const notFoundMt = "application/json"
			const headerName = "headerName"

			BeforeEach(func() {
				account := apidsl.MediaType("application/vnd.goa.test.account", func() {
					apidsl.Description("Account")
					apidsl.Attributes(func() {
						apidsl.Attribute("id", design.Integer)
						apidsl.Attribute("href", design.String)
					})
					apidsl.View("default", func() {
						apidsl.Attribute("id")
						apidsl.Attribute("href")
					})
					apidsl.View("link", func() {
						apidsl.Attribute("id")
						apidsl.Attribute("href")
					})
				})
				mt := apidsl.MediaType("application/vnd.goa.test.bottle", func() {
					apidsl.Description("A bottle of wine")
					apidsl.Attributes(func() {
						apidsl.Attribute("id", design.Integer, "ID of bottle")
						apidsl.Attribute("href", design.String, "API href of bottle")
						apidsl.Attribute("account", account, "Owner account")
						apidsl.Links(func() {
							apidsl.Link("account") // Defines a link to the Account media type
						})
						apidsl.Required("id", "href")
					})
					apidsl.View("default", func() {
						apidsl.Attribute("id")
						apidsl.Attribute("href")
						apidsl.Attribute("links") // Default view renders links
					})
					apidsl.View("extended", func() {
						apidsl.Attribute("id")
						apidsl.Attribute("href")
						apidsl.Attribute("account") // Extended view renders account inline
						apidsl.Attribute("links")   // Extended view also renders links
					})
				})
				base := design.Design.DSLFunc
				design.Design.DSLFunc = func() {
					base()
					apidsl.ResponseTemplate(okName, func() {
						apidsl.Description(okDesc)
						apidsl.Status(404)
						apidsl.Media(mt)
						apidsl.Headers(func() {
							apidsl.Header(headerName, func() {
								apidsl.Format("hostname")
							})
						})
					})
					apidsl.ResponseTemplate(notFoundName, func() {
						apidsl.Description(notFoundDesc)
						apidsl.Status(404)

						apidsl.Media(notFoundMt)
					})
				}
			})

			It("sets the Responses fields", func() {
				Ω(newErr).ShouldNot(HaveOccurred())
				Ω(swagger.Responses).Should(HaveLen(2))
				Ω(swagger.Responses[notFoundName]).ShouldNot(BeNil())
				Ω(swagger.Responses[notFoundName].Description).Should(Equal(notFoundDesc))
				Ω(swagger.Responses[okName]).ShouldNot(BeNil())
				Ω(swagger.Responses[okName].Description).Should(Equal(okDesc))
			})

			It("serializes into valid swagger JSON", func() { validateSwagger(swagger) })
		})

		Context("with resources", func() {
			var (
				minLength1  = 1
				maxLength10 = 10
				minimum_2   = -2.0
				maximum2    = 2.0
				minItems1   = 1
				maxItems5   = 5
			)
			BeforeEach(func() {
				Country := apidsl.MediaType("application/vnd.goa.example.origin", func() {
					apidsl.Description("Origin of bottle")
					apidsl.Attributes(func() {
						apidsl.Attribute("id")
						apidsl.Attribute("href")
						apidsl.Attribute("country")
					})
					apidsl.View("default", func() {
						apidsl.Attribute("id")
						apidsl.Attribute("href")
						apidsl.Attribute("country")
					})
					apidsl.View("tiny", func() {
						apidsl.Attribute("id")
					})
				})
				BottleMedia := apidsl.MediaType("application/vnd.goa.example.bottle", func() {
					apidsl.Description("A bottle of wine")
					apidsl.Attributes(func() {
						apidsl.Attribute("id", design.Integer, "ID of bottle")
						apidsl.Attribute("href", design.String, "API href of bottle")
						apidsl.Attribute("origin", Country, "Details on wine origin")
						apidsl.Links(func() {
							apidsl.Link("origin", "tiny")
						})
						apidsl.Required("id", "href")
					})
					apidsl.View("default", func() {
						apidsl.Attribute("id")
						apidsl.Attribute("href")
						apidsl.Attribute("links")
					})
					apidsl.View("extended", func() {
						apidsl.Attribute("id")
						apidsl.Attribute("href")
						apidsl.Attribute("origin")
						apidsl.Attribute("links")
					})
				})
				UpdatePayload := apidsl.Type("UpdatePayload", func() {
					apidsl.Description("Type of create and upload action payloads")
					apidsl.Attribute("name", design.String, "name of bottle")
					apidsl.Attribute("origin", Country, "Details on wine origin")
					apidsl.Required("name")
				})
				apidsl.Resource("res", func() {
					apidsl.Metadata("swagger:tag:res")
					apidsl.Description("A wine bottle")
					apidsl.DefaultMedia(BottleMedia)
					apidsl.BasePath("/bottles")
					apidsl.UseTrait("Authenticated")

					apidsl.Action("Update", func() {
						apidsl.Metadata("swagger:tag:Update")
						apidsl.Metadata("swagger:summary", "a summary")
						apidsl.Description("Update account")
						apidsl.Docs(func() {
							apidsl.Description("docs")
							apidsl.URL("http://cellarapi.com/docs/actions/update")
						})
						apidsl.Routing(
							apidsl.PUT("/:id"),
							apidsl.PUT("//orgs/:org/accounts/:id"),
						)
						apidsl.Params(func() {
							apidsl.Param("org", design.String)
							apidsl.Param("id", design.Integer)
							apidsl.Param("sort", func() {
								apidsl.Enum("asc", "desc")
							})
						})
						apidsl.Headers(func() {
							apidsl.Header("Authorization", design.String)
							apidsl.Header("X-Account", design.Integer)
							apidsl.Header("OptionalBoolWithDefault", design.Boolean, "defaults true", func() {
								apidsl.Default(true)
							})
							apidsl.Header("OptionalRegex", design.String, func() {
								apidsl.Pattern(`[a-z]\d+`)
								apidsl.MinLength(minLength1)
								apidsl.MaxLength(maxLength10)
							})
							apidsl.Header("OptionalInt", design.Integer, func() {
								apidsl.Minimum(minimum_2)
								apidsl.Maximum(maximum2)
							})
							apidsl.Header("OptionalArray", apidsl.ArrayOf(design.String), func() {
								// interpreted as MinItems & MaxItems:
								apidsl.MinLength(minItems1)
								apidsl.MaxLength(maxItems5)
							})
							apidsl.Header("OverrideRequiredHeader")
							apidsl.Header("OverrideOptionalHeader")
							apidsl.Required("Authorization", "X-Account", "OverrideOptionalHeader")
						})
						apidsl.Payload(UpdatePayload)
						apidsl.Response(design.OK, func() {
							apidsl.Media(apidsl.CollectionOf(BottleMedia), "extended")
						})
						apidsl.Response(design.NoContent)
						apidsl.Response(design.NotFound, design.ErrorMedia)
						apidsl.Response(design.BadRequest, design.ErrorMedia)
					})

					apidsl.Action("hidden", func() {
						apidsl.Description("Does not show up in Swagger spec")
						apidsl.Metadata("swagger:generate", "false")
						apidsl.Routing(apidsl.GET("/hidden"))
						apidsl.Response(design.OK)
					})
				})
				base := design.Design.DSLFunc
				design.Design.DSLFunc = func() {
					base()
					apidsl.Trait("Authenticated", func() {
						apidsl.Headers(func() {
							apidsl.Header("header")
							apidsl.Header("OverrideRequiredHeader", design.String, "to be overridden in Action and not marked Required")
							apidsl.Header("OverrideOptionalHeader", design.String, "to be overridden in Action and marked Required")
							apidsl.Header("OptionalResourceHeaderWithEnum", func() {
								apidsl.Enum("a", "b")
							})
							apidsl.Required("header", "OverrideRequiredHeader")
						})
					})
				}
			})

			It("sets the Path fields", func() {
				Ω(newErr).ShouldNot(HaveOccurred())
				Ω(swagger.Paths).Should(HaveLen(2))
				Ω(swagger.Paths["/orgs/{org}/accounts/{id}"]).ShouldNot(BeNil())
				a := swagger.Paths["/orgs/{org}/accounts/{id}"].(*genswagger.Path)
				Ω(a.Put).ShouldNot(BeNil())
				ps := a.Put.Parameters
				Ω(ps).Should(HaveLen(14))
				// check Headers in detail
				Ω(ps[3]).Should(Equal(&genswagger.Parameter{In: "header", Name: "Authorization", Type: "string", Required: true}))
				Ω(ps[4]).Should(Equal(&genswagger.Parameter{In: "header", Name: "OptionalArray", Type: "array", CollectionFormat: "multi",
					Items: &genswagger.Items{Type: "string"}, MinItems: &minItems1, MaxItems: &maxItems5}))
				Ω(ps[5]).Should(Equal(&genswagger.Parameter{In: "header", Name: "OptionalBoolWithDefault", Type: "boolean",
					Description: "defaults true", Default: true}))
				Ω(ps[6]).Should(Equal(&genswagger.Parameter{In: "header", Name: "OptionalInt", Type: "integer", Minimum: &minimum_2, Maximum: &maximum2}))
				Ω(ps[7]).Should(Equal(&genswagger.Parameter{In: "header", Name: "OptionalRegex", Type: "string",
					Pattern: `[a-z]\d+`, MinLength: &minLength1, MaxLength: &maxLength10}))
				Ω(ps[8]).Should(Equal(&genswagger.Parameter{In: "header", Name: "OptionalResourceHeaderWithEnum", Type: "string",
					Enum: []interface{}{"a", "b"}}))
				Ω(ps[9]).Should(Equal(&genswagger.Parameter{In: "header", Name: "OverrideOptionalHeader", Type: "string", Required: true}))
				Ω(ps[10]).Should(Equal(&genswagger.Parameter{In: "header", Name: "OverrideRequiredHeader", Type: "string", Required: true}))
				Ω(ps[11]).Should(Equal(&genswagger.Parameter{In: "header", Name: "X-Account", Type: "integer", Required: true}))
				Ω(ps[12]).Should(Equal(&genswagger.Parameter{In: "header", Name: "header", Type: "string", Required: true}))
				Ω(swagger.Paths["/base/bottles/{id}"]).ShouldNot(BeNil())
				b := swagger.Paths["/base/bottles/{id}"].(*genswagger.Path)
				Ω(b.Put).ShouldNot(BeNil())
				Ω(b.Put.Parameters).Should(HaveLen(14))
				Ω(b.Put.Produces).Should(Equal([]string{"application/vnd.goa.error", "application/vnd.goa.example.bottle; type=collection"}))
			})

			It("should set the inherited tag and the action tag", func() {
				tags := []string{"res", "Update"}
				a := swagger.Paths["/orgs/{org}/accounts/{id}"].(*genswagger.Path)
				Ω(a.Put).ShouldNot(BeNil())
				Ω(a.Put.Tags).Should(Equal(tags))
				b := swagger.Paths["/base/bottles/{id}"].(*genswagger.Path)
				Ω(b.Put.Tags).Should(Equal(tags))
			})

			It("sets the summary from the summary tag", func() {
				a := swagger.Paths["/orgs/{org}/accounts/{id}"].(*genswagger.Path)
				Ω(a.Put.Summary).Should(Equal("a summary"))
			})

			It("generates the media type collection schema", func() {
				Ω(swagger.Definitions).Should(HaveLen(7))
				Ω(swagger.Definitions).Should(HaveKey("GoaExampleBottleExtendedCollection"))
			})

			It("serializes into valid swagger JSON", func() { validateSwagger(swagger) })
		})

		Context("with metadata", func() {
			const gat = "gat"
			const extension = `{"foo":"bar"}`
			const stringExtension = "foo"

			var (
				unmarshaled map[string]interface{}
				_           = json.Unmarshal([]byte(extension), &unmarshaled)
			)

			BeforeEach(func() {
				apidsl.Resource("res", func() {
					apidsl.Metadata("swagger:tag:res")
					apidsl.Metadata("struct:tag:json", "resource")
					apidsl.Metadata("swagger:extension:x-resource", extension)
					apidsl.Metadata("swagger:extension:x-string", stringExtension)
					apidsl.Action("act", func() {
						apidsl.Metadata("swagger:tag:Update")
						apidsl.Metadata("struct:tag:json", "action")
						apidsl.Metadata("swagger:extension:x-action", extension)
						apidsl.Security("password", func() {
							apidsl.Metadata("swagger:extension:x-security", extension)
						})
						apidsl.Routing(
							apidsl.PUT("/", func() {
								apidsl.Metadata("swagger:extension:x-put", extension)
							}),
						)
						apidsl.Params(func() {
							apidsl.Param("param", func() {
								apidsl.Metadata("swagger:extension:x-param", extension)
							})
						})
						apidsl.Response(design.NoContent, func() {
							apidsl.Metadata("swagger:extension:x-response", extension)
						})
					})
				})
				base := design.Design.DSLFunc
				design.Design.DSLFunc = func() {
					base()
					apidsl.Metadata("swagger:tag:" + gat)
					apidsl.Metadata("struct:tag:json", "api")
					apidsl.Metadata("swagger:extension:x-api", extension)
					apidsl.BasicAuthSecurity("password")
				}
			})

			It("should set the swagger object tags", func() {
				Ω(swagger.Tags).Should(HaveLen(2))
				tags := []*genswagger.Tag{
					{Name: gat, Description: "", ExternalDocs: nil, Extensions: map[string]interface{}{"x-api": unmarshaled}},
					{Name: tag, Description: "Tag desc.", ExternalDocs: &genswagger.ExternalDocs{URL: "http://example.com/tag", Description: "Huge docs"}, Extensions: map[string]interface{}{"x-api": unmarshaled}},
				}
				Ω(swagger.Tags).Should(Equal(tags))
			})

			It("should set the action tags", func() {
				p := swagger.Paths["/"].(*genswagger.Path)
				Ω(p.Put.Tags).Should(HaveLen(2))
				tags := []string{"res", "Update"}
				Ω(p.Put.Tags).Should(Equal(tags))
			})

			It("should set the swagger extensions", func() {
				Ω(swagger.Info.Extensions).Should(HaveLen(1))
				Ω(swagger.Info.Extensions["x-api"]).Should(Equal(unmarshaled))
				p := swagger.Paths["/"].(*genswagger.Path)
				Ω(p.Extensions).Should(HaveLen(1))
				Ω(p.Extensions["x-action"]).Should(Equal(unmarshaled))
				Ω(p.Put.Extensions).Should(HaveLen(1))
				Ω(p.Put.Extensions["x-put"]).Should(Equal(unmarshaled))
				Ω(p.Put.Parameters[0].Extensions).Should(HaveLen(1))
				Ω(p.Put.Parameters[0].Extensions["x-param"]).Should(Equal(unmarshaled))
				Ω(p.Put.Responses["204"].Extensions).Should(HaveLen(1))
				Ω(p.Put.Responses["204"].Extensions["x-response"]).Should(Equal(unmarshaled))
				Ω(swagger.Paths["x-resource"]).ShouldNot(BeNil())
				rs := swagger.Paths["x-resource"].(map[string]interface{})
				Ω(rs).Should(Equal(unmarshaled))
				rs2 := swagger.Paths["x-string"].(string)
				Ω(rs2).Should(Equal(stringExtension))
				Ω(swagger.SecurityDefinitions["password"].Extensions).Should(HaveLen(1))
				Ω(swagger.SecurityDefinitions["password"].Extensions["x-security"]).Should(Equal(unmarshaled))
			})

		})
	})
})
