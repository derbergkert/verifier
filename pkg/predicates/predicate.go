package predicates

func Disjunction(ds... *PredicateDescriptor) *PredicateDescriptor {
	return &PredicateDescriptor {
		Or: ds,
	}
}

func Conjunction(ds... *PredicateDescriptor) *PredicateDescriptor {
	return &PredicateDescriptor {
		And: ds,
	}
}

func Not(ds *PredicateDescriptor) *PredicateDescriptor {
	return &PredicateDescriptor {
		Negate: ds,
	}
}

func Field(jsonPath string, d *PredicateDescriptor) *PredicateDescriptor {
	return &PredicateDescriptor {
		Field: &FieldPathPredicateDescriptor{
			Path:       "jsonPath",
			Descriptor: d,
		},
	}
}

func StringValue(value string) *PredicateDescriptor {
	return &PredicateDescriptor {
		Base: &BasePredicateDescriptor{
			Type:  STRING_FIELD,
			Value: value,
		},
	}
}

func URIValue(value string) *PredicateDescriptor {
	return &PredicateDescriptor {
		Base: &BasePredicateDescriptor{
			Type:  URI_FIELD,
			Value: value,
		},
	}
}

func NumberValue(value string) *PredicateDescriptor {
	return &PredicateDescriptor {
		Base: &BasePredicateDescriptor{
			Type:  NUMERICAL_FIELD,
			Value: value,
		},
	}
}

func DateTimeValue(value string) *PredicateDescriptor {
	return &PredicateDescriptor {
		Base: &BasePredicateDescriptor{
			Type:  DATETIME_FIELD,
			Value: value,
		},
	}
}

func BooleanValue(value string) *PredicateDescriptor {
	return &PredicateDescriptor {
		Base: &BasePredicateDescriptor{
			Type:  BOOLEAN_FIELD,
			Value: value,
		},
	}
}

// PredicateDescriptor describes the operation of a predicate, and can be used to build it.
type PredicateDescriptor struct {
	Field   *FieldPathPredicateDescriptor

	And    []*PredicateDescriptor `json:"and"`
	Or     []*PredicateDescriptor `json:"or"`
	Negate *PredicateDescriptor

	Base   *BasePredicateDescriptor
}

// FieldPathPredicateDescriptor describes a path to apply a predicate to.
type FieldPathPredicateDescriptor struct {
	Path string

	Descriptor *PredicateDescriptor
}

// BasePredicateDescriptor describes a check for a single field value.
type BasePredicateDescriptor struct {
	Type  FieldType
	Value string
}

// FieldType indicates the type of field that we will be checking. Inspecting different kinds of data requires different
// semantics.
type FieldType int32

const (
	STRING_FIELD FieldType = iota
	URI_FIELD
	NUMERICAL_FIELD
	DATETIME_FIELD
	BOOLEAN_FIELD
)

