package diagnosis

import "github.com/boserwuge/k8s-cluster-inspector/internal/model"

type Engine struct {
	rules []Rule
}

func NewEngine(rules ...Rule) *Engine {
	return &Engine{
		rules: rules,
	}
}

func (e *Engine) Run(data ClusterData) []model.Diagnosis {
	results := make([]model.Diagnosis, 0)

	for _, rule := range e.rules {
		if !rule.Enabled() {
			continue
		}

		results = append(results, rule.Evaluate(data)...)
	}

	return results
}
