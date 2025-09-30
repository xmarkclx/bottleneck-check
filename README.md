# System Bottleneck Checker

A comprehensive macOS system performance analyzer that identifies CPU, Memory, and GPU bottlenecks and provides specific upgrade recommendations.

## Features

- **Real-time Performance Analysis**: Monitors CPU usage, load averages, memory consumption, swap usage, and memory pressure
- **Smart Recommendations**: Provides prioritized upgrade suggestions with detailed explanations
- **Color-coded Output**: Easy-to-read results with severity indicators
- **Comprehensive Coverage**: Analyzes CPU, Memory (RAM), and GPU components
- **macOS Optimized**: Uses native macOS system tools and APIs for accurate readings

## Installation

1. Clone or download this project
2. Build the executable:
   ```bash
   cd bottleneck-check
   go build -o bottleneck-check .
   ```
3. Run the analyzer:
   ```bash
   ./bottleneck-check
   ```

## Usage

Simply run the executable to start continuous monitoring:
```bash
./bottleneck-check
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

- macOS (tested on macOS 13+)
- Go 1.21+ (for building from source)
- Admin privileges may be required for some system metrics

## Tips for Better Performance

- Run regularly to monitor system trends
- Close unnecessary applications before intensive tasks
- Address recommendations in order of severity
- Use Activity Monitor to identify specific resource-heavy processes
- Consider the tool's suggestions in context of your specific workload

## Technical Details

The tool uses various macOS system utilities:
- `sysctl` for CPU and system information
- `vm_stat` for memory statistics
- `iostat` for CPU usage sampling
- `uptime` for load averages
- `system_profiler` for hardware details
- `memory_pressure` for memory pressure status

All metrics are collected using native macOS APIs to ensure accuracy and reliability.