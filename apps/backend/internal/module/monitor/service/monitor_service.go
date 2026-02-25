package service

import (
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemInfo 系统信息
type SystemInfo struct {
	Hostname   string `json:"hostname"`
	OS         string `json:"os"`
	Platform   string `json:"platform"`
	KernelVer  string `json:"kernel_version"`
	Uptime     uint64 `json:"uptime"`
	UptimeDesc string `json:"uptime_desc"` // 可读的运行时间
}

// CPUInfo CPU信息
type CPUInfo struct {
	Count       int     `json:"count"`        // CPU核心数
	UsedPercent float64 `json:"used_percent"` // CPU使用率
	ModelName   string  `json:"model_name"`   // CPU型号
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total       uint64  `json:"total"`        // 总内存(字节)
	Used        uint64  `json:"used"`         // 已用内存(字节)
	Available   uint64  `json:"available"`    // 可用内存(字节)
	UsedPercent float64 `json:"used_percent"` // 内存使用率
	TotalGB     float64 `json:"total_gb"`     // 总内存(GB)
	UsedGB      float64 `json:"used_gb"`      // 已用内存(GB)
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Total       uint64  `json:"total"`        // 总空间(字节)
	Used        uint64  `json:"used"`         // 已用空间(字节)
	Free        uint64  `json:"free"`         // 可用空间(字节)
	UsedPercent float64 `json:"used_percent"` // 磁盘使用率
	TotalGB     float64 `json:"total_gb"`     // 总空间(GB)
	UsedGB      float64 `json:"used_gb"`      // 已用空间(GB)
}

// LoadInfo 系统负载
type LoadInfo struct {
	Load1  float64 `json:"load1"`  // 1分钟负载
	Load5  float64 `json:"load5"`  // 5分钟负载
	Load15 float64 `json:"load15"` // 15分钟负载
}

// RuntimeInfo 运行时信息
type RuntimeInfo struct {
	GoVersion    string `json:"go_version"`
	GoroutineNum int    `json:"goroutine_num"`
	HeapAlloc    uint64 `json:"heap_alloc"`
	HeapSys      uint64 `json:"heap_sys"`
	HeapIdle     uint64 `json:"heap_idle"`
	HeapInuse    uint64 `json:"heap_inuse"`
	StackInuse   uint64 `json:"stack_inuse"`
	StackSys     uint64 `json:"stack_sys"`
	NumGC        uint32 `json:"num_gc"`
}

// SystemStats 系统统计数据
type SystemStats struct {
	Timestamp int64        `json:"timestamp"`
	System    SystemInfo   `json:"system"`
	CPU       CPUInfo      `json:"cpu"`
	Memory    MemoryInfo   `json:"memory"`
	Disk      DiskInfo     `json:"disk"`
	Load      LoadInfo     `json:"load"`
	Runtime   RuntimeInfo  `json:"runtime"`
}

// cachedSystemInfo 缓存的系统信息
type cachedSystemInfo struct {
	info      SystemInfo
	timestamp time.Time
}

// MonitorService 监控服务
type MonitorService struct {
	// 系统信息缓存（变化不频繁）
	systemInfoCache     cachedSystemInfo
	systemInfoCacheMu   sync.RWMutex
	systemInfoCacheTTL  time.Duration
}

// NewMonitorService 创建监控服务实例
func NewMonitorService() *MonitorService {
	return &MonitorService{
		systemInfoCacheTTL: 5 * time.Minute, // 系统信息缓存5分钟
	}
}

// GetSystemStats 获取系统统计数据
func (s *MonitorService) GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{
		Timestamp: time.Now().Unix(),
	}

	// 使用 WaitGroup 并发获取数据
	var wg sync.WaitGroup
	var mu sync.Mutex

	// 并发获取各项数据
	wg.Add(6)

	// 系统信息
	go func() {
		defer wg.Done()
		systemInfo := s.getSystemInfo()
		mu.Lock()
		stats.System = systemInfo
		mu.Unlock()
	}()

	// CPU信息
	go func() {
		defer wg.Done()
		cpuInfo := s.getCPUInfo()
		mu.Lock()
		stats.CPU = cpuInfo
		mu.Unlock()
	}()

	// 内存信息
	go func() {
		defer wg.Done()
		memInfo := s.getMemoryInfo()
		mu.Lock()
		stats.Memory = memInfo
		mu.Unlock()
	}()

	// 磁盘信息
	go func() {
		defer wg.Done()
		diskInfo := s.getDiskInfo()
		mu.Lock()
		stats.Disk = diskInfo
		mu.Unlock()
	}()

	// 系统负载
	go func() {
		defer wg.Done()
		loadInfo := s.getLoadInfo()
		mu.Lock()
		stats.Load = loadInfo
		mu.Unlock()
	}()

	// 运行时信息
	go func() {
		defer wg.Done()
		runtimeInfo := s.getRuntimeInfo()
		mu.Lock()
		stats.Runtime = runtimeInfo
		mu.Unlock()
	}()

	wg.Wait()

	return stats, nil
}

// getSystemInfo 获取系统信息（带缓存）
func (s *MonitorService) getSystemInfo() SystemInfo {
	// 检查缓存
	s.systemInfoCacheMu.RLock()
	if time.Since(s.systemInfoCache.timestamp) < s.systemInfoCacheTTL {
		info := s.systemInfoCache.info
		// 更新运行时间（这个需要实时获取）
		if hostInfo, err := host.Info(); err == nil {
			info.Uptime = hostInfo.Uptime
			info.UptimeDesc = formatUptime(hostInfo.Uptime)
		}
		s.systemInfoCacheMu.RUnlock()
		return info
	}
	s.systemInfoCacheMu.RUnlock()

	// 获取新的系统信息
	info := SystemInfo{}

	hostInfo, err := host.Info()
	if err == nil {
		info.Hostname = hostInfo.Hostname
		info.OS = hostInfo.OS
		info.Platform = hostInfo.Platform
		info.KernelVer = hostInfo.KernelVersion
		info.Uptime = hostInfo.Uptime
		info.UptimeDesc = formatUptime(hostInfo.Uptime)
	}

	// 更新缓存
	s.systemInfoCacheMu.Lock()
	s.systemInfoCache = cachedSystemInfo{
		info:      info,
		timestamp: time.Now(),
	}
	s.systemInfoCacheMu.Unlock()

	return info
}

// getCPUInfo 获取CPU信息
func (s *MonitorService) getCPUInfo() CPUInfo {
	info := CPUInfo{}

	// 获取CPU核心数
	info.Count = runtime.NumCPU()

	// 获取CPU使用率（使用较短的超时时间）
	percent, err := cpu.Percent(500*time.Millisecond, false)
	if err == nil && len(percent) > 0 {
		info.UsedPercent = percent[0]
	}

	// 获取CPU型号
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		info.ModelName = cpuInfo[0].ModelName
	}

	return info
}

// getMemoryInfo 获取内存信息
func (s *MonitorService) getMemoryInfo() MemoryInfo {
	info := MemoryInfo{}

	memInfo, err := mem.VirtualMemory()
	if err == nil {
		info.Total = memInfo.Total
		info.Used = memInfo.Used
		info.Available = memInfo.Available
		info.UsedPercent = memInfo.UsedPercent
		// 转换为GB
		info.TotalGB = float64(memInfo.Total) / 1024 / 1024 / 1024
		info.UsedGB = float64(memInfo.Used) / 1024 / 1024 / 1024
	}

	return info
}

// getDiskInfo 获取磁盘信息
func (s *MonitorService) getDiskInfo() DiskInfo {
	info := DiskInfo{}

	diskInfo, err := disk.Usage("/")
	if err == nil {
		info.Total = diskInfo.Total
		info.Used = diskInfo.Used
		info.Free = diskInfo.Free
		info.UsedPercent = diskInfo.UsedPercent
		// 转换为GB
		info.TotalGB = float64(diskInfo.Total) / 1024 / 1024 / 1024
		info.UsedGB = float64(diskInfo.Used) / 1024 / 1024 / 1024
	}

	return info
}

// getLoadInfo 获取系统负载
func (s *MonitorService) getLoadInfo() LoadInfo {
	info := LoadInfo{}

	loadInfo, err := load.Avg()
	if err == nil {
		info.Load1 = loadInfo.Load1
		info.Load5 = loadInfo.Load5
		info.Load15 = loadInfo.Load15
	}

	return info
}

// getRuntimeInfo 获取运行时信息
func (s *MonitorService) getRuntimeInfo() RuntimeInfo {
	info := RuntimeInfo{
		GoVersion: runtime.Version(),
	}

	// 获取Goroutine数量
	info.GoroutineNum = runtime.NumGoroutine()

	// 获取内存统计
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	info.HeapAlloc = m.HeapAlloc
	info.HeapSys = m.HeapSys
	info.HeapIdle = m.HeapIdle
	info.HeapInuse = m.HeapInuse
	info.StackInuse = m.StackInuse
	info.StackSys = m.StackSys
	info.NumGC = m.NumGC

	return info
}

// formatUptime 格式化运行时间为可读字符串
func formatUptime(uptime uint64) string {
	seconds := uptime
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	var parts []string
	if days > 0 {
		parts = append(parts, formatNumber(days, "天"))
	}
	if hours%24 > 0 {
		parts = append(parts, formatNumber(hours%24, "小时"))
	}
	if minutes%60 > 0 || len(parts) == 0 {
		parts = append(parts, formatNumber(minutes%60, "分钟"))
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " "
		}
		result += part
	}
	return result
}

// formatNumber 格式化数字和单位
func formatNumber(n uint64, unit string) string {
	return strconv.FormatUint(n, 10) + unit
}
