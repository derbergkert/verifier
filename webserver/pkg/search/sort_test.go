package search

import (
	"testing"

	"github.com/blevesearch/bleve/search"
	"github.com/stretchr/testify/assert"
)

func TestSortConversion(t *testing.T) {
	index := BuildIndex(
		"predicate",
		WithDocument(
			BuildDocument(
				"predicate",
				WithField("userId", KEYWORD_FTYPE),
				WithField("tags", KEYWORD_FTYPE),
				WithField("createdTimestamp", NUMBER_FTYPE),
				WithField("updatedTimestamp", NUMBER_FTYPE),
			),
		),
	)

	sort := &Sort{
		Fields: []string{
			"createdTimestamp:<",
			"updatedTimestamp:>",
			"userId:<",
		},
	}

	bSort, err := ParseSort(sort, index)
	assert.Nil(t, err)
	assert.Equal(t, search.SortOrder{
		&search.SortField{
			Field: "createdTimestamp",
			Desc:  false,
			Type:  2,
		},
		&search.SortField{
			Field: "updatedTimestamp",
			Desc:  true,
			Type:  2,
		},
		&search.SortField{
			Field: "userId",
			Desc:  false,
			Type:  1,
		},
	}, bSort)
}
