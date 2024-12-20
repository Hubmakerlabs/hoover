package monitor_checker

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	monitor *Monitor

	// Run
	UpForSeconds *prometheus.Desc

	BundlesTakenFromDb     *prometheus.Desc
	AllCheckedBundles      *prometheus.Desc
	FinishedBundles        *prometheus.Desc
	UnfinishedBundles      *prometheus.Desc
	IrysUnfinishedBundles  *prometheus.Desc
	TurboUnfinishedBundles *prometheus.Desc
	DbStateUpdated         *prometheus.Desc
	IrysGetStatusError     *prometheus.Desc
	TurboGetStatusError    *prometheus.Desc
	DbStateUpdateError     *prometheus.Desc
}

func NewCollector() *Collector {
	return &Collector{
		UpForSeconds: prometheus.NewDesc("up_for_seconds", "", nil,
			nil),
		BundlesTakenFromDb: prometheus.NewDesc("bundles_taken_from_db", "",
			nil, nil),
		AllCheckedBundles: prometheus.NewDesc("all_checked_bundles", "",
			nil, nil),
		FinishedBundles: prometheus.NewDesc("finished_bundles", "", nil,
			nil),
		UnfinishedBundles: prometheus.NewDesc("unfinished_bundles", "",
			nil, nil),
		IrysUnfinishedBundles: prometheus.NewDesc("irys_unfinished_bundles",
			"", nil, nil),
		TurboUnfinishedBundles: prometheus.NewDesc("turbo_unfinished_bundles",
			"", nil, nil),
		DbStateUpdated: prometheus.NewDesc("db_state_updated", "", nil,
			nil),
		IrysGetStatusError: prometheus.NewDesc("irys_check_state_error", "",
			nil, nil),
		TurboGetStatusError: prometheus.NewDesc("turbo_check_state_error",
			"", nil, nil),
		DbStateUpdateError: prometheus.NewDesc("db_state_update_error", "",
			nil, nil),
	}
}

func (self *Collector) WithMonitor(m *Monitor) *Collector {
	self.monitor = m
	return self
}

func (self *Collector) Describe(ch chan<- *prometheus.Desc) {
	// Run
	ch <- self.UpForSeconds

	// Checker
	ch <- self.BundlesTakenFromDb
	ch <- self.AllCheckedBundles
	ch <- self.FinishedBundles
	ch <- self.UnfinishedBundles
	ch <- self.IrysUnfinishedBundles
	ch <- self.TurboUnfinishedBundles
	ch <- self.DbStateUpdated
	ch <- self.IrysGetStatusError
	ch <- self.TurboGetStatusError
	ch <- self.DbStateUpdateError
}

// Collect implements required collect function for all promehteus collectors
func (self *Collector) Collect(ch chan<- prometheus.Metric) {
	// Run
	ch <- prometheus.MustNewConstMetric(self.UpForSeconds,
		prometheus.GaugeValue,
		float64(self.monitor.Report.Run.State.UpForSeconds.Load()))

	// Checker
	ch <- prometheus.MustNewConstMetric(self.BundlesTakenFromDb,
		prometheus.CounterValue,
		float64(self.monitor.Report.Checker.State.BundlesTakenFromDb.Load()))
	ch <- prometheus.MustNewConstMetric(self.AllCheckedBundles,
		prometheus.CounterValue,
		float64(self.monitor.Report.Checker.State.AllCheckedBundles.Load()))
	ch <- prometheus.MustNewConstMetric(self.FinishedBundles,
		prometheus.CounterValue,
		float64(self.monitor.Report.Checker.State.FinishedBundles.Load()))
	ch <- prometheus.MustNewConstMetric(self.UnfinishedBundles,
		prometheus.CounterValue,
		float64(self.monitor.Report.Checker.State.UnfinishedBundles.Load()))
	ch <- prometheus.MustNewConstMetric(self.IrysUnfinishedBundles,
		prometheus.CounterValue,
		float64(self.monitor.Report.Checker.State.IrysUnfinishedBundles.Load()))
	ch <- prometheus.MustNewConstMetric(self.TurboUnfinishedBundles,
		prometheus.CounterValue,
		float64(self.monitor.Report.Checker.State.TurboUnfinishedBundles.Load()))
	ch <- prometheus.MustNewConstMetric(self.DbStateUpdated,
		prometheus.CounterValue,
		float64(self.monitor.Report.Checker.State.DbStateUpdated.Load()))
	ch <- prometheus.MustNewConstMetric(self.IrysGetStatusError,
		prometheus.CounterValue,
		float64(self.monitor.Report.Checker.Errors.IrysGetStatusError.Load()))
	ch <- prometheus.MustNewConstMetric(self.TurboGetStatusError,
		prometheus.CounterValue,
		float64(self.monitor.Report.Checker.Errors.TurboGetStatusError.Load()))
	ch <- prometheus.MustNewConstMetric(self.DbStateUpdateError,
		prometheus.CounterValue,
		float64(self.monitor.Report.Checker.Errors.DbStateUpdateError.Load()))
}
