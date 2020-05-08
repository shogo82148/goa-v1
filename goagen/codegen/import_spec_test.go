package codegen_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/goa-v1/design"
	"github.com/shogo82148/goa-v1/dslengine"
	"github.com/shogo82148/goa-v1/goagen/codegen"
)

var _ = Describe("AttributeImports", func() {
	Context("given an attribute definition with fields", func() {
		var att *design.AttributeDefinition
		var st string
		var object design.Object

		Context("of object", func() {

			It("produces the import slice", func() {
				var imports []*codegen.ImportSpec
				object = design.Object{
					"foo": &design.AttributeDefinition{Type: design.String},
					"bar": &design.AttributeDefinition{Type: design.Integer},
				}
				object["foo"].Metadata = dslengine.MetadataDefinition{
					"struct:field:type": []string{"json.RawMessage", "encoding/json"},
				}
				att = new(design.AttributeDefinition)
				att.Type = object
				imports = codegen.AttributeImports(att, imports, nil)

				i := []*codegen.ImportSpec{{
					Path: "encoding/json",
				},
				}
				st = i[0].Path

				Ω(st).Should(Equal(imports[0].Path))
			})
		})

		Context("of recursive object", func() {

			It("produces the import slice", func() {
				var imports []*codegen.ImportSpec
				o := design.Object{
					"foo": &design.AttributeDefinition{Type: design.String},
				}
				o["foo"].Metadata = dslengine.MetadataDefinition{
					"struct:field:type": []string{"json.RawMessage", "encoding/json"},
				}
				child := &design.AttributeDefinition{Type: o}

				po := design.Object{"child": child}
				po["child"].Metadata = dslengine.MetadataDefinition{
					"struct:field:type": []string{"json.RawMessage", "encoding/json"},
				}
				parent := &design.AttributeDefinition{Type: po}

				o["parent"] = parent

				att = new(design.AttributeDefinition)
				att.Type = po
				imports = codegen.AttributeImports(att, imports, nil)

				i := []*codegen.ImportSpec{{
					Path: "encoding/json",
				},
				}
				st = i[0].Path
				l := len(imports)

				Ω(st).Should(Equal(imports[0].Path))
				Ω(l).Should(Equal(1))
			})
		})

		Context("of hash", func() {

			It("produces the import slice", func() {
				var imports []*codegen.ImportSpec
				elemType := &design.AttributeDefinition{Type: design.Integer}
				elemType.Metadata = dslengine.MetadataDefinition{
					"struct:field:type": []string{"json.RawMessage", "encoding/json"},
				}
				keyType := &design.AttributeDefinition{Type: design.Integer}
				elemType.Metadata = dslengine.MetadataDefinition{
					"struct:field:type": []string{"json.RawMessage", "encoding/json"},
				}
				hash := &design.Hash{KeyType: keyType, ElemType: elemType}

				att = new(design.AttributeDefinition)
				att.Type = hash
				imports = codegen.AttributeImports(att, imports, nil)

				i := []*codegen.ImportSpec{{
					Path: "encoding/json",
				},
				}
				st = i[0].Path
				l := len(imports)

				Ω(st).Should(Equal(imports[0].Path))
				Ω(l).Should(Equal(1))
			})
		})

		Context("of array", func() {
			It("produces the import slice", func() {
				var imports []*codegen.ImportSpec
				elemType := &design.AttributeDefinition{Type: design.Integer}
				elemType.Metadata = dslengine.MetadataDefinition{
					"struct:field:type": []string{"json.RawMessage", "encoding/json"},
				}
				array := &design.Array{ElemType: elemType}

				att = new(design.AttributeDefinition)
				att.Type = array
				imports = codegen.AttributeImports(att, imports, nil)

				i := []*codegen.ImportSpec{{
					Path: "encoding/json",
				},
				}
				st = i[0].Path

				Ω(st).Should(Equal(imports[0].Path))
			})
		})

		Context("of UserTypeDefinition", func() {

			It("produces the import slice", func() {
				var imports []*codegen.ImportSpec
				object = design.Object{
					"bar": &design.AttributeDefinition{Type: design.String},
				}
				object["bar"].Metadata = dslengine.MetadataDefinition{
					"struct:field:type": []string{"json.RawMessage", "encoding/json"},
				}

				u := &design.UserTypeDefinition{
					AttributeDefinition: &design.AttributeDefinition{Type: object},
				}

				att = u.AttributeDefinition
				imports = codegen.AttributeImports(att, imports, nil)

				i := []*codegen.ImportSpec{{
					Path: "encoding/json",
				},
				}
				st = i[0].Path

				Ω(st).Should(Equal(imports[0].Path))
			})
		})

		Context("of MediaTypeDefinition", func() {
			It("produces the import slice", func() {
				var imports []*codegen.ImportSpec
				elemType := &design.AttributeDefinition{Type: design.Integer}
				elemType.Metadata = dslengine.MetadataDefinition{
					"struct:field:type": []string{"json.RawMessage", "encoding/json"},
				}
				array := &design.Array{ElemType: elemType}
				u := &design.UserTypeDefinition{
					AttributeDefinition: &design.AttributeDefinition{Type: array},
				}
				m := &design.MediaTypeDefinition{
					UserTypeDefinition: u,
				}

				att = m.AttributeDefinition
				imports = codegen.AttributeImports(att, imports, nil)

				i := []*codegen.ImportSpec{{
					Path: "encoding/json",
				},
				}
				st = i[0].Path

				Ω(st).Should(Equal(imports[0].Path))
			})
		})
	})
})
