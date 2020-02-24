package search

import (
	"fmt"
	"strings"

	"github.com/blevesearch/bleve/search"
)

type Sort struct {
	Fields []string `json:"fields"`
}

// ParseSort takes in a sort string and returns the sort order for it.
func ParseSort(sSort *Sort, index *Index) (search.SortOrder, error) {
	if sSort == nil {
		return nil, nil
	}

	var bleveSortFields []search.SearchSort
	for _, sf := range sSort.Fields {
		bsf, err := ParseSortField(sf, index)
		if err != nil {
			return nil, err
		}
		bleveSortFields = append(bleveSortFields, bsf)
	}
	return search.SortOrder(bleveSortFields), nil
}

// ParseSortField takes in a sort string and returns the sort order for it.
func ParseSortField(sortField string, index *Index) (*search.SortField, error) {
	fixes := strings.Split(sortField, ":")
	if len(fixes) != 2 {
		return nil, fmt.Errorf("filter must have a field and value with ':' in between: %s", sortField)
	}

	for _, doc := range index.Documents {
		if fType, hasField := doc.Fields[fixes[0]]; hasField {
			t, err := convertFType(fType)
			if err != nil {
				return nil, err
			}
			desc, err := convertOrder(fixes[1])
			if err != nil {
				return nil, err
			}
			return &search.SortField{
				Field: fixes[0],
				Desc:  desc,
				Type:  t,
			}, nil
		}
	}
	return nil, fmt.Errorf("cannot find indexed field: %s", fixes[0])
}

func convertFType(fType FType) (search.SortFieldType, error) {
	switch fType {
	case TEXT_FTYPE:
		return search.SortFieldAsString, nil
	case KEYWORD_FTYPE:
		return search.SortFieldAsString, nil
	case KEYWORD_LIST_FTYPE:
		return search.SortFieldAsString, nil
	case NUMBER_FTYPE:
		return search.SortFieldAsNumber, nil
	case DATE_FTYPE:
		return search.SortFieldAsDate, nil
	default:
		return 0, fmt.Errorf("cannot handle filter of type %d", fType)
	}
}

func convertOrder(suffix string) (bool, error) {
	if suffix == ">" {
		return true, nil
	} else if suffix == "<" {
		return false, nil
	}
	return false, fmt.Errorf("unrecognized sort suffix: %s", suffix)
}
