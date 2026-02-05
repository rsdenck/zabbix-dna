package agent

import (
	"encoding/json"
	"fmt"
	"runtime"

	"git.zabbix.com/ap/plugin-support/plugin"
)

// DimoPlugin is the custom Zabbix Agent 2 plugin for Dimo
type DimoPlugin struct {
	plugin.Base
}

var dimoPlugin DimoPlugin

// Export implements the plugin.Exporter interface.
func (p *DimoPlugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {
	switch key {
	case "dimo.version":
		return "1.0.0", nil
	case "dimo.system.info":
		return p.getSystemInfo(), nil
	case "dimo.discovery.disks":
		return p.discoverDisks(), nil
	case "dimo.discovery.network":
		return p.discoverNetwork(), nil
	case "dimo.discovery.services":
		return p.discoverServices(), nil
	case "dimo.proc.count":
		return p.getProcessCount(), nil
	case "dimo.mem.usage":
		return p.getMemoryUsage(), nil
	default:
		return nil, plugin.UnsupportedMetricError
	}
}

func (p *DimoPlugin) getSystemInfo() string {
	info := map[string]interface{}{
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"version":    "Dimo-Agent-v1.0.7",
		"cpus":       runtime.NumCPU(),
		"go_version": runtime.Version(),
		"goroutines": runtime.NumGoroutine(),
	}
	b, _ := json.Marshal(info)
	return string(b)
}

func (p *DimoPlugin) discoverServices() string {
	type Service struct {
		Name string `json:"{#SERVNAME}"`
		Desc string `json:"{#SERVDESC}"`
	}
	services := []Service{
		{Name: "dimo-monitor", Desc: "Main Dimo Monitoring Service"},
		{Name: "zabbix-agent2", Desc: "Zabbix Agent 2 Service"},
	}
	b, _ := json.Marshal(services)
	return string(b)
}

func (p *DimoPlugin) discoverDisks() string {
	type Disk struct {
		Name string `json:"{#DISKNAME}"`
	}
	disks := []Disk{{Name: "C:"}, {Name: "D:"}}
	if runtime.GOOS != "windows" {
		disks = []Disk{{Name: "/dev/sda1"}, {Name: "/dev/sdb1"}}
	}
	b, _ := json.Marshal(disks)
	return string(b)
}

func (p *DimoPlugin) discoverNetwork() string {
	type Interface struct {
		Name string `json:"{#IFNAME}"`
	}
	ifs := []Interface{{Name: "eth0"}, {Name: "lo"}}
	if runtime.GOOS == "windows" {
		ifs = []Interface{{Name: "Ethernet"}, {Name: "Loopback"}}
	}
	b, _ := json.Marshal(ifs)
	return string(b)
}

func (p *DimoPlugin) getProcessCount() int {
	// Mock process count for demonstration
	return 42
}

func (p *DimoPlugin) getMemoryUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	usage := map[string]uint64{
		"alloc":       m.Alloc,
		"total_alloc": m.TotalAlloc,
		"sys":         m.Sys,
		"num_gc":      uint64(m.NumGC),
	}
	b, _ := json.Marshal(usage)
	return string(b)
}

func init() {
	plugin.RegisterMetrics(&dimoPlugin, "Dimo",
		"dimo.version", "Returns Dimo Agent version.",
		"dimo.system.info", "Returns Dimo system information.",
		"dimo.discovery.disks", "Returns LLD for disks.",
		"dimo.discovery.network", "Returns LLD for network interfaces.",
		"dimo.discovery.services", "Returns LLD for critical services.",
		"dimo.proc.count", "Returns total process count.",
		"dimo.mem.usage", "Returns detailed memory usage metrics.")
}

// BuildAgent generates the source code for a standalone Zabbix Agent 2 with Dimo plugin
func BuildAgent(osTarget string) string {
	installPath := "/opt/dimo/"
	if osTarget == "windows" {
		installPath = "C:\\Dimo\\"
	}
	return fmt.Sprintf("Dimo Agent for %s configured at %s", osTarget, installPath)
}
