package util

import (
	"math"
	"runtime"
	"sync"
)

func ComputeCosine(query []float64, dataset [][]float64) []float64 {
	result := make([]float64, len(dataset))
	resultLock := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	thread := runtime.NumCPU()
	chunkSize := (len(dataset) + thread - 1) / thread
	for i := 0; i < thread; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(dataset) {
			end = len(dataset)
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for iY := start; iY < end; iY++ {
				resultLock.Lock()
				result[iY] = computeCosine(query, dataset[iY])
				resultLock.Unlock()
			}
		}(start, end)
	}
	wg.Wait()
	return result
}

func computeCosine(x, y []float64) float64 {
	var sum, s1, s2 float64
	for i := 0; i < len(x); i++ {
		sum += x[i] * y[i]
		s1 += math.Pow(x[i], 2)
		s2 += math.Pow(y[i], 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0
	}
	return sum / (math.Sqrt(s1) * math.Sqrt(s2))
}
