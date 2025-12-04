// Package admission provides addmission plugins for GIE.
package admission

import (
	"context"
	"encoding/json"

	"github.com/llm-d/llm-d-inference-scheduler/pkg/metrics"
	"github.com/llm-d/llm-d-kv-cache-manager/pkg/utils/logging"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/plugins"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/requestcontrol"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/scheduling/types"
)

const (
	// ScaleFromZeroType is the AdmissionPlugin type that is used in plugins registry.
	ScaleFromZeroType = "scale-from-zero-handler"
)

var _ requestcontrol.AdmissionPlugin = &ScaleFromZeroTypeHandler{}

// ScaleFromZeroTypeHandle Admission plugin
type ScaleFromZeroTypeHandler struct {
	typedName plugins.TypedName
}

// TypedName returns the type and name tuple of this plugin instance.
func (f *ScaleFromZeroTypeHandler) TypedName() plugins.TypedName {
	return f.typedName
}

// WithName sets the name of the plugin.
func (f *ScaleFromZeroTypeHandler) WithName(name string) *ScaleFromZeroTypeHandler {
	f.typedName.Name = name
	return f
}

// DestinationEndpointServedVerifierFactory defines the factory function for DestinationEndpointServedVerifier.
func ScaleFromZeroTypeHandlerFactory(name string, _ json.RawMessage, _ plugins.Handle) (plugins.Plugin, error) {
	return NewScaleFromZeroTypeHandler().WithName(name), nil
}

func NewScaleFromZeroTypeHandler() *ScaleFromZeroTypeHandler {
	return &ScaleFromZeroTypeHandler{
		typedName: plugins.TypedName{Type: ScaleFromZeroType},
	}
}

// ScaleFromZeroTypeHandleris the handler for the ScaleFromZeroTypeHandler extension point.
func (p *ScaleFromZeroTypeHandler) AdmitRequest(ctx context.Context, request *types.LLMRequest, pods []types.Pod) error {
	logger := log.FromContext(ctx).WithName(p.TypedName().String())
	logger.V(logging.DEBUG).Info("Verifying if inferencePool is not empty")

	metrics.RecordScaleZeroWaitingRequest("samplePool", request.TargetModel, float64(1.0))

	return nil
}
