package prometheuser

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type DeletableMetric interface {
	Collect(ch chan<- prometheus.Metric)
	Delete(labels prometheus.Labels) bool
}

func DeleteAllGaugeVecValues(gv DeletableMetric) error {
	ch := make(chan prometheus.Metric)
	go func() {
		gv.Collect(ch)
		close(ch)
	}()

	for metric := range ch {
		m := &dto.Metric{}
		err := metric.Write(m)
		if err != nil {
			log.Printf("Error writing metric: %v", err)
			continue
		}

		labels := make(prometheus.Labels)
		for _, label := range m.Label {
			labels[label.GetName()] = label.GetValue()
		}

		// Delete the label set
		deleted := gv.Delete(labels)
		if !deleted {
			log.Printf("Warning: failed to delete label set: %v", labels)
		}
	}

	return nil
}
