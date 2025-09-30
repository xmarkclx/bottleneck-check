package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// SystemMetrics holds all the system performance data
type SystemMetrics struct {
	CPUUsage     float64
	LoadAverage  [3]float64
	MemoryUsed   uint64
	MemoryTotal  uint64
	SwapUsed     uint64
	SwapTotal    uint64
	MemPressure  string
	GPUUsage     float64
	CPUCores     int
	CPUModel     string
	MemorySpeed  string
	GPUModel     string
}

// Recommendation represents an upgrade suggestion
type Recommendation struct {
	Component string
	Severity  string // "LOW", "MEDIUM", "HIGH", "CRITICAL"
	Reason    string
	Suggestion string
	Color     string
}

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

func main() {
	fmt.Printf("%s%s🔍 System Bottleneck Monitor%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("═══════════════════════════════════\n")
	fmt.Printf("Running continuous system monitoring...\n")
	fmt.Printf("Press %s[Enter]%s for menu options\n\n", ColorYellow, ColorReset)

	// Start continuous monitoring
	runContinuousMonitor()
}

// Global variables for continuous monitoring
var (
	lastMetrics *SystemMetrics
	lastRecommendations []Recommendation
	monitoringActive = true
	lastUpdate time.Time
)

func runContinuousMonitor() {
	// Channel to handle user input
	inputChan := make(chan string)
	ticker := time.NewTicker(10 * time.Second) // Update every 10 seconds
	defer ticker.Stop()

	// Start goroutine to handle user input
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			if scanner.Scan() {
				inputChan <- scanner.Text()
			}
		}
	}()

	// Initial data collection
	updateSystemData()
	displayStatus()

	// Main monitoring loop
	for monitoringActive {
		select {
		case <-ticker.C:
			// Update system data every 10 seconds
			updateSystemData()
			displayStatus()

		case input := <-inputChan:
			handleUserInput(strings.TrimSpace(strings.ToLower(input)))
		}
	}
}

func updateSystemData() {
	metrics, err := collectSystemMetrics()
	if err != nil {
		fmt.Printf("%sError collecting metrics: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	lastMetrics = metrics
	lastRecommendations = analyzeSystem(metrics)
	lastUpdate = time.Now()
}

func displayStatus() {
	if lastMetrics == nil {
		return
	}

	// Clear screen (move cursor to top)
	fmt.Print("\033[H\033[2J")

	// Header
	fmt.Printf("%s%s🔍 System Bottleneck Monitor%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("═══════════════════════════════════\n")
	fmt.Printf("Last updated: %s | Press %s[Enter]%s for menu\n\n", 
		lastUpdate.Format("15:04:05"), ColorYellow, ColorReset)

	// Quick status indicators
	displayQuickStatus(lastMetrics)

	// Show detailed system status
	fmt.Printf("\n")
	displaySystemStatus(lastMetrics)

	// Show detailed recommendations by default
	displayRecommendations(lastRecommendations)

	// Status bar
	fmt.Printf("\n%s──────────────────────────────────────────────────%s\n", ColorBlue, ColorReset)
	fmt.Printf("Monitoring active... Press %s[Enter]%s for menu\n", ColorGreen, ColorReset)
}

func displayQuickStatus(metrics *SystemMetrics) {
	// CPU Status with color coding
	cpuColor := ColorGreen
	if metrics.CPUUsage > 80 {
		cpuColor = ColorRed
	} else if metrics.CPUUsage > 60 {
		cpuColor = ColorYellow
	}

	// Memory Status with color coding
	memUsagePercent := float64(metrics.MemoryUsed) / float64(metrics.MemoryTotal) * 100
	memColor := ColorGreen
	if memUsagePercent > 90 {
		memColor = ColorRed
	} else if memUsagePercent > 75 {
		memColor = ColorYellow
	}

	// Load Average Status
	loadColor := ColorGreen
	if metrics.LoadAverage[0] > float64(metrics.CPUCores)*1.5 {
		loadColor = ColorRed
	} else if metrics.LoadAverage[0] > float64(metrics.CPUCores) {
		loadColor = ColorYellow
	}

	// Display compact status
	fmt.Printf("📊 %sQuick Status%s\n", ColorBold, ColorReset)
	fmt.Printf("─────────────\n")
	fmt.Printf("CPU: %s%.1f%%%s | Load: %s%.2f%s | Memory: %s%.1f%%%s", 
		cpuColor, metrics.CPUUsage, ColorReset,
		loadColor, metrics.LoadAverage[0], ColorReset,
		memColor, memUsagePercent, ColorReset)

	if metrics.SwapUsed > 0 {
		swapGB := float64(metrics.SwapUsed) / (1024 * 1024 * 1024)
		swapColor := ColorYellow
		if swapGB > 2 {
			swapColor = ColorRed
		}
		fmt.Printf(" | Swap: %s%.1fGB%s", swapColor, swapGB, ColorReset)
	}
	fmt.Println()
}

func showCriticalAlerts(recommendations []Recommendation) {
	criticalCount := 0
	highCount := 0

	for _, rec := range recommendations {
		if rec.Severity == "CRITICAL" {
			criticalCount++
		} else if rec.Severity == "HIGH" {
			highCount++
		}
	}

	if criticalCount > 0 || highCount > 0 {
		fmt.Printf("\n🚨 %sActive Alerts%s\n", ColorBold, ColorReset)
		fmt.Printf("──────────────\n")
		if criticalCount > 0 {
			fmt.Printf("%s● CRITICAL: %d issue(s) need immediate attention%s\n", ColorRed, criticalCount, ColorReset)
		}
		if highCount > 0 {
			fmt.Printf("%s● HIGH: %d issue(s) affecting performance%s\n", ColorYellow, highCount, ColorReset)
		}
	} else {
		fmt.Printf("\n✅ %sSystem Status: Good%s\n", ColorGreen, ColorReset)
	}
}

func handleUserInput(input string) {
	switch input {
	case "":
		// Show menu on Enter
		showMenu()
	case "q", "quit", "exit":
		fmt.Printf("\n%sExiting monitor...%s\n", ColorCyan, ColorReset)
		monitoringActive = false
	case "s", "status":
		// Force refresh display
		updateSystemData()
		displayStatus()
	case "a", "advice":
		// Detailed advice is now shown by default, so just refresh
		updateSystemData()
		displayStatus()
	case "d", "details":
		showDetailedSystemInfo()
	case "h", "help":
		showHelp()
	case "c", "clear":
		// Clear screen and refresh
		displayStatus()
	default:
		fmt.Printf("%sUnknown command: '%s'. Press Enter for menu.%s\n", ColorRed, input, ColorReset)
	}
}

func showMenu() {
	fmt.Printf("\n%s⚙️  Menu Options%s\n", ColorBold, ColorReset)
	fmt.Printf("───────────────\n")
	fmt.Printf("%s[a]%s - Refresh advice and recommendations\n", ColorGreen, ColorReset)
	fmt.Printf("%s[d]%s - Show detailed system information\n", ColorGreen, ColorReset)
	fmt.Printf("%s[s]%s - Refresh system status\n", ColorGreen, ColorReset)
	fmt.Printf("%s[c]%s - Clear screen\n", ColorGreen, ColorReset)
	fmt.Printf("%s[h]%s - Show help\n", ColorGreen, ColorReset)
	fmt.Printf("%s[q]%s - Quit monitor\n", ColorGreen, ColorReset)
	fmt.Printf("───────────────\n")
	fmt.Printf("Enter command: ")
}

// showDetailedAdvice is now integrated into the main display
// This function is kept for compatibility but just refreshes the main display
func showDetailedAdvice() {
	updateSystemData()
	displayStatus()
}

func showDetailedSystemInfo() {
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Printf("%s%s📊 Detailed System Information%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("═══════════════════════════════════\n\n")

	if lastMetrics != nil {
		// Show detailed system information
		fmt.Printf("%sSystem Hardware:%s\n", ColorBlue, ColorReset)
		fmt.Printf("  CPU: %s (%d cores)\n", lastMetrics.CPUModel, lastMetrics.CPUCores)
		if lastMetrics.GPUModel != "" {
			fmt.Printf("  GPU: %s\n", lastMetrics.GPUModel)
		}
		fmt.Printf("  Total RAM: %.1fGB\n", float64(lastMetrics.MemoryTotal)/(1024*1024*1024))

		fmt.Printf("\n%sPerformance Metrics:%s\n", ColorPurple, ColorReset)
		fmt.Printf("  CPU Usage: %.1f%%\n", lastMetrics.CPUUsage)
		fmt.Printf("  Load Averages: %.2f (1m), %.2f (5m), %.2f (15m)\n", 
			lastMetrics.LoadAverage[0], lastMetrics.LoadAverage[1], lastMetrics.LoadAverage[2])
		fmt.Printf("  Memory Usage: %.1f%% (%.1fGB used)\n", 
			float64(lastMetrics.MemoryUsed)/float64(lastMetrics.MemoryTotal)*100,
			float64(lastMetrics.MemoryUsed)/(1024*1024*1024))
		if lastMetrics.SwapUsed > 0 {
			fmt.Printf("  Swap Usage: %.1fGB\n", float64(lastMetrics.SwapUsed)/(1024*1024*1024))
		}
		fmt.Printf("  Memory Pressure: %s\n", lastMetrics.MemPressure)

		// Show uptime
		cmd := exec.Command("uptime")
		if output, err := cmd.Output(); err == nil {
			fmt.Printf("\n%sSystem Uptime:%s\n", ColorGreen, ColorReset)
			fmt.Printf("  %s\n", strings.TrimSpace(string(output)))
		}
	} else {
		fmt.Printf("%sNo system data available yet.%s\n", ColorRed, ColorReset)
	}

	fmt.Printf("\n%sPress Enter to return to monitor...%s", ColorYellow, ColorReset)
	fmt.Scanln() // Wait for user input
	displayStatus() // Return to main display
}

func showHelp() {
	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Printf("%s%s❓ Help & Usage Guide%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("══════════════════════\n\n")

	fmt.Printf("%sWhat This Tool Does:%s\n", ColorBlue, ColorReset)
	fmt.Printf("• Continuously monitors CPU, Memory, and GPU performance\n")
	fmt.Printf("• Provides real-time bottleneck detection\n")
	fmt.Printf("• Shows detailed upgrade recommendations by default\n")
	fmt.Printf("• Updates every 10 seconds automatically\n\n")

	fmt.Printf("%sStatus Indicators:%s\n", ColorPurple, ColorReset)
	fmt.Printf("%s• Green%s - Good performance\n", ColorGreen, ColorReset)
	fmt.Printf("%s• Yellow%s - Moderate usage/warning\n", ColorYellow, ColorReset)
	fmt.Printf("%s• Red%s - High usage/critical issue\n\n", ColorRed, ColorReset)

	fmt.Printf("%sRecommendation Levels:%s\n", ColorYellow, ColorReset)
	fmt.Printf("🚨 CRITICAL - Immediate action required\n")
	fmt.Printf("⚠️  HIGH - Should address soon\n")
	fmt.Printf("📋 MEDIUM - Consider for future upgrades\n")
	fmt.Printf("💡 LOW - Optional improvements\n\n")

	fmt.Printf("%sTips for Best Results:%s\n", ColorGreen, ColorReset)
	fmt.Printf("• Let it run for a few minutes to see usage patterns\n")
	fmt.Printf("• Detailed advice is shown automatically - watch for changes\n")
	fmt.Printf("• Use during your typical workload for accurate assessment\n")
	fmt.Printf("• Address CRITICAL issues first for best performance gains\n")

	fmt.Printf("\n%sPress Enter to return to monitor...%s", ColorYellow, ColorReset)
	fmt.Scanln() // Wait for user input
	displayStatus() // Return to main display
}

func collectSystemMetrics() (*SystemMetrics, error) {
	metrics := &SystemMetrics{}

	// Get CPU information
	metrics.CPUCores = runtime.NumCPU()
	cpuModel, _ := getCPUModel()
	metrics.CPUModel = cpuModel

	// Get CPU usage (average over 2 seconds)
	cpuUsage, _ := getCPUUsage()
	metrics.CPUUsage = cpuUsage

	// Get load averages
	loadAvg, _ := getLoadAverages()
	metrics.LoadAverage = loadAvg

	// Get memory information
	memInfo, _ := getMemoryInfo()
	metrics.MemoryUsed = memInfo.Used
	metrics.MemoryTotal = memInfo.Total
	metrics.SwapUsed = memInfo.SwapUsed
	metrics.SwapTotal = memInfo.SwapTotal

	// Get memory pressure
	memPressure, _ := getMemoryPressure()
	metrics.MemPressure = memPressure

	// Get GPU information
	gpuModel, _ := getGPUModel()
	metrics.GPUModel = gpuModel

	return metrics, nil
}

func getCPUModel() (string, error) {
	cmd := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getCPUUsage() (float64, error) {
	// Sample CPU usage over 2 seconds using iostat
	cmd := exec.Command("iostat", "-c", "2", "-n", "2")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], "avg-cpu") {
			continue
		}
		fields := strings.Fields(lines[i])
		if len(fields) >= 6 {
			// iostat format: %user %nice %system %interrupt %idle
			idle, err := strconv.ParseFloat(fields[len(fields)-1], 64)
			if err == nil {
				return 100 - idle, nil
			}
		}
	}
	return 0, fmt.Errorf("could not parse CPU usage")
}

func getLoadAverages() ([3]float64, error) {
	cmd := exec.Command("uptime")
	output, err := cmd.Output()
	if err != nil {
		return [3]float64{}, err
	}

	// Parse load averages from uptime output
	uptimeStr := string(output)
	loadIndex := strings.Index(uptimeStr, "load averages:")
	if loadIndex == -1 {
		return [3]float64{}, fmt.Errorf("could not find load averages")
	}

	loadStr := strings.TrimSpace(uptimeStr[loadIndex+14:])
	loadParts := strings.Split(loadStr, " ")
	
	var loads [3]float64
	for i := 0; i < 3 && i < len(loadParts); i++ {
		load, err := strconv.ParseFloat(strings.TrimSpace(loadParts[i]), 64)
		if err == nil {
			loads[i] = load
		}
	}

	return loads, nil
}

type MemoryInfo struct {
	Total    uint64
	Used     uint64
	SwapUsed uint64
	SwapTotal uint64
}

func getMemoryInfo() (*MemoryInfo, error) {
	info := &MemoryInfo{}

	// Get physical memory info
	cmd := exec.Command("vm_stat")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse vm_stat output
	lines := strings.Split(string(output), "\n")
	pageSize := uint64(4096) // Default page size for macOS

	var free, active, inactive, wired, compressed uint64

	for _, line := range lines {
		if strings.Contains(line, "page size of") {
			// Extract page size
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "of" && i+1 < len(parts) {
					if ps, err := strconv.ParseUint(parts[i+1], 10, 64); err == nil {
						pageSize = ps
					}
				}
			}
		} else if strings.Contains(line, "Pages free:") {
			if val := extractNumber(line); val > 0 {
				free = val
			}
		} else if strings.Contains(line, "Pages active:") {
			if val := extractNumber(line); val > 0 {
				active = val
			}
		} else if strings.Contains(line, "Pages inactive:") {
			if val := extractNumber(line); val > 0 {
				inactive = val
			}
		} else if strings.Contains(line, "Pages wired down:") {
			if val := extractNumber(line); val > 0 {
				wired = val
			}
		} else if strings.Contains(line, "Pages occupied by compressor:") {
			if val := extractNumber(line); val > 0 {
				compressed = val
			}
		}
	}

	totalPages := free + active + inactive + wired + compressed
	usedPages := active + inactive + wired + compressed

	info.Total = totalPages * pageSize
	info.Used = usedPages * pageSize

	// Get swap info
	cmd = exec.Command("sysctl", "-n", "vm.swapusage")
	output, err = cmd.Output()
	if err == nil {
		swapStr := string(output)
		// Parse format like "total = 2048.00M  used = 0.00M  free = 2048.00M  (encrypted)"
		if strings.Contains(swapStr, "used = ") {
			parts := strings.Split(swapStr, "used = ")
			if len(parts) > 1 {
				usedPart := strings.Fields(parts[1])[0]
				if swapUsed := parseMemorySize(usedPart); swapUsed > 0 {
					info.SwapUsed = swapUsed
				}
			}
		}
		if strings.Contains(swapStr, "total = ") {
			parts := strings.Split(swapStr, "total = ")
			if len(parts) > 1 {
				totalPart := strings.Fields(parts[1])[0]
				if swapTotal := parseMemorySize(totalPart); swapTotal > 0 {
					info.SwapTotal = swapTotal
				}
			}
		}
	}

	return info, nil
}

func extractNumber(line string) uint64 {
	parts := strings.Fields(line)
	for _, part := range parts {
		// Remove trailing period and parse
		clean := strings.TrimSuffix(part, ".")
		if val, err := strconv.ParseUint(clean, 10, 64); err == nil {
			return val
		}
	}
	return 0
}

func parseMemorySize(sizeStr string) uint64 {
	sizeStr = strings.TrimSpace(sizeStr)
	if len(sizeStr) == 0 {
		return 0
	}

	// Get the numeric part and unit
	unit := sizeStr[len(sizeStr)-1:]
	numStr := sizeStr[:len(sizeStr)-1]
	
	size, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}

	switch strings.ToUpper(unit) {
	case "K":
		return uint64(size * 1024)
	case "M":
		return uint64(size * 1024 * 1024)
	case "G":
		return uint64(size * 1024 * 1024 * 1024)
	case "T":
		return uint64(size * 1024 * 1024 * 1024 * 1024)
	default:
		return uint64(size)
	}
}

func getMemoryPressure() (string, error) {
	cmd := exec.Command("memory_pressure")
	output, err := cmd.Output()
	if err != nil {
		return "unknown", nil // memory_pressure might not be available
	}

	pressure := strings.TrimSpace(string(output))
	if strings.Contains(strings.ToLower(pressure), "normal") {
		return "normal", nil
	} else if strings.Contains(strings.ToLower(pressure), "warn") {
		return "warning", nil
	} else if strings.Contains(strings.ToLower(pressure), "urgent") || strings.Contains(strings.ToLower(pressure), "critical") {
		return "critical", nil
	}
	return pressure, nil
}

func getGPUModel() (string, error) {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType", "-json")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Simple parsing - look for GPU model in the output
	outputStr := string(output)
	if strings.Contains(outputStr, "chipset_model") {
		lines := strings.Split(outputStr, "\n")
		for _, line := range lines {
			if strings.Contains(line, "chipset_model") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					model := strings.Trim(strings.TrimSpace(parts[1]), `",`)
					return model, nil
				}
			}
		}
	}
	return "Unknown GPU", nil
}

func displaySystemStatus(metrics *SystemMetrics) {
	fmt.Printf("%s📊 Current System Status%s\n", ColorBold, ColorReset)
	fmt.Printf("─────────────────────────\n")
	
	// CPU Status
	fmt.Printf("%sCPU:%s %s (%d cores)\n", ColorBlue, ColorReset, metrics.CPUModel, metrics.CPUCores)
	fmt.Printf("  Usage: %.1f%%\n", metrics.CPUUsage)
	fmt.Printf("  Load Average: %.2f, %.2f, %.2f\n", metrics.LoadAverage[0], metrics.LoadAverage[1], metrics.LoadAverage[2])
	
	// Memory Status
	memUsagePercent := float64(metrics.MemoryUsed) / float64(metrics.MemoryTotal) * 100
	fmt.Printf("%sMemory:%s %.1fGB used / %.1fGB total (%.1f%%)\n", 
		ColorPurple, ColorReset,
		float64(metrics.MemoryUsed)/(1024*1024*1024),
		float64(metrics.MemoryTotal)/(1024*1024*1024),
		memUsagePercent)
	
	if metrics.SwapUsed > 0 {
		fmt.Printf("  Swap: %.1fMB used / %.1fMB total\n",
			float64(metrics.SwapUsed)/(1024*1024),
			float64(metrics.SwapTotal)/(1024*1024))
	}
	fmt.Printf("  Memory Pressure: %s\n", metrics.MemPressure)
	
	// GPU Status
	if metrics.GPUModel != "" {
		fmt.Printf("%sGPU:%s %s\n", ColorYellow, ColorReset, metrics.GPUModel)
	}
	
	fmt.Println()
}

func analyzeSystem(metrics *SystemMetrics) []Recommendation {
	var recommendations []Recommendation

	// Analyze CPU
	cpuRecommendations := analyzeCPU(metrics)
	recommendations = append(recommendations, cpuRecommendations...)

	// Analyze Memory
	memRecommendations := analyzeMemory(metrics)
	recommendations = append(recommendations, memRecommendations...)

	// Analyze GPU
	gpuRecommendations := analyzeGPU(metrics)
	recommendations = append(recommendations, gpuRecommendations...)

	return recommendations
}

func analyzeCPU(metrics *SystemMetrics) []Recommendation {
	var recommendations []Recommendation

	// Check CPU usage
	if metrics.CPUUsage > 90 {
		recommendations = append(recommendations, Recommendation{
			Component:  "CPU",
			Severity:   "CRITICAL",
			Reason:     fmt.Sprintf("CPU usage is very high (%.1f%%)", metrics.CPUUsage),
			Suggestion: "Consider upgrading to a faster CPU or adding more cores. Close unnecessary applications.",
			Color:      ColorRed,
		})
	} else if metrics.CPUUsage > 70 {
		recommendations = append(recommendations, Recommendation{
			Component:  "CPU",
			Severity:   "HIGH",
			Reason:     fmt.Sprintf("CPU usage is high (%.1f%%)", metrics.CPUUsage),
			Suggestion: "Monitor CPU usage patterns. Consider CPU upgrade if consistently high.",
			Color:      ColorYellow,
		})
	}

	// Check load average relative to CPU cores
	if metrics.LoadAverage[0] > float64(metrics.CPUCores)*1.5 {
		recommendations = append(recommendations, Recommendation{
			Component:  "CPU",
			Severity:   "HIGH",
			Reason:     fmt.Sprintf("Load average (%.2f) is high for %d cores", metrics.LoadAverage[0], metrics.CPUCores),
			Suggestion: "System is overloaded. Consider upgrading to more CPU cores or optimizing running processes.",
			Color:      ColorYellow,
		})
	}

	// Check for old CPU architectures (basic heuristic)
	if strings.Contains(strings.ToLower(metrics.CPUModel), "intel") && 
	   (strings.Contains(strings.ToLower(metrics.CPUModel), "core 2") || 
		strings.Contains(strings.ToLower(metrics.CPUModel), "core i3") ||
		strings.Contains(strings.ToLower(metrics.CPUModel), "core i5") && 
		!strings.Contains(strings.ToLower(metrics.CPUModel), "11th") &&
		!strings.Contains(strings.ToLower(metrics.CPUModel), "12th") &&
		!strings.Contains(strings.ToLower(metrics.CPUModel), "13th")) {
		recommendations = append(recommendations, Recommendation{
			Component:  "CPU",
			Severity:   "MEDIUM",
			Reason:     "CPU model appears to be older generation",
			Suggestion: "Consider upgrading to a newer CPU for better performance and efficiency.",
			Color:      ColorYellow,
		})
	}

	return recommendations
}

func analyzeMemory(metrics *SystemMetrics) []Recommendation {
	var recommendations []Recommendation

	memUsagePercent := float64(metrics.MemoryUsed) / float64(metrics.MemoryTotal) * 100
	currentRAMGB := float64(metrics.MemoryTotal) / (1024 * 1024 * 1024)

	// Calculate recommended RAM based on usage patterns
	recommendedRAM := calculateRecommendedRAM(currentRAMGB, memUsagePercent, float64(metrics.SwapUsed)/(1024*1024*1024))

	// Check memory usage percentage
	usedRAMGB := (memUsagePercent / 100) * currentRAMGB
	if memUsagePercent > 95 {
		// For critical usage, be more conservative
		conservativeRAM := calculateConservativeRAM(currentRAMGB, memUsagePercent, float64(metrics.SwapUsed)/(1024*1024*1024))
		recommendations = append(recommendations, Recommendation{
			Component:  "Memory",
			Severity:   "CRITICAL",
			Reason:     fmt.Sprintf("Memory usage is critical (%.1fGB/%.1fGB = %.1f%% used)", usedRAMGB, currentRAMGB, memUsagePercent),
			Suggestion: fmt.Sprintf("Urgently need more RAM. Minimum upgrade: %.0fGB (gives you %.1fGB headroom). Close applications immediately.", conservativeRAM, conservativeRAM-usedRAMGB),
			Color:      ColorRed,
		})
	} else if memUsagePercent > 85 {
		recommendations = append(recommendations, Recommendation{
			Component:  "Memory",
			Severity:   "HIGH",
			Reason:     fmt.Sprintf("Memory usage is high (%.1fGB/%.1fGB = %.1f%% used)", usedRAMGB, currentRAMGB, memUsagePercent),
			Suggestion: fmt.Sprintf("Consider upgrading to %.0fGB RAM to prevent slowdowns (provides %.1fGB buffer).", recommendedRAM, recommendedRAM-usedRAMGB),
			Color:      ColorYellow,
		})
	} else if memUsagePercent > 70 {
		recommendations = append(recommendations, Recommendation{
			Component:  "Memory",
			Severity:   "MEDIUM",
			Reason:     fmt.Sprintf("Memory usage is moderate (%.1fGB/%.1fGB = %.1f%% used)", usedRAMGB, currentRAMGB, memUsagePercent),
			Suggestion: fmt.Sprintf("Monitor memory usage. Consider %.0fGB for intensive tasks.", recommendedRAM),
			Color:      ColorYellow,
		})
	}

	// Check swap usage
	if metrics.SwapUsed > 0 {
		swapUsageGB := float64(metrics.SwapUsed) / (1024 * 1024 * 1024)
		usedRAMGB := (memUsagePercent / 100) * currentRAMGB
		totalMemoryNeed := usedRAMGB + swapUsageGB
		
		// Calculate both conservative and optimal recommendations
		conservativeRAM := calculateConservativeRAM(currentRAMGB, memUsagePercent, swapUsageGB)
		optimalRAM := calculateRecommendedRAM(currentRAMGB, memUsagePercent, swapUsageGB)
		
		if swapUsageGB > 2 {
			recommendations = append(recommendations, Recommendation{
				Component:  "Memory",
				Severity:   "HIGH",
				Reason:     fmt.Sprintf("Heavy swap usage (%.1fGB) - system is using disk as memory", swapUsageGB),
				Suggestion: fmt.Sprintf("Add more RAM immediately. Memory needed: %.1fGB (%.1fGB used + %.1fGB swap). Minimum: %.0fGB, Optimal: %.0fGB for headroom.", totalMemoryNeed, usedRAMGB, swapUsageGB, conservativeRAM, optimalRAM),
				Color:      ColorRed,
			})
		} else if swapUsageGB > 0.5 {
			recommendations = append(recommendations, Recommendation{
				Component:  "Memory",
				Severity:   "MEDIUM",
				Reason:     fmt.Sprintf("Moderate swap usage (%.1fGB)", swapUsageGB),
				Suggestion: fmt.Sprintf("Consider upgrading to %.0fGB RAM to eliminate swap (total need: %.1fGB with buffer).", conservativeRAM, totalMemoryNeed*1.15),
				Color:      ColorYellow,
			})
		}
	}

	// Check memory pressure
	if metrics.MemPressure == "critical" || metrics.MemPressure == "urgent" {
		recommendations = append(recommendations, Recommendation{
			Component:  "Memory",
			Severity:   "CRITICAL",
			Reason:     fmt.Sprintf("Memory pressure is %s", metrics.MemPressure),
			Suggestion: "System is under severe memory pressure. Upgrade RAM immediately.",
			Color:      ColorRed,
		})
	} else if metrics.MemPressure == "warning" {
		recommendations = append(recommendations, Recommendation{
			Component:  "Memory",
			Severity:   "HIGH",
			Reason:     "Memory pressure warning detected",
			Suggestion: "Consider upgrading RAM to prevent performance issues.",
			Color:      ColorYellow,
		})
	}

	// Check total memory amount
	if currentRAMGB < 8 {
		recommendations = append(recommendations, Recommendation{
			Component:  "Memory",
			Severity:   "HIGH",
			Reason:     fmt.Sprintf("Total RAM (%.1fGB) is below modern standards", currentRAMGB),
			Suggestion: fmt.Sprintf("Upgrade to at least 16GB RAM for modern applications (current: %.1fGB → recommended: 16GB+).", currentRAMGB),
			Color:      ColorYellow,
		})
	} else if currentRAMGB < 16 {
		recommendations = append(recommendations, Recommendation{
			Component:  "Memory",
			Severity:   "MEDIUM",
			Reason:     fmt.Sprintf("Total RAM (%.1fGB) may be limiting for intensive tasks", currentRAMGB),
			Suggestion: fmt.Sprintf("Consider upgrading to 32GB RAM for development/content creation (current: %.1fGB → recommended: 32GB).", currentRAMGB),
			Color:      ColorYellow,
		})
	} else if currentRAMGB >= 16 && currentRAMGB < 32 && (memUsagePercent > 80 || metrics.SwapUsed > 0) {
		// For systems with 16-32GB that are still running out of memory
		recommendations = append(recommendations, Recommendation{
			Component:  "Memory",
			Severity:   "MEDIUM",
			Reason:     fmt.Sprintf("Despite having %.1fGB RAM, still experiencing memory pressure", currentRAMGB),
			Suggestion: fmt.Sprintf("Upgrade to 64GB RAM for heavy workloads (current: %.1fGB → recommended: 64GB).", currentRAMGB),
			Color:      ColorYellow,
		})
	}

	return recommendations
}

// calculateRecommendedRAM suggests optimal RAM based on current usage patterns
func calculateRecommendedRAM(currentRAMGB, memUsagePercent, swapUsageGB float64) float64 {
	// More conservative approach: actual memory need + reasonable buffer
	usedRAMGB := (memUsagePercent / 100) * currentRAMGB
	totalMemoryNeed := usedRAMGB + swapUsageGB
	
	// Add a reasonable buffer (25%) but don't be overly generous
	targetRAM := totalMemoryNeed * 1.25
	
	// Round up to common RAM sizes, but prefer the next logical step
	commonSizes := []float64{8, 16, 24, 32, 48, 64, 96, 128, 192, 256}
	
	for _, size := range commonSizes {
		if targetRAM <= size {
			return size
		}
	}
	
	// If we need more than 256GB, round up to nearest 32GB
	return math.Ceil(targetRAM/32) * 32
}

// calculateConservativeRAM provides a more conservative estimate
func calculateConservativeRAM(currentRAMGB, memUsagePercent, swapUsageGB float64) float64 {
	// Simple calculation: current usage + swap + 15% buffer
	usedRAMGB := (memUsagePercent / 100) * currentRAMGB
	totalNeed := usedRAMGB + swapUsageGB
	conservativeTarget := totalNeed * 1.15
	
	// Round to next common size
	commonSizes := []float64{16, 24, 32, 48, 64, 96, 128}
	for _, size := range commonSizes {
		if conservativeTarget <= size {
			return size
		}
	}
	return 128
}

// max returns the maximum of two float64 values
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func analyzeGPU(metrics *SystemMetrics) []Recommendation {
	var recommendations []Recommendation

	// Basic GPU analysis based on model name
	if metrics.GPUModel != "" && metrics.GPUModel != "Unknown GPU" {
		gpuLower := strings.ToLower(metrics.GPUModel)
		
		// Check for integrated vs dedicated GPU
		if strings.Contains(gpuLower, "intel") && (strings.Contains(gpuLower, "hd") || strings.Contains(gpuLower, "iris")) {
			recommendations = append(recommendations, Recommendation{
				Component:  "GPU",
				Severity:   "MEDIUM",
				Reason:     "Using integrated Intel graphics",
				Suggestion: "For gaming or graphics-intensive work, consider a system with dedicated GPU.",
				Color:      ColorYellow,
			})
		}
		
		// Check for old AMD integrated graphics
		if strings.Contains(gpuLower, "radeon") && (strings.Contains(gpuLower, "r5") || strings.Contains(gpuLower, "r7")) {
			recommendations = append(recommendations, Recommendation{
				Component:  "GPU",
				Severity:   "MEDIUM",
				Reason:     "Using older integrated AMD graphics",
				Suggestion: "Consider upgrading to a system with newer integrated or dedicated graphics.",
				Color:      ColorYellow,
			})
		}

		// Note about Apple Silicon
		if strings.Contains(gpuLower, "apple") {
			recommendations = append(recommendations, Recommendation{
				Component:  "GPU",
				Severity:   "LOW",
				Reason:     "Using Apple Silicon integrated GPU",
				Suggestion: "Apple Silicon GPUs are generally excellent. Consider Mac Studio/Pro for intensive GPU work.",
				Color:      ColorGreen,
			})
		}
	}

	return recommendations
}

func displayRecommendations(recommendations []Recommendation) {
	if len(recommendations) == 0 {
		fmt.Printf("%s✅ Great! No bottlenecks detected%s\n", ColorGreen, ColorReset)
		fmt.Printf("Your system appears to be running optimally.\n")
		return
	}

	fmt.Printf("%s🔧 Upgrade Recommendations%s\n", ColorBold, ColorReset)
	fmt.Printf("═══════════════════════════\n\n")

	// Group by severity
	critical := []Recommendation{}
	high := []Recommendation{}
	medium := []Recommendation{}
	low := []Recommendation{}

	for _, rec := range recommendations {
		switch rec.Severity {
		case "CRITICAL":
			critical = append(critical, rec)
		case "HIGH":
			high = append(high, rec)
		case "MEDIUM":
			medium = append(medium, rec)
		case "LOW":
			low = append(low, rec)
		}
	}

	// Display by priority
	displayRecommendationGroup("🚨 CRITICAL", critical, ColorRed)
	displayRecommendationGroup("⚠️  HIGH", high, ColorYellow)
	displayRecommendationGroup("📋 MEDIUM", medium, ColorYellow)
	displayRecommendationGroup("💡 LOW", low, ColorGreen)

	fmt.Printf("\n%s💡 Pro Tips:%s\n", ColorBold, ColorReset)
	fmt.Printf("• Run this tool regularly to monitor system performance\n")
	fmt.Printf("• Close unnecessary applications before intensive tasks\n")
	fmt.Printf("• Consider upgrading components in order of severity\n")
	fmt.Printf("• Monitor Activity Monitor for specific resource-heavy processes\n")
}

func displayRecommendationGroup(title string, recommendations []Recommendation, color string) {
	if len(recommendations) == 0 {
		return
	}

	fmt.Printf("%s%s%s\n", color, title, ColorReset)
	fmt.Printf("────────────────\n")

	for _, rec := range recommendations {
		fmt.Printf("%s• %s (%s)%s\n", rec.Color, rec.Component, rec.Reason, ColorReset)
		fmt.Printf("  → %s\n", rec.Suggestion)
		fmt.Println()
	}
}