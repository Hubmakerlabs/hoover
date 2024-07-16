package bundle

import (
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/bundlr"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/config"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/listener"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/model"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/monitoring"
	monitor_bundler "github.com/Hubmakerlabs/hoover/pkg/arweave/utils/monitoring/bundler"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/task"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/turbo"
)

type Controller struct {
	*task.Task
}

// +---------------+
// |   Collector   |
// |               |
// |               |
// | +-----------+ |
// | |  Poller   | |             +----------+         +-----------+                +-----------------+
// | +-----------+ |     tx      |          |   pd    |           |  network_info  |                 |
// |               +------------>| Bundler  +-------->| Confirmer |<-------------- | Network Monitor |
// | +-----------+ |             |          |         |           |                |                 |
// | |  Notifier | |             +----------+         +-----------+                +-----------------+
// | +-----------+ |
// |               |
// +---------------+
// Main class that orchestrates main syncer functionalities
func NewController(config *config.Config) (self *Controller, err error) {
	self = new(Controller)
	self.Task = task.NewTask(config, "bundle-controller")

	// SQL database
	db, err := model.NewConnection(self.Ctx, config, "bundler")
	if err != nil {
		return
	}

	// Arweave client
	arweaveClient := arweave.NewClient(self.Ctx, config)

	// Bundlr client
	irysClient := bundlr.NewClient(self.Ctx, &config.Bundlr)
	turboClient := turbo.NewClient(self.Ctx, &config.Bundlr)

	// Monitoring
	monitor := monitor_bundler.NewMonitor()
	server := monitoring.NewServer(config).
		WithMonitor(monitor)

	// Gets interactions to bundle from the database
	collector := NewCollector(config, db).
		WithMonitor(monitor)

	// Monitors latest Arweave network block height
	networkMonitor := listener.NewNetworkMonitor(config).
		WithClient(arweaveClient).
		WithMonitor(monitor).
		WithInterval(config.NetworkMonitor.Period).
		WithRequiredConfirmationBlocks(0).
		WithEnableOutput(false /*disable output channel to avoid blocking*/)

	// Sends interactions to bundlr.network
	bundler := NewBundler(config, db).
		WithInputChannel(collector.Output).
		WithMonitor(monitor).
		WithIrysClient(irysClient).
		WithTurboClient(turboClient)

	// Confirmer periodically updates the state of the bundled interactions
	confirmer := NewConfirmer(config).
		WithDB(db).
		WithMonitor(monitor).
		WithNetworkMonitor(networkMonitor).
		WithInputChannel(bundler.Output)

	// Periodically run queries. Results stored in the monitor.
	dbPoller := monitoring.NewDbPoller(config).
		WithDB(db).
		WithQuery(config.Bundler.DBPollerInterval,
			&monitor.GetReport().Bundler.State.PendingBundleItems,
			"SELECT count(1) FROM bundle_items WHERE state='PENDING'")

	// Setup everything, will start upon calling Controller.Start()
	self.Task.
		WithConditionalSubtask(!config.Bundler.NotifierDisabled && config.Bundler.PollerDisabled,
			dbPoller.Task).
		WithSubtask(confirmer.Task).
		WithSubtask(bundler.Task).
		WithSubtask(monitor.Task).
		WithSubtask(networkMonitor.Task).
		WithSubtask(server.Task).
		WithSubtask(collector.Task)
	return
}
