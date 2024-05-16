package metrics

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/talkincode/logsight/common/web"
	"github.com/talkincode/logsight/webserver"
)

func InitRouter() {

	webserver.GET("/admin/metrics/system/hostname", func(c echo.Context) error {
		hinfo, err := host.Info()
		_host := "unknow"
		if err == nil {
			_host = hinfo.Hostname
		}
		return c.Render(http.StatusOK, "metrics", web.NewMetrics("mdi mdi-server", _host, "主机名"))
	})

	webserver.GET("/admin/metrics/system/os", func(c echo.Context) error {
		hinfo, err := host.Info()
		_os := "unknow"
		if err == nil {
			_os = hinfo.OS
		}
		return c.Render(http.StatusOK, "metrics",
			web.NewMetrics("mdi mdi-lifebuoy", _os, "操作系统"))
	})

	webserver.GET("/admin/metrics/system/diskuse", func(c echo.Context) error {
		partitions, err := disk.Partitions(true)
		if err != nil {
			return c.Render(http.StatusOK, "metrics",
				web.NewMetrics("mdi mdi-lifebuoy", "N/A", "Unknow"))
		}

		var totalUsed uint64
		var totalSize uint64

		for _, p := range partitions {
			// 过滤掉虚拟文件系统和 LVM 逻辑卷
			if strings.HasPrefix(p.Device, "/dev") && !strings.HasPrefix(p.Fstype, "tmpfs") && !strings.HasPrefix(p.Fstype, "devtmpfs") && !strings.Contains(p.Device, "mapper") {
				usage, err := disk.Usage(p.Mountpoint)
				if err != nil {
					continue
				}
				totalUsed += usage.Used
				totalSize += usage.Total
			}
		}

		var totalUsagePercent float64
		if totalSize > 0 {
			totalUsagePercent = (float64(totalUsed) / float64(totalSize)) * 100
		}

		totalSizeGB := totalSize / (1024 * 1024 * 1024)

		return c.Render(http.StatusOK, "metrics",
			web.NewMetrics("mdi mdi-harddisk", fmt.Sprintf("%.2f%%", totalUsagePercent),
				fmt.Sprintf("磁盘总占用 (总大小: %d G)", totalSizeGB)))

	})

	webserver.GET("/admin/metrics/system/cpuusage", func(c echo.Context) error {
		_cpuuse, err := cpu.Percent(0, false)
		_cpucount, _ := cpu.Counts(false)
		if err != nil {
			_cpuuse = []float64{0}
		}
		return c.Render(http.StatusOK, "metrics",
			web.NewMetrics("mdi mdi-circle-slice-2",
				fmt.Sprintf("%.2f%%", _cpuuse[0]),
				fmt.Sprintf("Cpu %d Core", _cpucount)))
	})

	webserver.GET("/admin/metrics/system/main/cpuusage", func(c echo.Context) error {
		var cpuuse float64
		p, err := process.NewProcess(int32(os.Getpid()))
		if err != nil {
			cpuuse, _ = p.CPUPercent()
		}
		return c.Render(http.StatusOK, "metrics",
			web.NewMetrics("mdi mdi-circle-slice-2", fmt.Sprintf("%.2f%%", cpuuse), "主程序Cpu负载"))
	})

	webserver.GET("/admin/metrics/system/memusage", func(c echo.Context) error {
		_meminfo, err := mem.VirtualMemory()
		_usage := 0.0
		_total := uint64(0)
		if err == nil {
			_usage = _meminfo.UsedPercent
			_total = _meminfo.Total / (1000 * 1000 * 1000)
		}
		return c.Render(http.StatusOK, "metrics",
			web.NewMetrics("mdi mdi-memory", fmt.Sprintf("%.2f%%", _usage),
				fmt.Sprintf("Memory Total: %d G", _total)))
	})

	webserver.GET("/admin/metrics/system/main/memusage", func(c echo.Context) error {
		var memuse uint64
		p, err := process.NewProcess(int32(os.Getpid()))
		if err == nil {
			meminfo, err := p.MemoryInfo()
			if err == nil {
				memuse = meminfo.RSS / 1024 / 1024
			}
		}

		return c.Render(http.StatusOK, "metrics",
			web.NewMetrics("mdi mdi-memory", fmt.Sprintf("%d MB", memuse), "主程序内存使用"))
	})

	webserver.GET("/admin/metrics/system/uptime", func(c echo.Context) error {
		hinfo, err := host.Info()
		_hour := uint64(0)
		if err == nil {
			_hour = hinfo.Uptime
		}
		return c.Render(http.StatusOK, "metrics",
			web.NewMetrics("mdi mdi-clock",
				fmt.Sprintf("%.1f Hour",
					float64(_hour)/float64(3600)), "运行时长"))
	})

	webserver.GET("/admin/metrics/unknow", func(c echo.Context) error {
		return c.Render(http.StatusOK, "metrics",
			web.NewMetrics("mdi mdi-lifebuoy", "N/A", "Unknow"))
	})

}
