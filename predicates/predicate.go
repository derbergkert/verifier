package predicates

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

// PredicateType
type PredicateType int32

const (
	FIELD PredicateType = iota
	CONJUNCT
	DISJUNCT
	NEGATE
)

// FieldType
type FieldType int32

const (
	STRING FieldType = iota
	NUMERICAL
	DATETIME
)

// Predicate
type Predicate struct {
	Type PredicateType

	Field   *FieldPredicate
	And    []*Predicate `json:"and"`
	Or     []*Predicate `json:"or"`
	Negate *Predicate
}

// FieldPredicate
type FieldPredicate struct {
	Type FieldType

	Field      string
	Comparator string
}

// ParsePredicates converts a Predicate object into an executable function which takes in a JSON object and returns
// an error message if the input object matched the predicate, which describes the matching fields.
func ParsePredicate(sPredicate *Predicate) (func(map[string]interface{}) error, error) {
	if sPredicate == nil {
		return bleve.NewMatchAllQuery(), nil
	}

	if sPredicate.Base != "" {
		if len(sPredicate.And) != 0 || len(sPredicate.Or) != 0 {
			return nil, errors.New("Predicate should be of a single type")
		}
		return ParseBasePredicate(sPredicate.Base, index)
	} else if len(sPredicate.And) != 0 {
		if sPredicate.Or != nil {
			return nil, errors.New("Predicate should be of a single type")
		}
		return ParseAndPredicate(sPredicate.And, index)
	} else if len(sPredicate.Or) != 0 {
		return ParseOrPredicate(sPredicate.Or, index)
	}
	return nil, errors.New("empty Predicate")
}

// ParseAndPredicate parses the ANDs of a Predicate.
func ParseAndPredicate(andPredicates []*Predicate, index *Index) (query.Query, error) {
	var ands []query.Query
	for _, Predicate := range andPredicates {
		q, err := ParsePredicate(Predicate, index)
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

// ParseOrPredicate parses the ORs of a Predicate.
func ParseOrPredicate(orPredicates []*Predicate, index *Index) (query.Query, error) {
	var ors []query.Query
	for _, Predicate := range orPredicates {
		q, err := ParsePredicate(Predicate, index)
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

// ParseBasePredicate parses a string formatted, single field, query.
func ParseBasePredicate(basePredicate string, index *Index) (query.Query, error) {
	fixes := strings.SplitN(basePredicate, ":", 2)
	if len(fixes) != 2 {
		return nil, fmt.Errorf("Predicate must have a field and value with ':' in between: %s", basePredicate)
	}
	for _, doc := range index.Documents {
		if fType, hasField := doc.Fields[fixes[0]]; hasField {
			return parsePredicateOfType(fixes[0], fixes[1], fType)
		}
	}
	return nil, fmt.Errorf("cannot find indexed field: %s", fixes[0])
}

// Helper functions for parsing a string query.
///////////////////////////////////////////////

func parsePredicateOfType(prefix, suffix string, fType FType) (query.Query, error) {
	switch fType {
	case TEXT_FTYPE:
		return parseTextPredicate(prefix, suffix)
	case KEYWORD_FTYPE:
		return parseKeywordPredicate(prefix, suffix)
	case KEYWORD_LIST_FTYPE:
		return parseKeywordPredicate(prefix, suffix)
	case NUMBER_FTYPE:
		return parseNumberPredicate(prefix, suffix)
	case DATE_FTYPE:
		return parseDatePredicate(prefix, suffix)
	default:
		return nil, fmt.Errorf("cannot handle Predicate of type %d", fType)
	}
}

func parseTextPredicate(prefix, suffix string) (query.Query, error) {
	mq := bleve.NewMatchQuery(suffix)
	mq.SetField(prefix)
	return mq, nil
}

func parseKeywordPredicate(prefix, suffix string) (query.Query, error) {
	mq := bleve.NewMatchQuery(suffix)
	mq.SetField(prefix)
	return mq, nil
}

func parseNumberPredicate(prefix, suffix string) (query.Query, error) {
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

func parseDatePredicate(prefix, suffix string) (query.Query, error) {
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
