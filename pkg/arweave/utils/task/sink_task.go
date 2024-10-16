package task

import (
	"math"
	"time"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/config"
	"github.com/cenkalti/backoff/v4"

	"github.com/gammazero/deque"
)

// Store handles saving data to the database in na robust way.
// - groups incoming Interactions into batches,
// - ensures data isn't stuck even if a batch isn't big enough
type SinkTask[In any] struct {
	*Task

	// Channel for the data to be processed
	input chan In

	// Periodically called to handle a batch of processed data
	onFlush func([]In) error

	// Queue for the processed data
	queue deque.Deque[In]

	// Batch size that will trigger the onFlush function
	batchSize int

	// Flush interval
	flushInterval time.Duration

	// Max time flush should be retried. 0 means no limit.
	maxElapsedTime time.Duration

	// Max times between flush retries
	maxInterval time.Duration
}

func NewSinkTask[In any](config *config.Config,
	name string) (self *SinkTask[In]) {
	self = new(SinkTask[In])

	// Defaults
	self.maxElapsedTime = 0
	self.maxInterval = 15 * time.Second

	self.Task = NewTask(config, name).
		WithSubtaskFunc(self.run).
		WithWorkerPool(1, 1)

	return
}

func (self *SinkTask[In]) WithBatchSize(batchSize int) *SinkTask[In] {
	self.batchSize = batchSize
	exp := uint(math.Round(math.Logb(float64(batchSize)))) + 1
	self.queue.SetMinCapacity(exp)
	return self
}

func (self *SinkTask[In]) WithInputChannel(v chan In) *SinkTask[In] {
	self.input = v
	return self
}

func (self *SinkTask[In]) WithOnFlush(interval time.Duration,
	f func([]In) error) *SinkTask[In] {
	self.flushInterval = interval
	self.onFlush = f
	return self
}

func (self *SinkTask[In]) WithBackoff(maxElapsedTime, maxInterval time.Duration) *SinkTask[In] {
	self.maxElapsedTime = maxElapsedTime
	self.maxInterval = maxInterval
	return self
}

func (self *SinkTask[In]) flush() {
	size := self.queue.Len()
	data := make([]In, 0, size)
	for i := 0; i < size; i++ {
		data = append(data, self.queue.PopFront())
	}

	self.SubmitToWorker(func() {
		// Expotentially increase the interval between retries
		// Never stop retrying
		// Wait at most maxBackoffInterval between retries
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = self.maxElapsedTime
		b.MaxInterval = self.maxInterval

		err := backoff.Retry(func() error {
			return self.onFlush(data)
		}, b)

		if err != nil {
			self.Log.WithError(err).Error("Failed to flush data to sink")
		}
	})
}

// Receives data from the input channel and saves in the database
func (self *SinkTask[In]) run() (err error) {
	// Used to ensure data isn't stuck in Processor for too long
	timer := time.NewTimer(self.flushInterval)

	for {
		select {
		case in, ok := <-self.input:
			if !ok {
				// The only way input channel is closed is that the Processor's source is stopping
				// There will be no more data, flush everything there is and quit.
				self.flush()

				return
			}

			self.queue.PushBack(in)

			if self.queue.Len() >= self.batchSize {
				self.flush()
			}

		case <-timer.C:
			// Flush is called even if the queue is empty
			self.flush()
			timer = time.NewTimer(self.flushInterval)
		}
	}
}
