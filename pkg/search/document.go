package search

import (
	"fmt"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

type DocumentOption func(*Document) error

func BuildDocument(name string, opts ...DocumentOption) *Document {
	doc := &Document{
		Name:    name,
		Fields:  make(map[string]FType),
		SubDocs: make(map[string]*Document),
	}
	for _, opt := range opts {
		if err := opt(doc); err != nil {
			panic(err)
		}
	}
	return doc
}

func WithField(name string, fType FType) DocumentOption {
	return func(doc *Document) error {
		if _, hasField := doc.Fields[name]; hasField {
			return fmt.Errorf("duplicate field: %s", name)
		}
		if _, hasDoc := doc.SubDocs[name]; hasDoc {
			return fmt.Errorf("field overlaps subdocument field: %s", name)
		}
		doc.Fields[name] = fType
		return nil
	}
}

func WithSubDocument(subDoc *Document) DocumentOption {
	return func(doc *Document) error {
		if _, hasField := doc.Fields[subDoc.Name]; hasField {
			return fmt.Errorf("document overlaps field: %s", subDoc.Name)
		}
		if _, hasDoc := doc.SubDocs[subDoc.Name]; hasDoc {
			return fmt.Errorf("duplicate subdocument: %s", subDoc.Name)
		}
		doc.SubDocs[subDoc.Name] = subDoc
		return nil
	}
}

type Document struct {
	Name    string
	Fields  map[string]FType
	SubDocs map[string]*Document
}

func (d *Document) toBleveMapping() *mapping.DocumentMapping {
	docMapping := bleve.NewDocumentMapping()
	for fName, fType := range d.Fields {
		fm := fType.toBleveMapping()
		docMapping.AddFieldMappingsAt(fName, fm)
	}
	for dName, subDoc := range d.SubDocs {
		sd := subDoc.toBleveMapping()
		docMapping.AddSubDocumentMapping(dName, sd)
	}
	return docMapping
}
