package metrix

import (
	"fmt"
	"errors"
	"bytes"

	"github.com/rcrowley/go-metrics"
)

var registry metrics.Registry

func init() {
	registry = metrics.NewRegistry()
}

func GetRegistry() *metrics.Registry {
	return &registry
}

func GetCounter( name string ) ( metrics.Counter, error) {
	if registry != nil {
		return metrics.GetOrRegisterCounter( name, registry ), nil
	}

	s := fmt.Sprintf("get counter failed (no registry)" )
	return nil, errors.New(s)
}

func GetGauge( name string ) ( metrics.Gauge, error) {
	if registry != nil {
		return metrics.GetOrRegisterGauge( name, registry ), nil
	}

	s := fmt.Sprintf("get gauge failed (no registry)" )
	return nil, errors.New(s)
}

func GetHistogram( name string, samples int ) (metrics.Histogram, error ) {
	if registry != nil {
		s := metrics.NewUniformSample( samples )
		return metrics.GetOrRegisterHistogram( name, registry, s ), nil
	}

	s := fmt.Sprintf("get histogram failed (no registry)" )
	return nil, errors.New(s)
}

func GetJson() ( []byte, error ) {
	if registry != nil {
		b := &bytes.Buffer{}
		metrics.WriteJSONOnce( registry, b )

		return b.Bytes(), nil
	}

	s := fmt.Sprintf("get json string failed (no registry)" )
	return nil, errors.New(s)
}
