package utils

import (
	"errors"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func IsTrue(s string) bool {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "true":
		return true
	case "false":
		return false
	default:
		panic(errors.New(s + " is not a boolean"))
	}
}

type MemoryMonitor struct {
	softLimit int64
	hardLimit int64
}

func NewMemoryMonitor(softLimitMB, hardLimitMB int64) *MemoryMonitor {
	return &MemoryMonitor{
		softLimit: softLimitMB * 1024 * 1024,
		hardLimit: hardLimitMB * 1024 * 1024,
	}
}

func (m *MemoryMonitor) Check() (bool, int64) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	current := int64(memStats.HeapAlloc) // Usamos HeapAlloc para maior precisão
	
	if current > m.hardLimit {
		return false, current
	}
	
	if current > m.softLimit {
		runtime.GC()
		time.Sleep(100 * time.Millisecond) // Espera o GC completar
		runtime.ReadMemStats(&memStats)
		current = int64(memStats.HeapAlloc)
	}
	
	return current <= m.hardLimit, current
}

func CopyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
