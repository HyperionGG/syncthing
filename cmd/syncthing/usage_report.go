package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"runtime"
	"strings"

	"github.com/calmh/syncthing/model"
)

// Current version number of the usage report, for acceptance purposes. If
// fields are added or changed this integer must be incremented so that users
// are prompted for acceptance of the new report.
const usageReportVersion = 1

func reportData(m *model.Model) map[string]interface{} {
	res := make(map[string]interface{})
	res["uniqueID"] = strings.ToLower(certID([]byte(myID)))[:6]
	res["version"] = Version
	res["platform"] = runtime.GOOS + "-" + runtime.GOARCH
	res["numRepos"] = len(cfg.Repositories)
	res["numNodes"] = len(cfg.Nodes)

	var totFiles, maxFiles int
	var totBytes, maxBytes int64
	for _, repo := range cfg.Repositories {
		files, _, bytes := m.GlobalSize(repo.ID)
		totFiles += files
		totBytes += bytes
		if files > maxFiles {
			maxFiles = files
		}
		if bytes > maxBytes {
			maxBytes = bytes
		}
	}

	res["totFiles"] = totFiles
	res["repoMaxFiles"] = maxFiles
	res["totMiB"] = totBytes / 1024 / 1024
	res["repoMaxMiB"] = maxBytes / 1024 / 1024

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	res["memoryUsageMiB"] = mem.Sys / 1024 / 1024

	return res
}

func sendUsageRport(m *model.Model) error {
	d := reportData(m)
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(d)
	_, err := http.Post("https://data.syncthing.net:8443/newdata", "application/json", &b)
	return err
}
