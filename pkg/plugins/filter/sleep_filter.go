package filter

import (
	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/plugins"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/scheduling/framework"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/scheduling/types"
)

const (
	//Label name
	LabelKey = "vllm.ai/state"
	// Value set for designated sleeping workers
	LabelValue = "sleep"

	// SleepFilterType is the type of the SleepFilter
	SleepFilterType = "sleep-filter"
)

// compile-time type assertion
var _ framework.Filter = &SleepFilter{}

// SleepFilterFactory defines the factory function for the PrefillFilter
func SleepFilterFactory(name string, _ json.RawMessage, _ plugins.Handle) (plugins.Plugin, error) {
	return NewSleepFilter().WithName(name), nil
}

// NewSleepFilter creates and returns an instance of the Filter configured for SleepFilter role
func NewSleepFilter() *SleepFilter {

	selector := metav1.LabelSelector{
		MatchLabels: map[string]string{LabelKey: LabelValue},
	}

	labelSelector, _ := metav1.LabelSelectorAsSelector(&selector)

	return &SleepFilter{
		typedName: plugins.TypedName{Type: SleepFilterType},
		selector:  labelSelector,
	}
}

// ByLabelSelector filters out pods that do not match its label selector criteria
type SleepFilter struct {
	typedName plugins.TypedName
	selector  labels.Selector
}

// TypedName returns the typed name of the plugin
func (f *SleepFilter) TypedName() plugins.TypedName {
	return f.typedName
}

// WithName sets the name of the plugin.
func (f *SleepFilter) WithName(name string) *SleepFilter {
	f.typedName.Name = name
	return f
}

// Filter filters out all pods that do not satisfy the label selector
func (blf *SleepFilter) Filter(_ context.Context, _ *types.CycleState, _ *types.LLMRequest, pods []types.Pod) []types.Pod {
	filtered := []types.Pod{}

	for _, pod := range pods {
		labels := labels.Set(pod.GetPod().Labels)
		if !blf.selector.Matches(labels) {
			filtered = append(filtered, pod)
		}
	}
	return filtered
}
