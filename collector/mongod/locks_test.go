package mongod

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestExportHandlesInvalidUTF8(t *testing.T) {
	// setup
	invalidUTF8LockStatMap := LockStatsMap{
		"\x80": {
			TimeLockedMicros: ReadWriteLockTimes{
				Read:       100,
				Write:      100,
				ReadLower:  100,
				WriteLower: 100,
			},
			TimeAcquiringMicros: ReadWriteLockTimes{
				Read:       100,
				Write:      100,
				ReadLower:  100,
				WriteLower: 100,
			},
		},
	}
	c := make(chan prometheus.Metric, 6)

	expectedName := "database"
	expectedValue := ""
	expected := &io_prometheus_client.LabelPair{
		Name: &expectedName,
		Value: &expectedValue,
	}

	// run
	go func() {
		invalidUTF8LockStatMap.Export(c)
		close(c)
	}()

	// test
	for range c {
		m := &io_prometheus_client.Metric{}
		received := <-c
		received.Write(m)
		assert.Contains(t, m.GetLabel(), expected)
	}
}
