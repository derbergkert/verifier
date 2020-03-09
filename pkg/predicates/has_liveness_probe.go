package predicates

import v1 "k8s.io/api/core/v1"

var PortValidator = Conjunction(
	NumberValue(`>=0 `),
	NumberValue(`<=9999 `),
)

// HasLivenessProbe is a predicate that determines if a pod has a liveness probe configured correctly.
var HasLivenessProbe = func() Predicate {
	pred, err := NewPredicateFactory(v1.Pod{}).Build(
		Field("spec.containers.livenessprobe",
			Conjunction(
				Disjunction(
					Field("exec",
						Conjunction(
							Field("command", StringValue(`\w`)),
						),
					),
					Field("httpGet",
						Conjunction(
							Field("port", PortValidator),
						),
					),
					Field("tcpSocket",
						Conjunction(
							Field("port", PortValidator),
						),
					),
				),
				Field("initialDelaySeconds", NumberValue("")),
				Field("periodSeconds", NumberValue("")),
			),
		),
	)
	if err != nil {
		panic(err)
	}
	return pred
}()

