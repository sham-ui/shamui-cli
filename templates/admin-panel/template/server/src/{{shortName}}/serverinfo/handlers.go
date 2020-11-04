package serverinfo

import (
	"assets"
	"encoding/json"
	"net/http"
	"os"
	"{{ shortName }}/sessions"
	"runtime"
	"sort"
	"time"
)

var (
	startTime = time.Now()
)

type resDetails struct {
	Status   string
	Messages []string
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	if session == nil {
		msg := resDetails{
			Status:   "Expired session or cookie",
			Messages: []string{"Session Expired.  Log out and log back in."},
		}
		json.NewEncoder(w).Encode(msg)
		return
	}
	if !session.IsSuperuser {
		msg := resDetails{
			Status:   "Only superuser can get server info",
			Messages: []string{"Only superuser can get server info"},
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(msg)
		return
	}
	info := newInfo()
	json.NewEncoder(w).Encode(info)
	w.WriteHeader(http.StatusOK)
}

type info struct {
	Host    string       `json:"Host"`
	Runtime *RuntimeInfo `json:"Runtime"`
	Files   []FileInfo   `json:"Files"`
}

// RuntimeInfo defines runtime part of service information
type RuntimeInfo struct {
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

type FileInfo struct {
	Name string `json:"Name"`
	Size int64  `string:"Size"`
	ModTime string `string:"ModTime"`
}

func newInfo() *info {
	host, _ := os.Hostname()
	memory := &runtime.MemStats{}
	runtime.ReadMemStats(memory)
	rt := &RuntimeInfo{
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

	var files []FileInfo
	names := assets.AssetNames()
	sort.Strings(names)
	for _, name := range names {
		fileInfo := FileInfo{
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
	}
}
