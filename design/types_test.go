package design_test

import (
	"errors"
	"mime"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/goa-v1/design"
	"github.com/shogo82148/goa-v1/design/apidsl"
	"github.com/shogo82148/goa-v1/dslengine"
)

var _ = Describe("IsObject", func() {
	var dt design.DataType
	var isObject bool

	JustBeforeEach(func() {
		isObject = dt.IsObject()
	})

	Context("with a primitive", func() {
		BeforeEach(func() {
			dt = design.String
		})

		It("returns false", func() {
			Ω(isObject).Should(BeFalse())
		})
	})

	Context("with an array", func() {
		BeforeEach(func() {
			dt = &design.Array{ElemType: &design.AttributeDefinition{Type: design.String}}
		})

		It("returns false", func() {
			Ω(isObject).Should(BeFalse())
		})
	})

	Context("with a hash", func() {
		BeforeEach(func() {
			dt = &design.Hash{
				KeyType:  &design.AttributeDefinition{Type: design.String},
				ElemType: &design.AttributeDefinition{Type: design.String},
			}
		})

		It("returns false", func() {
			Ω(isObject).Should(BeFalse())
		})
	})

	Context("with a nil user type type", func() {
		BeforeEach(func() {
			dt = &design.UserTypeDefinition{AttributeDefinition: &design.AttributeDefinition{Type: nil}}
		})

		It("returns false", func() {
			Ω(isObject).Should(BeFalse())
		})
	})

	Context("with an object", func() {
		BeforeEach(func() {
			dt = design.Object{}
		})

		It("returns true", func() {
			Ω(isObject).Should(BeTrue())
		})
	})
})

var _ = Describe("Project", func() {
	var mt *design.MediaTypeDefinition
	var view string

	var projected *design.MediaTypeDefinition
	var links *design.UserTypeDefinition
	var prErr error

	JustBeforeEach(func() {
		design.ProjectedMediaTypes = make(map[string]*design.MediaTypeDefinition)
		projected, links, prErr = mt.Project(view)
	})

	Context("with a media type with a default and a tiny view", func() {
		BeforeEach(func() {
			mt = &design.MediaTypeDefinition{
				UserTypeDefinition: &design.UserTypeDefinition{
					AttributeDefinition: &design.AttributeDefinition{
						Type: design.Object{
							"att1": &design.AttributeDefinition{Type: design.Integer},
							"att2": &design.AttributeDefinition{Type: design.String},
						},
					},
					TypeName: "Foo",
				},
				Identifier: "vnd.application/foo",
				Views: map[string]*design.ViewDefinition{
					"default": {
						Name: "default",
						AttributeDefinition: &design.AttributeDefinition{
							Type: design.Object{
								"att1": &design.AttributeDefinition{Type: design.String},
								"att2": &design.AttributeDefinition{Type: design.String},
							},
						},
					},
					"tiny": {
						Name: "tiny",
						AttributeDefinition: &design.AttributeDefinition{
							Type: design.Object{
								"att2": &design.AttributeDefinition{Type: design.String},
							},
						},
					},
				},
			}
		})

		Context("using the empty view", func() {
			BeforeEach(func() {
				view = ""
			})

			It("returns an error", func() {
				Ω(prErr).Should(HaveOccurred())
			})
		})

		Context("using the default view", func() {
			BeforeEach(func() {
				view = "default"
			})

			It("returns a media type with an identifier view param", func() {
				Ω(prErr).ShouldNot(HaveOccurred())
				_, params, err := mime.ParseMediaType(projected.Identifier)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(params).Should(HaveKeyWithValue("view", "default"))
			})

			It("returns a media type with only a default view", func() {
				Ω(prErr).ShouldNot(HaveOccurred())
				Ω(projected.Views).Should(HaveLen(1))
				Ω(projected.Views).Should(HaveKey("default"))
			})

			It("returns a media type with the default view attributes", func() {
				Ω(prErr).ShouldNot(HaveOccurred())
				Ω(projected).ShouldNot(BeNil())
				Ω(projected.Type).Should(BeAssignableToTypeOf(design.Object{}))
				Ω(projected.Type.ToObject()).Should(HaveKey("att1"))
				att := projected.Type.ToObject()["att1"]
				Ω(att).ShouldNot(BeNil())
				Ω(att.Type).ShouldNot(BeNil())
				Ω(att.Type.Kind()).Should(Equal(design.IntegerKind))
			})
		})

		Context("using the tiny view", func() {
			BeforeEach(func() {
				view = "tiny"
			})

			It("returns a media type with an identifier view param", func() {
				Ω(prErr).ShouldNot(HaveOccurred())
				_, params, err := mime.ParseMediaType(projected.Identifier)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(params).Should(HaveKeyWithValue("view", "tiny"))
			})

			It("returns a media type with only a default view", func() {
				Ω(prErr).ShouldNot(HaveOccurred())
				Ω(projected.Views).Should(HaveLen(1))
				Ω(projected.Views).Should(HaveKey("default"))
			})

			It("returns a media type with the default view attributes", func() {
				Ω(prErr).ShouldNot(HaveOccurred())
				Ω(projected).ShouldNot(BeNil())
				Ω(projected.Type).Should(BeAssignableToTypeOf(design.Object{}))
				Ω(projected.Type.ToObject()).Should(HaveKey("att2"))
				att := projected.Type.ToObject()["att2"]
				Ω(att).ShouldNot(BeNil())
				Ω(att.Type).ShouldNot(BeNil())
				Ω(att.Type.Kind()).Should(Equal(design.StringKind))
			})

			Context("on a collection", func() {
				BeforeEach(func() {
					mt = apidsl.CollectionOf(design.Dup(mt))
					dslengine.Execute(mt.DSL(), mt)
					mt.GenerateExample(design.NewRandomGenerator(""), nil)
				})

				It("resets the example", func() {
					Ω(prErr).ShouldNot(HaveOccurred())
					Ω(projected).ShouldNot(BeNil())
					Ω(projected.Example).Should(BeNil())
				})
			})
		})

	})

	Context("with a media type with a links attribute", func() {
		BeforeEach(func() {
			mt = &design.MediaTypeDefinition{
				UserTypeDefinition: &design.UserTypeDefinition{
					AttributeDefinition: &design.AttributeDefinition{
						Type: design.Object{
							"att1":  &design.AttributeDefinition{Type: design.Integer},
							"links": &design.AttributeDefinition{Type: design.String},
						},
					},
					TypeName: "Foo",
				},
				Identifier: "vnd.application/foo",
				Views: map[string]*design.ViewDefinition{
					"default": {
						Name: "default",
						AttributeDefinition: &design.AttributeDefinition{
							Type: design.Object{
								"att1":  &design.AttributeDefinition{Type: design.String},
								"links": &design.AttributeDefinition{Type: design.String},
							},
						},
					},
				},
			}
		})

		Context("using the default view", func() {
			BeforeEach(func() {
				view = "default"
			})

			It("uses the links attribute in the view", func() {
				Ω(prErr).ShouldNot(HaveOccurred())
				Ω(projected).ShouldNot(BeNil())
				Ω(projected.Type).Should(BeAssignableToTypeOf(design.Object{}))
				Ω(projected.Type.ToObject()).Should(HaveKey("links"))
				att := projected.Type.ToObject()["links"]
				Ω(att).ShouldNot(BeNil())
				Ω(att.Type).ShouldNot(BeNil())
				Ω(att.Type.Kind()).Should(Equal(design.StringKind))
			})
		})
	})

	Context("with media types with view attributes with a cyclical dependency", func() {
		const id = "vnd.application/MT1"
		const typeName = "Mt1"
		metadata := dslengine.MetadataDefinition{"foo": []string{"bar"}}

		BeforeEach(func() {
			dslengine.Reset()
			apidsl.API("test", func() {})
			mt = apidsl.MediaType(id, func() {
				apidsl.TypeName(typeName)
				apidsl.Attributes(func() {
					apidsl.Attribute("att", "vnd.application/MT2", func() {
						apidsl.Metadata("foo", "bar")
					})
				})
				apidsl.Links(func() {
					apidsl.Link("att", "default")
				})
				apidsl.View("default", func() {
					apidsl.Attribute("att")
					apidsl.Attribute("links")
				})
				apidsl.View("tiny", func() {
					apidsl.Attribute("att", func() {
						apidsl.View("tiny")
					})
				})
			})
			apidsl.MediaType("vnd.application/MT2", func() {
				apidsl.TypeName("Mt2")
				apidsl.Attributes(func() {
					apidsl.Attribute("att2", mt)
				})
				apidsl.Links(func() {
					apidsl.Link("att2", "default")
				})
				apidsl.View("default", func() {
					apidsl.Attribute("att2")
					apidsl.Attribute("links")
				})
				apidsl.View("tiny", func() {
					apidsl.Attribute("links")
				})
			})
			err := dslengine.Run()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(dslengine.Errors).ShouldNot(HaveOccurred())
		})

		Context("using the default view", func() {
			BeforeEach(func() {
				view = "default"
			})

			It("returns the projected media type with links", func() {
				Ω(prErr).ShouldNot(HaveOccurred())
				Ω(projected).ShouldNot(BeNil())
				Ω(projected.Type).Should(BeAssignableToTypeOf(design.Object{}))
				Ω(projected.Type.ToObject()).Should(HaveKey("att"))
				l := projected.Type.ToObject()["links"]
				Ω(l.Type.(*design.UserTypeDefinition).AttributeDefinition).Should(Equal(links.AttributeDefinition))
				Ω(links.Type.ToObject()).Should(HaveKey("att"))
				Ω(links.Type.ToObject()["att"].Metadata).Should(Equal(metadata))
			})
		})

		Context("using the tiny view", func() {
			BeforeEach(func() {
				view = "tiny"
			})

			It("returns the projected media type with links", func() {
				Ω(prErr).ShouldNot(HaveOccurred())
				Ω(projected).ShouldNot(BeNil())
				Ω(projected.Type).Should(BeAssignableToTypeOf(design.Object{}))
				Ω(projected.Type.ToObject()).Should(HaveKey("att"))
				att := projected.Type.ToObject()["att"]
				Ω(att.Type.ToObject()).Should(HaveKey("links"))
				Ω(att.Type.ToObject()).ShouldNot(HaveKey("att2"))
			})
		})
	})
})

var _ = Describe("UserTypes", func() {
	var (
		o         design.Object
		userTypes map[string]*design.UserTypeDefinition
	)

	JustBeforeEach(func() {
		userTypes = design.UserTypes(o)
	})

	Context("with an object not using user types", func() {
		BeforeEach(func() {
			o = design.Object{"foo": &design.AttributeDefinition{Type: design.String}}
		})

		It("returns nil", func() {
			Ω(userTypes).Should(BeNil())
		})
	})

	Context("with an object with an attribute using a user type", func() {
		var ut *design.UserTypeDefinition
		BeforeEach(func() {
			ut = &design.UserTypeDefinition{
				TypeName:            "foo",
				AttributeDefinition: &design.AttributeDefinition{Type: design.String},
			}

			o = design.Object{"foo": &design.AttributeDefinition{Type: ut}}
		})

		It("returns the user type", func() {
			Ω(userTypes).Should(HaveLen(1))
			Ω(userTypes[ut.TypeName]).Should(Equal(ut))
		})
	})

	Context("with an object with an attribute using recursive user types", func() {
		var ut, childut *design.UserTypeDefinition

		BeforeEach(func() {
			childut = &design.UserTypeDefinition{
				TypeName:            "child",
				AttributeDefinition: &design.AttributeDefinition{Type: design.String},
			}
			child := design.Object{"child": &design.AttributeDefinition{Type: childut}}
			ut = &design.UserTypeDefinition{
				TypeName:            "parent",
				AttributeDefinition: &design.AttributeDefinition{Type: child},
			}

			o = design.Object{"foo": &design.AttributeDefinition{Type: ut}}
		})

		It("returns the user types", func() {
			Ω(userTypes).Should(HaveLen(2))
			Ω(userTypes[ut.TypeName]).Should(Equal(ut))
			Ω(userTypes[childut.TypeName]).Should(Equal(childut))
		})
	})
})

var _ = Describe("MediaTypeDefinition", func() {
	Describe("IterateViews", func() {
		var (
			m  *design.MediaTypeDefinition
			it design.ViewIterator

			iteratedViews []string
		)
		BeforeEach(func() {
			m = &design.MediaTypeDefinition{}

			// setup iterator that just accumulates view names into iteratedViews
			iteratedViews = []string{}
			it = func(v *design.ViewDefinition) error {
				iteratedViews = append(iteratedViews, v.Name)
				return nil
			}
		})
		It("works with empty", func() {
			Expect(m.Views).To(BeEmpty())
			Expect(m.IterateViews(it)).To(Succeed())
			Expect(iteratedViews).To(BeEmpty())
		})
		Context("with non-empty views map", func() {
			BeforeEach(func() {
				m.Views = map[string]*design.ViewDefinition{
					"d": {Name: "d"},
					"c": {Name: "c"},
					"a": {Name: "a"},
					"b": {Name: "b"},
				}
			})
			It("sorts views", func() {
				Expect(m.IterateViews(it)).To(Succeed())
				Expect(iteratedViews).To(Equal([]string{"a", "b", "c", "d"}))
			})
			It("propagates error", func() {
				errIterator := func(v *design.ViewDefinition) error {
					if len(iteratedViews) > 2 {
						return errors.New("foo")
					}
					iteratedViews = append(iteratedViews, v.Name)
					return nil
				}
				Expect(m.IterateViews(errIterator)).To(MatchError("foo"))
				Expect(iteratedViews).To(Equal([]string{"a", "b", "c"}))
			})
		})
	})
})

var _ = Describe("Walk", func() {
	var target design.DataStructure
	var matchedName string
	var count int
	var matched bool

	counter := func(*design.AttributeDefinition) error {
		count++
		return nil
	}

	matcher := func(name string) func(*design.AttributeDefinition) error {
		done := errors.New("done")
		return func(att *design.AttributeDefinition) error {
			if u, ok := att.Type.(*design.UserTypeDefinition); ok {
				if u.TypeName == name {
					matched = true
					return done
				}
			} else if m, ok := att.Type.(*design.MediaTypeDefinition); ok {
				if m.TypeName == name {
					matched = true
					return done
				}
			}
			return nil
		}
	}

	BeforeEach(func() {
		matchedName = ""
		count = 0
		matched = false
	})

	JustBeforeEach(func() {
		target.Walk(counter)
		if matchedName != "" {
			target.Walk(matcher(matchedName))
		}
	})

	Context("with simple attribute", func() {
		BeforeEach(func() {
			target = &design.AttributeDefinition{Type: design.String}
		})

		It("walks it", func() {
			Ω(count).Should(Equal(1))
		})
	})

	Context("with an object attribute", func() {
		BeforeEach(func() {
			o := design.Object{"foo": &design.AttributeDefinition{Type: design.String}}
			target = &design.AttributeDefinition{Type: o}
		})

		It("walks it", func() {
			Ω(count).Should(Equal(2))
		})
	})

	Context("with an object attribute containing user types", func() {
		const typeName = "foo"
		BeforeEach(func() {
			matchedName = typeName
			at := &design.AttributeDefinition{Type: design.String}
			ut := &design.UserTypeDefinition{AttributeDefinition: at, TypeName: typeName}
			o := design.Object{"foo": &design.AttributeDefinition{Type: ut}}
			target = &design.AttributeDefinition{Type: o}
		})

		It("walks it", func() {
			Ω(count).Should(Equal(3))
			Ω(matched).Should(BeTrue())
		})
	})

	Context("with an object attribute containing recursive user types", func() {
		const typeName = "foo"
		BeforeEach(func() {
			matchedName = typeName
			co := design.Object{}
			at := &design.AttributeDefinition{Type: co}
			ut := &design.UserTypeDefinition{AttributeDefinition: at, TypeName: typeName}
			co["recurse"] = &design.AttributeDefinition{Type: ut}
			o := design.Object{"foo": &design.AttributeDefinition{Type: ut}}
			target = &design.AttributeDefinition{Type: o}
		})

		It("walks it", func() {
			Ω(count).Should(Equal(4))
			Ω(matched).Should(BeTrue())
		})
	})
})

var _ = Describe("Finalize", func() {
	BeforeEach(func() {
		dslengine.Reset()
		apidsl.MediaType("application/vnd.menu+json", func() {
			apidsl.Attributes(func() {
				apidsl.Attribute("name", design.String, "The name of an application")
				apidsl.Attribute("child", apidsl.CollectionOf("application/vnd.menu"))
			})

			apidsl.View("default", func() {
				apidsl.Attribute("name")
			})
		})
	})

	It("running the DSL should not loop indefinitely", func() {
		var mu sync.Mutex
		err := errors.New("infinite loop")
		go func() {
			err2 := dslengine.Run()
			mu.Lock()
			defer mu.Unlock()
			err = err2
		}()
		Eventually(func() error {
			mu.Lock()
			defer mu.Unlock()
			return err
		}).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("GenerateExample", func() {

	Context("Given a UUID", func() {
		It("generates a string example", func() {
			rand := design.NewRandomGenerator("foo")
			Ω(design.UUID.GenerateExample(rand, nil)).Should(BeAssignableToTypeOf("foo"))
		})
	})

	Context("Given a Hash keyed by UUIDs", func() {
		var h *design.Hash
		BeforeEach(func() {
			h = &design.Hash{
				KeyType:  &design.AttributeDefinition{Type: design.UUID},
				ElemType: &design.AttributeDefinition{Type: design.String},
			}
		})
		It("generates a serializable example", func() {
			rand := design.NewRandomGenerator("foo")
			Ω(h.GenerateExample(rand, nil)).Should(BeAssignableToTypeOf(map[string]string{"foo": "bar"}))
		})
	})
})
