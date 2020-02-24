package search

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

type IndexOption func(*Index) error

func BuildIndex(name string, opts ...IndexOption) *Index {
	index := &Index{
		Name: name,
	}
	for _, opt := range opts {
		if err := opt(index); err != nil {
			panic(err)
		}
	}
	return index
}

func WithDocument(doc *Document) IndexOption {
	return func(ind *Index) error {
		ind.Documents = append(ind.Documents, doc)
		return nil
	}
}

type Index struct {
	Name      string
	Documents []*Document
}

func (i *Index) ToBleveMapping() (mapping.IndexMapping, error) {
	indexMapping := bleve.NewIndexMapping()
	addCustomAnalyzers(indexMapping)
	for _, doc := range i.Documents {
		dm := doc.toBleveMapping()
		indexMapping.AddDocumentMapping(doc.Name, dm)
	}
	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"
	return indexMapping, nil
}
