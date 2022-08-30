package smsproxy

import (
	"gitlab.com/devskiller-tasks/messaging-app-golang/fastsmsing"
	"math"
	"sync"
)

type clientResult struct {
	messagesBatch  []fastsmsing.Message
	err            error
	currentAttempt int
	maxAttempts    int
}

type Stats struct {
	success           int
	failed            int
	successPercentage float64
}

type ClientStatistics interface {
	Send(result clientResult)
	GetStatistics() Stats
}

func newClientStatistics() ClientStatistics {
	return &inMemoryClientStatistics{clientResults: make([]clientResult, 0), lock: sync.RWMutex{}}
}

type inMemoryClientStatistics struct {
	lock          sync.RWMutex
	clientResults []clientResult
}

func (s *inMemoryClientStatistics) Send(result clientResult) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.clientResults = append(s.clientResults, result)
}

func (s *inMemoryClientStatistics) GetStatistics() Stats {
	s.lock.Lock()
	defer s.lock.Unlock()
	success, failed, percentage := calculateStats(s.clientResults)
	return Stats{
		success:           success,
		failed:            failed,
		successPercentage: percentage,
	}
}

func calculateStats(results []clientResult) (int, int, float64) {
	success := 0
	failed := 0
	for _, result := range results {
		if result.err != nil {
			failed += 1
		} else {
			success += 1
		}
	}
	return success, failed, toPercentage(success, success+failed)
}

func toPercentage(dividend, divisor int) float64 {
	percentage := float64(dividend) / float64(divisor) * 100.0
	roundedPercentage := math.Floor(percentage*100) / 100
	return roundedPercentage
}
