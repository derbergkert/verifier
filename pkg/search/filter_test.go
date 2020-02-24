package search

import (
	"testing"

	"github.com/blevesearch/bleve"
	"github.com/stretchr/testify/assert"
)

func TestFilterConversion(t *testing.T) {
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

	filter := &Filter{
		And: []*Filter{
			&Filter{
				Base: "userId:Rob",
			},
			&Filter{
				Or: []*Filter{
					{
						Base: "createdTimestamp:<1",
					},
					{
						Base: "createdTimestamp:>=5",
					},
				},
			},
		},
	}

	one := float64(1)
	five := float64(5)
	isInc := true
	isntInc := false

	bQuery, err := ParseFilter(filter, index)
	assert.Nil(t, err)
	assert.Equal(t, bleve.NewConjunctionQuery(
		bleve.NewMatchQuery("userId:Rob"),
		bleve.NewDisjunctionQuery(
			bleve.NewNumericRangeInclusiveQuery(nil, &one, nil, &isntInc),
			bleve.NewNumericRangeInclusiveQuery(&five, nil, &isInc, nil),
		),
	), bQuery)
}
