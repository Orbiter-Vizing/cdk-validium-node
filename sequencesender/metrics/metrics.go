package metrics

import (
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Prefix for the metrics of the sequencer package.
	Prefix = "sequencer_"
	// SequencesSentToL1CountName is the name of the metric that counts the sequences sent to L1.
	SequencesSentToL1CountName = Prefix + "sequences_sent_to_L1_count"
)

// Register the metrics for the sequencer package.
func Register() {
	var (
		counters []prometheus.CounterOpts
	)

	counters = []prometheus.CounterOpts{
		{
			Name: SequencesSentToL1CountName,
			Help: "[SEQUENCER] total count of sequences sent to L1",
		},
	}

	metrics.RegisterCounters(counters...)
}

// SequencesSentToL1 increases the counter by the provided number of sequences
// sent to L1.
func SequencesSentToL1(numSequences float64) {
	metrics.CounterAdd(SequencesSentToL1CountName, numSequences)
}
