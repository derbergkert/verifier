package search

import (
	"errors"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/tokenizer/whitespace"
	"github.com/blevesearch/bleve/mapping"
)

var (
	once sync.Once

	englishTextFieldMapping *mapping.FieldMapping
	keywordFieldMapping     *mapping.FieldMapping
	keywordListFieldMapping *mapping.FieldMapping
	numFieldMapping         *mapping.FieldMapping
	dateFieldMapping        *mapping.FieldMapping
	boolFieldMapping        *mapping.FieldMapping
)

func addCustomAnalyzers(im *mapping.IndexMappingImpl) {
	if err := im.AddCustomAnalyzer("keywordList", map[string]interface{}{
		"type":          custom.Name,
		"tokenizer":     whitespace.Name,
		"token_filters": []interface{}{lowercase.Name},
	}); err != nil {
		panic(err)
	}
}

func initializeMappings() {
	// test field mapping.
	englishTextFieldMapping = bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName
	englishTextFieldMapping.IncludeInAll = false
	englishTextFieldMapping.IncludeTermVectors = false
	englishTextFieldMapping.Store = false

	// Keyword field mapping.
	keywordFieldMapping = bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name
	keywordFieldMapping.IncludeInAll = false
	keywordFieldMapping.IncludeTermVectors = false
	keywordFieldMapping.Store = false

	// KeywordList field mapping.
	keywordListFieldMapping = bleve.NewTextFieldMapping()
	keywordListFieldMapping.Analyzer = "keywordList"
	keywordListFieldMapping.IncludeInAll = false
	keywordListFieldMapping.IncludeTermVectors = false
	keywordListFieldMapping.Store = false

	// Numbers.
	numFieldMapping = bleve.NewNumericFieldMapping()
	numFieldMapping.Store = false

	// Dates.
	dateFieldMapping = bleve.NewDateTimeFieldMapping()
	dateFieldMapping.Store = false

	// Booleans.
	boolFieldMapping = bleve.NewBooleanFieldMapping()
	boolFieldMapping.Store = false
}

type Field struct {
	Name string
	Type FType
}

type FType int

const (
	UNKNOWN_FTYPE FType = iota
	TEXT_FTYPE
	KEYWORD_FTYPE
	KEYWORD_LIST_FTYPE
	NUMBER_FTYPE
	DATE_FTYPE
	BOOLEAN_FTYPE
)

func (ft FType) toBleveMapping() *mapping.FieldMapping {
	once.Do(initializeMappings)

	switch ft {
	case TEXT_FTYPE:
		return englishTextFieldMapping
	case KEYWORD_FTYPE:
		return keywordFieldMapping
	case NUMBER_FTYPE:
		return numFieldMapping
	case DATE_FTYPE:
		return dateFieldMapping
	case BOOLEAN_FTYPE:
		return boolFieldMapping
	case KEYWORD_LIST_FTYPE:
		return keywordListFieldMapping
	default:
		panic(errors.New("unrecognized field type"))
	}
}
