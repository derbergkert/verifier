package search

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

// Filter represents a serializable query.
type Filter struct {
	Base string    `json:"base"`
	And  []*Filter `json:"and"`
	Or   []*Filter `json:"or"`
}

// ParseFilters returns the AND of all the input filters.
func ParseFilter(sFilter *Filter, index *Index) (query.Query, error) {
	if sFilter == nil {
		return bleve.NewMatchAllQuery(), nil
	}

	if sFilter.Base != "" {
		if len(sFilter.And) != 0 || len(sFilter.Or) != 0 {
			return nil, errors.New("filter should be of a single type")
		}
		return ParseBaseFilter(sFilter.Base, index)
	} else if len(sFilter.And) != 0 {
		if sFilter.Or != nil {
			return nil, errors.New("filter should be of a single type")
		}
		return ParseAndFilter(sFilter.And, index)
	} else if len(sFilter.Or) != 0 {
		return ParseOrFilter(sFilter.Or, index)
	}
	return nil, errors.New("empty filter")
}

// ParseAndFilter parses the ANDs of a filter.
func ParseAndFilter(andFilters []*Filter, index *Index) (query.Query, error) {
	var ands []query.Query
	for _, filter := range andFilters {
		q, err := ParseFilter(filter, index)
		if err != nil {
			return nil, err
		}
		ands = append(ands, q)
	}
	if len(ands) == 1 {
		return ands[0], nil
	}
	return bleve.NewConjunctionQuery(ands...), nil
}

// ParseOrFilter parses the ORs of a filter.
func ParseOrFilter(orFilters []*Filter, index *Index) (query.Query, error) {
	var ors []query.Query
	for _, filter := range orFilters {
		q, err := ParseFilter(filter, index)
		if err != nil {
			return nil, err
		}
		ors = append(ors, q)
	}
	if len(ors) == 1 {
		return ors[0], nil
	}
	return bleve.NewDisjunctionQuery(ors...), nil
}

// ParseBaseFilter parses a string formatted, single field, query.
func ParseBaseFilter(baseFilter string, index *Index) (query.Query, error) {
	fixes := strings.SplitN(baseFilter, ":", 2)
	if len(fixes) != 2 {
		return nil, fmt.Errorf("filter must have a field and value with ':' in between: %s", baseFilter)
	}
	for _, doc := range index.Documents {
		if fType, hasField := doc.Fields[fixes[0]]; hasField {
			return parseFilterOfType(fixes[0], fixes[1], fType)
		}
	}
	return nil, fmt.Errorf("cannot find indexed field: %s", fixes[0])
}

// Helper functions for parsing a string query.
///////////////////////////////////////////////

func parseFilterOfType(prefix, suffix string, fType FType) (query.Query, error) {
	switch fType {
	case TEXT_FTYPE:
		return parseTextFilter(prefix, suffix)
	case KEYWORD_FTYPE:
		return parseKeywordFilter(prefix, suffix)
	case KEYWORD_LIST_FTYPE:
		return parseKeywordFilter(prefix, suffix)
	case NUMBER_FTYPE:
		return parseNumberFilter(prefix, suffix)
	case DATE_FTYPE:
		return parseDateFilter(prefix, suffix)
	default:
		return nil, fmt.Errorf("cannot handle filter of type %d", fType)
	}
}

func parseTextFilter(prefix, suffix string) (query.Query, error) {
	mq := bleve.NewMatchQuery(suffix)
	mq.SetField(prefix)
	return mq, nil
}

func parseKeywordFilter(prefix, suffix string) (query.Query, error) {
	mq := bleve.NewMatchQuery(suffix)
	mq.SetField(prefix)
	return mq, nil
}

func parseNumberFilter(prefix, suffix string) (query.Query, error) {
	if strings.HasPrefix(suffix, "==") {
		num, err := parseNum(strings.TrimPrefix(suffix, "=="))
		if err != nil {
			return nil, err
		}
		inc := true
		return bleve.NewNumericRangeInclusiveQuery(&num, &num, &inc, &inc), nil
	} else if strings.HasPrefix(suffix, "!=") {
		num, err := parseNum(strings.TrimPrefix(suffix, "!="))
		if err != nil {
			return nil, err
		}
		inc := false
		return bleve.NewConjunctionQuery(
			bleve.NewNumericRangeInclusiveQuery(&num, nil, &inc, nil),
			bleve.NewNumericRangeInclusiveQuery(nil, &num, nil, &inc),
		), nil
	} else if strings.HasPrefix(suffix, "<=") {
		num, err := parseNum(strings.TrimPrefix(suffix, "<="))
		if err != nil {
			return nil, err
		}
		inc := true
		return bleve.NewNumericRangeInclusiveQuery(nil, &num, nil, &inc), nil
	} else if strings.HasPrefix(suffix, ">=") {
		num, err := parseNum(strings.TrimPrefix(suffix, ">="))
		if err != nil {
			return nil, err
		}
		inc := true
		return bleve.NewNumericRangeInclusiveQuery(&num, nil, &inc, nil), nil
	} else if strings.HasPrefix(suffix, "<") {
		num, err := parseNum(strings.TrimPrefix(suffix, "<"))
		if err != nil {
			return nil, err
		}
		inc := false
		return bleve.NewNumericRangeInclusiveQuery(nil, &num, nil, &inc), nil
	} else if strings.HasPrefix(suffix, ">") {
		num, err := parseNum(strings.TrimPrefix(suffix, ">"))
		if err != nil {
			return nil, err
		}
		inc := false
		return bleve.NewNumericRangeInclusiveQuery(&num, nil, &inc, nil), nil
	} else {
		return nil, fmt.Errorf("incorrectly formatted number query: %s", suffix)
	}
}

func parseNum(strNum string) (float64, error) {
	num, err := strconv.ParseFloat(strNum, 64)
	if err != nil {
		return 0, err
	}
	return num, err
}

func parseDateFilter(prefix, suffix string) (query.Query, error) {
	if strings.HasPrefix(suffix, "==") {
		num, err := parseDate(strings.TrimPrefix(suffix, "=="))
		if err != nil {
			return nil, err
		}
		inc := true
		return bleve.NewDateRangeInclusiveQuery(num, num, &inc, &inc), nil
	} else if strings.HasPrefix(suffix, "!=") {
		num, err := parseDate(strings.TrimPrefix(suffix, "!="))
		if err != nil {
			return nil, err
		}
		inc := false
		return bleve.NewConjunctionQuery(
			bleve.NewDateRangeInclusiveQuery(num, time.Time{}, &inc, nil),
			bleve.NewDateRangeInclusiveQuery(time.Time{}, num, nil, &inc),
		), nil
	} else if strings.HasPrefix(suffix, "<=") {
		num, err := parseDate(strings.TrimPrefix(suffix, "<="))
		if err != nil {
			return nil, err
		}
		inc := true
		return bleve.NewDateRangeInclusiveQuery(time.Time{}, num, nil, &inc), nil
	} else if strings.HasPrefix(suffix, ">=") {
		num, err := parseDate(strings.TrimPrefix(suffix, ">="))
		if err != nil {
			return nil, err
		}
		inc := true
		return bleve.NewDateRangeInclusiveQuery(num, time.Time{}, &inc, nil), nil
	} else if strings.HasPrefix(suffix, "<") {
		num, err := parseDate(strings.TrimPrefix(suffix, "<"))
		if err != nil {
			return nil, err
		}
		inc := false
		return bleve.NewDateRangeInclusiveQuery(time.Time{}, num, nil, &inc), nil
	} else if strings.HasPrefix(suffix, ">") {
		num, err := parseDate(strings.TrimPrefix(suffix, ">"))
		if err != nil {
			return nil, err
		}
		inc := false
		return bleve.NewDateRangeInclusiveQuery(num, time.Time{}, &inc, nil), nil
	} else {
		return nil, fmt.Errorf("incorrectly formatted number query: %s", suffix)
	}
}

func parseDate(strDate string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000Z", strDate)
}
