// Copyright 2017 Percona LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongod

import (
	"unicode/utf8"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	locksTimeLockedGlobalMicrosecondsTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "locks_time_locked_global_microseconds_total"),
		"amount of time in microseconds that any database has held the global lock",
		[]string{"type", "database"},
		nil,
	)
)
var (
	locksTimeLockedLocalMicrosecondsTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "locks_time_locked_local_microseconds_total"),
		"amount of time in microseconds that any database has held the local lock",
		[]string{"type", "database"},
		nil,
	)
)
var (
	locksTimeAcquiringGlobalMicrosecondsTotalDesc = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "locks_time_acquiring_global_microseconds_total"),
		"amount of time in microseconds that any database has spent waiting for the global lock",
		[]string{"type", "database"},
		nil,
	)
)

// LockStatsMap is a map of lock stats
type LockStatsMap map[string]LockStats

// ReadWriteLockTimes information about the lock
type ReadWriteLockTimes struct {
	Read       float64 `bson:"R"`
	Write      float64 `bson:"W"`
	ReadLower  float64 `bson:"r"`
	WriteLower float64 `bson:"w"`
}

// LockStats lock stats
type LockStats struct {
	TimeLockedMicros    ReadWriteLockTimes `bson:"timeLockedMicros"`
	TimeAcquiringMicros ReadWriteLockTimes `bson:"timeAcquiringMicros"`
}

// Export exports the data to prometheus.
func (locks LockStatsMap) Export(ch chan<- prometheus.Metric) {
	for key, locks := range locks {
		// Remove non UTF8 bytes from keys.
		if !utf8.ValidString(key) {
			v := make([]rune, 0, len(key))
			for i, r := range key {
				if r == utf8.RuneError {
					_, size := utf8.DecodeRuneInString(key[i:])
					if size == 1 {
						continue
					}
				}
				v = append(v, r)
			}
			key = string(v)
		}

		if key == "." {
			key = "dot"
		}

		ch <- prometheus.MustNewConstMetric(locksTimeLockedGlobalMicrosecondsTotalDesc, prometheus.CounterValue, locks.TimeLockedMicros.Read, "read", key)
		ch <- prometheus.MustNewConstMetric(locksTimeLockedGlobalMicrosecondsTotalDesc, prometheus.CounterValue, locks.TimeLockedMicros.Write, "write", key)

		ch <- prometheus.MustNewConstMetric(locksTimeLockedLocalMicrosecondsTotalDesc, prometheus.CounterValue, locks.TimeLockedMicros.ReadLower, "read", key)
		ch <- prometheus.MustNewConstMetric(locksTimeLockedLocalMicrosecondsTotalDesc, prometheus.CounterValue, locks.TimeLockedMicros.WriteLower, "write", key)

		ch <- prometheus.MustNewConstMetric(locksTimeAcquiringGlobalMicrosecondsTotalDesc, prometheus.CounterValue, locks.TimeAcquiringMicros.ReadLower, "read", key)
		ch <- prometheus.MustNewConstMetric(locksTimeAcquiringGlobalMicrosecondsTotalDesc, prometheus.CounterValue, locks.TimeAcquiringMicros.WriteLower, "write", key)
	}
}

// Describe describes the metrics for prometheus
func (locks LockStatsMap) Describe(ch chan<- *prometheus.Desc) {
	ch <- locksTimeLockedGlobalMicrosecondsTotalDesc
	ch <- locksTimeLockedLocalMicrosecondsTotalDesc
	ch <- locksTimeAcquiringGlobalMicrosecondsTotalDesc
}
