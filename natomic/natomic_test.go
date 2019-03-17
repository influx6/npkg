package natomic

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestAtomicCounter(t *testing.T) {
	var ws IntCounter
	var waiter sync.WaitGroup
	for i := 0; i < 100; i++ {
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			for i := int64(0); i < 1000; i++ {
				ws.Add(i)
			}
		}()

		waiter.Add(1)
		go func() {
			defer waiter.Done()
			for i := 0; i < 1000; i++ {
				ws.Read()
			}
		}()
	}

	for i := 0; i < 1000; i++ {
		ws.Read()
	}

	waiter.Wait()
}

func BenchmarkIntSwitchWrite(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var ws IntSwitch
	for i := 0; i < b.N; i++ {
		ws.Flip(int64(i))
	}
}

func BenchmarkAtomicCounter(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var ws IntCounter
	var waiter sync.WaitGroup
	for i := 0; i < 100; i++ {
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			for i := int64(0); i < 1000; i++ {
				ws.Add(i)
			}
		}()

		waiter.Add(1)
		go func() {
			defer waiter.Done()
			for i := 0; i < 1000; i++ {
				ws.Read()
			}
		}()
	}

	for i := 0; i < 1000; i++ {
		ws.Read()
	}

	waiter.Wait()
}

func BenchmarkAtomicInt(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var ws int64
	var waiter sync.WaitGroup
	for i := 0; i < 100; i++ {
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			for i := 0; i < 1000; i++ {
				atomic.AddInt64(&ws, 1)
			}
		}()
	}

	for i := 0; i < 100; i++ {
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			for i := 0; i < 1000; i++ {
				atomic.LoadInt64(&ws)
			}
		}()
	}

	for i := 0; i < 1000; i++ {
		atomic.LoadInt64(&ws)
	}

	waiter.Wait()
}

func BenchmarkUintSwitchWrite(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var ws UintSwitch
	for i := 0; i < b.N; i++ {
		ws.Flip(uint64(i))
	}
}

func BenchmarkBoolSwitchWrite(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var ws BoolSwitch
	for i := 0; i < b.N; i++ {
		ws.Flip(true)
	}
}

func BenchmarkIntSwitchRead(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var ws IntSwitch
	ws.Flip(int64(1))
	for i := 0; i < b.N; i++ {
		ws.Read()
	}
}

func BenchmarkUintSwitchRead(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var ws UintSwitch
	ws.Flip(uint64(1))
	for i := 0; i < b.N; i++ {
		ws.Read()
	}
}

func BenchmarkBoolSwitchRead(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var ws BoolSwitch
	ws.Flip(true)
	for i := 0; i < b.N; i++ {
		ws.Read()
	}
}
