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
	internal, err :=  parsePredicate("", pf.example, d)
	if err != nil {
		return nil, err
	}
	return func(input interface{}) bool {
		return internal(reflect.ValueOf(input))
	}, nil
}


func parsePredicate(currentPath string, currentExample interface{}, predD *PredicateDescriptor) (internalPredicate, error) {
	if predD == nil {
		return nil, errors.New("received a nil descriptor")
	}
	if predD.Field != nil {
		return parseFieldPredicate(currentPath, currentExample, predD)
	} else if len(predD.And) != 0 {
		return parseAndPredicate(currentPath, currentExample, predD)
	} else if len(predD.Or) != 0 {
		return parseOrPredicate(currentPath, currentExample, predD)
	}
	return nil, errors.New(fmt.Sprintf("empty descriptor at path %s", currentPath))
}

func parseAndPredicate(currentPath string, currentExample interface{}, andPredicate *PredicateDescriptor) (internalPredicate, error) {
	var ands []internalPredicate
	for _, pred := range andPredicate.And {
		q, err := parsePredicate(currentPath, currentExample, pred)
		if err != nil {
			return nil, err
		}
		ands = append(ands, q)
	}
	if len(ands) == 0 {
		return nil, errors.New(fmt.Sprintf("empty conjunction predicate encountered at path %s", currentPath))
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

func parseOrPredicate(currentPath string, currentExample interface{}, orPredicate *PredicateDescriptor) (internalPredicate, error) {
	var ors []internalPredicate
	for _, pred := range orPredicate.Or{
		q, err := parsePredicate(currentPath, currentExample, pred)
		if err != nil {
			return nil, err
		}
		ors = append(ors, q)
	}
	if len(ors) == 0 {
		return nil, errors.New(fmt.Sprintf("empty disjunction predicate encountered at path %s", currentPath))
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

func parseFieldPredicate(currentPath string, currentExample interface{}, pred *PredicateDescriptor) (internalPredicate, error) {
	newPath, newExample, extractor, err := fieldExtractor(currentPath, currentExample, pred.Field.Path)
	if err != nil {
		return nil, err
	}
	child, err := parsePredicate(newPath, newExample, pred.Field.Descriptor)
	if err != nil {
		return nil, err
	}
	return func(input reflect.Value) bool {
		return child(extractor(input))
	}, nil
}

func fieldExtractor(currentPath string, currentExample interface{}, jsonPath string) (string, interface{}, func(reflect.Value) reflect.Value, error) {
	steps := strings.Split(jsonPath, ".")
	if len(steps) == 0 {
		return "", nil, nil, errors.New(fmt.Sprintf("empty json path for field after: %s", currentPath))
	}
	for _, step := range steps {
		currType := reflect.TypeOf(currentExample)
		for fieldIndex := 0; fieldIndex < currType.NumField(); fieldIndex++ {
			currField := currType.FieldByIndex(fieldIndex)
			jsonTags := currField.
		}
		currentPath = fmt.Sprintf("%s.%s", currentPath, step)
	}
	return currentPath, nil, nil, nil
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
