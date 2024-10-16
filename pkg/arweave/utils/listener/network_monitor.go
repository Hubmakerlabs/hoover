package listener

import (
	"context"
	"sync"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/arweave"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/config"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/monitoring"
	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/task"
)

// Task that periodically checks for new arweave network info.
// Optionally waits for a number of required confirmation blocks before emitting the info
type NetworkMonitor struct {
	*task.Task

	// Runtime configuration
	requiredConfirmationBlocks int64
	client                     *arweave.Client
	monitor                    monitoring.Monitor
	cond                       *sync.Cond

	// Output channel
	Output chan *arweave.NetworkInfo

	// Runtime variables
	lastHeight      int64
	lastNetworkInfo *arweave.NetworkInfo
	isOutputEnabled bool
}

// Using Arweave client periodically checks for blocks of transactions
func NewNetworkMonitor(config *config.Config) (self *NetworkMonitor) {
	self = new(NetworkMonitor)

	self.cond = sync.NewCond(&sync.Mutex{})
	self.Output = make(chan *arweave.NetworkInfo)
	self.isOutputEnabled = true

	self.Task = task.NewTask(config, "network-monitor").
		WithOnAfterStop(func() {
			close(self.Output)
		})
	return
}

func (self *NetworkMonitor) WithMonitor(monitor monitoring.Monitor) *NetworkMonitor {
	self.monitor = monitor
	return self
}

func (self *NetworkMonitor) WithClient(client *arweave.Client) *NetworkMonitor {
	self.client = client
	return self
}

func (self *NetworkMonitor) WithRequiredConfirmationBlocks(requiredConfirmationBlocks int64) *NetworkMonitor {
	self.requiredConfirmationBlocks = requiredConfirmationBlocks
	return self
}

func (self *NetworkMonitor) WithInterval(interval time.Duration) *NetworkMonitor {
	self.Task = self.Task.WithPeriodicSubtaskFunc(interval,
		self.runPeriodically)
	return self
}

func (self *NetworkMonitor) WithEnableOutput(v bool) *NetworkMonitor {
	self.isOutputEnabled = v
	return self
}

func (self *NetworkMonitor) GetLastNetworkInfo() *arweave.NetworkInfo {
	self.cond.L.Lock()
	for self.lastNetworkInfo == nil {
		self.cond.Wait()
	}
	self.cond.L.Unlock()
	return self.lastNetworkInfo
}

// Periodically checks Arweave network info for updated height
func (self *NetworkMonitor) runPeriodically() error {
	// Use a specific URL as the source of truth, to avoid race conditions with SDK
	ctx := context.WithValue(self.Ctx, arweave.ContextForcePeer,
		self.Config.NetworkMonitor.Url)
	ctx = context.WithValue(ctx, arweave.ContextDisablePeers, true)

	networkInfo, err := self.client.GetNetworkInfo(ctx)
	if err != nil {
		self.Log.WithError(err).Error("Failed to get Arweave network info")
		self.monitor.GetReport().NetworkInfo.Errors.NetworkInfoDownloadErrors.Inc()
		return nil
	}

	self.monitor.GetReport().NetworkInfo.State.ArweaveCurrentHeight.Store(uint64(networkInfo.Height))
	self.monitor.GetReport().NetworkInfo.State.ArweaveLastNetworkInfoTimestamp.Store(uint64(time.Now().Unix()))

	// This is the last block height we consider stable
	stableHeight := networkInfo.Height - self.requiredConfirmationBlocks
	// self.Log.WithField("stable_height", stableHeight).
	// 	WithField("last_height", self.lastHeight).
	// 	WithField("height", networkInfo.Height).
	// 	WithField("delta", self.requiredConfirmationBlocks).
	// 	Debug("Network info")

	if stableHeight <= self.lastHeight {
		// Nothing changed, retry later
		return nil
	}

	self.lastHeight = stableHeight

	self.cond.L.Lock()
	self.lastNetworkInfo = &arweave.NetworkInfo{
		Height: stableHeight,
	}
	self.cond.Broadcast()
	self.cond.L.Unlock()

	if self.isOutputEnabled {
		select {
		case <-self.StopChannel:
		case self.Output <- &arweave.NetworkInfo{
			Height: stableHeight,
		}:
		}

		// self.Log.WithField("stable_height", stableHeight).
		// 	Debug("Emited stable network info")
	}

	return nil
}
