package predicates

import (
	"errors"
	"fmt"
	"github.com/blevesearch/bleve"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Predicate defines a function which takes in an object, and returns an error indicating some issue.
// Super generic.
type Predicate func(interface{}) bool

// Predicate defines a function which takes in an object, and returns an error indicating some issue.
// Super generic.
type internalPredicate func(value reflect.Value) bool

type PredicateFactory interface {
	Build(d *PredicateDescriptor) (Predicate, error)
}

func NewPredicateFactory(example interface{}) PredicateFactory {
	return &predicateFactoryImpl{
		example: example,
	}
}

type predicateFactoryImpl struct {
	example interface{}
}

func (pf *predicateFactoryImpl) Build(d *PredicateDescriptor) (Predicate, error) {
	internal, err :=  parsePredicate(pf.example, d)
	if err != nil {
		return nil, err
	}
	return func(input interface{}) bool {
		return internal(reflect.ValueOf(input))
	}, nil
}

// ParsePredicates converts a Predicate object into an executable function which takes in a JSON object and returns
// an error message if the input object matched the predicate, which describes the matching fields.
func parsePredicate(example interface{}, predD *PredicateDescriptor) (internalPredicate, error) {
	if predD == nil {
		return nil, errors.New("received a nil descriptor")
	}

	if predD.Field != nil {
		return parseFieldPredicate(example, predD)
	} else if len(predD.And) != 0 {
		if predD.Or != nil {
			return nil, errors.New("Predicate should be of a single type")
		}
		return parseAndPredicate(example, predD)
	} else if len(predD.Or) != 0 {
		return parseOrPredicate(example, predD)
	}
	return nil, errors.New("empty Predicate")
}

// ParseAndPredicate parses the ANDs of a Predicate.
func parseAndPredicate(example interface{}, andPredicate *PredicateDescriptor) (internalPredicate, error) {
	var ands []internalPredicate
	for _, pred := range andPredicate.And {
		q, err := parsePredicate(example, pred)
		if err != nil {
			return nil, err
		}
		ands = append(ands, q)
	}
	if len(ands) == 0 {
		return nil, errors.New("empty CONJUNCTION predicate encountered")
	}
	if len(ands) == 1 {
		return ands[0], nil
	}
	return func(input reflect.Value) bool {
		for _, a := range ands {
			if pass := a(input); !pass {
				return false
			}
		}
		return true
	}, nil
}

func parseOrPredicate(example interface{}, orPredicate *PredicateDescriptor) (internalPredicate, error) {
	var ors []internalPredicate
	for _, pred := range orPredicate.Or{
		q, err := parsePredicate(example, pred)
		if err != nil {
			return nil, err
		}
		ors = append(ors, q)
	}
	if len(ors) == 0 {
		return nil, errors.New("empty DISJUNCTION predicate encountered")
	}
	if len(ors) == 1 {
		return ors[0], nil
	}
	return func(input reflect.Value) bool {
		for _, a := range ors {
			if pass := a(input); pass {
				return true
			}
		}
		return false
	}, nil
}

func parseFieldPredicate(example interface{}, pred *PredicateDescriptor) (internalPredicate, error) {
	child, err := parsePredicate(example, pred.Field.Descriptor)
	if err != nil {
		return nil, err
	}
	return func(input reflect.Value) bool {
		steps := strings.Split(pred.Field.Path, ".")
		for _, step := range steps {
			input.F
		}
		return child(v)
	}, nil
}

func parseBasePredicate(example interface{}, prefix, suffix string, fType FieldType) (internalPredicate, error) {
	switch fType {
	case STRING_FIELD:
		return parseTextPredicate(prefix, suffix)
	case URI_FIELD:
		return parseKeywordPredicate(prefix, suffix)
	case NUMERICAL_FIELD:
		return parseKeywordPredicate(prefix, suffix)
	case DATETIME_FIELD:
		return parseNumberPredicate(prefix, suffix)
	case BOOLEAN_FIELD:
		return parseDatePredicate(prefix, suffix)
	default:
		return nil, fmt.Errorf("cannot handle field of type %d", fType)
	}
}

func parseTextPredicate(prefix, suffix string) (internalPredicate, error) {
	mq := bleve.NewMatchQuery(suffix)
	mq.SetField(prefix)
	return mq, nil
}

func parseKeywordPredicate(prefix, suffix string) (internalPredicate, error) {
	mq := bleve.NewMatchQuery(suffix)
	mq.SetField(prefix)
	return mq, nil
}

func parseNumberPredicate(prefix, suffix string) (internalPredicate, error) {
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

func parseDatePredicate(prefix, suffix string) (internalPredicate, error) {
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
