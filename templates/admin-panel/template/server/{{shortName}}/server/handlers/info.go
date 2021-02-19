package handlers

import (
	"net/http"
	"os"
	"{{shortName}}/assets"
	"{{shortName}}/core/handler"
	"{{shortName}}/core/sessions"
	"runtime"
	"sort"
	"time"
)

var (
	startTime = time.Now()
)

type info struct {
	Host    string       `json:"Host"`
	Runtime *runtimeInfo `json:"Runtime"`
	Files   []fileInfo   `json:"Files"`
}

// runtimeInfo defines runtime part of service information
type runtimeInfo struct {
	NumCPU       int    `json:"NumCPU"`
	Memory       uint64 `json:"Memory"`
	MemSys       uint64 `json:"MemSys"`
	HeapAlloc    uint64 `json:"HeapAlloc"`
	HeapSys      uint64 `json:"HeapSys"`
	HeapIdle     uint64 `json:"HeapIdle"`
	HeapInuse    uint64 `json:"HeapInuse"`
	HeapReleased uint64 `json:"HeapRealease"`
	NextGC       uint64 `json:"NextGC"`
	Goroutines   int    `json:"Goroutines"`
	UpTime       uint64 `json:"UpTime"`
	Time         string `json:"Time"`
}

type fileInfo struct {
	Name    string `json:"Name"`
	Size    int64  `string:"Size"`
	ModTime string `string:"ModTime"`
}

func infoHandler(_ *handler.Context, _ interface{}) (interface{}, error) {
	host, _ := os.Hostname()
	memory := &runtime.MemStats{}
	runtime.ReadMemStats(memory)
	rt := &runtimeInfo{
		NumCPU:       runtime.NumCPU(),
		Memory:       memory.Alloc,
		MemSys:       memory.Sys / 1024,
		HeapAlloc:    memory.HeapAlloc / 1024,
		HeapSys:      memory.HeapSys / 1024,
		HeapIdle:     memory.HeapIdle / 1024,
		HeapInuse:    memory.HeapInuse / 1024,
		HeapReleased: memory.HeapReleased / 1024,
		NextGC:       memory.NextGC / 1024,
		Goroutines:   runtime.NumGoroutine(),
		UpTime:       uint64(time.Since(startTime).Seconds()),
		Time:         time.Now().Format(time.RFC1123Z),
	}

	var files []fileInfo
	names := assets.AssetNames()
	sort.Strings(names)
	for _, name := range names {
		fileInfo := fileInfo{
			Name: name,
		}
		assetInfo, err := assets.AssetInfo(name)
		if nil == err {
			fileInfo.Size = assetInfo.Size()
			fileInfo.ModTime = assetInfo.ModTime().Format(time.RFC1123Z)
		}
		files = append(files, fileInfo)
	}

	return &info{
		Host:    host,
		Runtime: rt,
		Files:   files,
	}, nil
}

func NewInfoHandler(sessionStore *sessions.Store) http.HandlerFunc {
	return handler.CreateFromProcessFunc(infoHandler, handler.WithOnlyForSuperuser(sessionStore))
}
