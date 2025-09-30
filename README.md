# System Bottleneck Checker

A comprehensive cross-platform system performance analyzer that identifies CPU, Memory, and GPU bottlenecks and provides specific upgrade recommendations.

**Supported Platforms:** macOS, Linux, Windows

## Features

- **Real-time Performance Analysis**: Monitors CPU usage, load averages, memory consumption, swap usage, and memory pressure
- **Smart Recommendations**: Provides prioritized upgrade suggestions with detailed explanations
- **Color-coded Output**: Easy-to-read results with severity indicators
- **Comprehensive Coverage**: Analyzes CPU, Memory (RAM), and GPU components
- **Cross-Platform**: Works on macOS, Linux, and Windows using gopsutil for system metrics

## Installation

1. **Clone or download this project**
   ```bash
   git clone https://github.com/xmarkclx/bottleneck-check.git
   cd bottleneck-check
   ```

2. **Build the executable**
   
   **macOS/Linux:**
   ```bash
   go build -o bottleneck-check .
   ```
   
   **Windows:**
   ```cmd
   go build -o bottleneck-check.exe .
   ```

3. **Run the analyzer**
   
   **macOS/Linux:**
   ```bash
   ./bottleneck-check
   ```
   
   **Windows:**
   ```cmd
   bottleneck-check.exe
   ```

## Usage

Simply run the executable to start continuous monitoring:

**macOS/Linux:**
```bash
./bottleneck-check
```

**Windows:**
```cmd
bottleneck-check.exe
```

The tool will:
1. Start continuous real-time monitoring (updates every 10 seconds)
2. Display live system performance metrics with detailed advice
3. Show specific upgrade recommendations automatically
4. Provide an interactive menu for additional features

### Interactive Commands
Once running, press **Enter** to access the menu with these options:
- **[a]** - Refresh advice and recommendations (shown by default)
- **[d]** - Show detailed system information and hardware specs  
- **[s]** - Force refresh system status
- **[c]** - Clear screen
- **[h]** - Show comprehensive help guide
- **[q]** - Quit monitor

**Note:** Detailed upgrade recommendations are now displayed automatically - no need to press 'a' unless you want to force a refresh!

## What It Checks

### CPU Analysis
- Current usage percentage
- Load averages (1, 5, 15 minutes)
- Core count vs load ratio
- CPU model and generation detection

### Memory Analysis
- RAM usage percentage
- Swap file usage
- Memory pressure status
- Total memory adequacy for modern workloads

### GPU Analysis
- Graphics card model detection
- Integrated vs dedicated GPU identification
- Performance recommendations based on use case

## Recommendation Levels

- 🚨 **CRITICAL**: Immediate action required - system severely impacted
- ⚠️ **HIGH**: Should address soon - noticeable performance impact
- 📋 **MEDIUM**: Consider for future upgrades
- 💡 **LOW**: Optional improvements or informational

## Example Output

### Live Monitoring Display
```
🔍 System Bottleneck Monitor
═══════════════════════════════════
Last updated: 14:25:30 | Press [Enter] for menu

📊 Quick Status
─────────────
CPU: 85.2% | Load: 20.18 | Memory: 99.7% | Swap: 11.6GB

🚨 Active Alerts
──────────────
● CRITICAL: 1 issue(s) need immediate attention
● HIGH: 2 issue(s) affecting performance

──────────────────────────────────────────────────
Monitoring active... Press [Enter] for options
```

### Detailed Recommendations (press 'a')
```
🔧 Detailed Recommendations
════════════════════════════════

📊 Current System Status
─────────────────────────
CPU: Apple M3 Pro (12 cores)
  Usage: 85.2%
  Load Average: 20.18, 18.28, 14.83
Memory: 17.2GB used / 17.3GB total (99.7%)
  Swap: 11.6GB used / 13.0GB total
  Memory Pressure: critical
GPU: Apple M3 Pro

🔧 Upgrade Recommendations
═══════════════════════════

🚨 CRITICAL
────────────────
• Memory (Memory usage is critical (99.7%))
  → Urgently need more RAM. Close applications or upgrade memory.

⚠️  HIGH
────────────────
• Memory (Heavy swap usage (11.6GB) - system is using disk as memory)
  → Add more RAM immediately. Swap usage causes significant slowdowns.
```

## Requirements

### Supported Operating Systems
- **macOS**: 10.15+ (tested on macOS 13+)
- **Linux**: Most distributions with kernel 3.0+
- **Windows**: Windows 7/Server 2008R2 and later

### Build Requirements
- Go 1.21+ (for building from source)
- Internet connection (for downloading dependencies)

### Runtime
- No admin privileges required for basic metrics
- Some advanced features may require elevated permissions on certain platforms

## Tips for Better Performance

- Run regularly to monitor system trends
- Close unnecessary applications before intensive tasks
- Address recommendations in order of severity
- Use system monitoring tools to identify specific resource-heavy processes:
  - **macOS**: Activity Monitor
  - **Windows**: Task Manager or Resource Monitor
  - **Linux**: htop, top, or system monitor GUI
- Consider the tool's suggestions in context of your specific workload

## Technical Details

The tool uses the **[gopsutil](https://github.com/shirou/gopsutil)** library for cross-platform system monitoring:

### Cross-Platform Data Sources
- **CPU Information**: `/proc/cpuinfo` (Linux), `sysctl` (macOS), WMI (Windows)
- **Memory Statistics**: `/proc/meminfo` (Linux), `vm_stat` (macOS), Performance Counters (Windows)
- **CPU Usage**: `/proc/stat` (Linux), `iostat` (macOS), Performance Counters (Windows)
- **Load Averages**: `/proc/loadavg` (Linux), `uptime` (macOS), CPU percentage estimation (Windows)
- **System Information**: Various platform-specific APIs

### Dependencies
- **[github.com/shirou/gopsutil/v3](https://github.com/shirou/gopsutil)**: Cross-platform system and process monitoring library
- **Standard Go libraries**: For core functionality and UI

All metrics are collected using platform-appropriate APIs through gopsutil to ensure accuracy and reliability across different operating systems.
