# Citrix Endpoint Metrics Collector

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

This Go project collects metrics from Citrix endpoint sessions using WMI queries and maps the SessionIDs to their corresponding SessionGUIDs from the Windows registry. **This tool is designed to work in conjunction with uberAgent** to provide comprehensive endpoint performance metrics for Citrix sessions.

## Overview

This application retrieves endpoint performance metrics that are available when using Citrix Virtual Apps and Desktops. The metrics provide valuable insights into the end-user experience including network latency, throughput, GPU usage, and system performance from the endpoint perspective.

**Important**: This tool is specifically designed to complement uberAgent deployments by providing additional Citrix endpoint metrics that can be correlated with uberAgent's comprehensive endpoint monitoring data.

## Citrix Requirements

To have these metrics available, you need to meet specific Citrix infrastructure requirements:

- **Citrix Virtual Apps and Desktops** with appropriate VDA versions
- **Citrix Workspace App** with specific minimum versions  

For detailed requirements regarding Workspace App versions and VDA compatibility, please refer to the official Citrix documentation:
**[Session Performance Metrics Requirements](https://docs.citrix.com/en-us/citrix-virtual-apps-desktops/director/troubleshoot-deployments/user-issues/session-performance.html#session-performance-metrics)**

## Features

- **WMI Integration**: Retrieves Citrix endpoint metrics from WMI namespace `ROOT\Citrix\EUEM\Citrix_Euem_EndpointMetrics`
- **SessionGUID Mapping**: Maps Citrix SessionIDs to uberAgent SessionGUIDs using registry values from `HKEY_LOCAL_MACHINE\SOFTWARE\vast limits\uberAgent\SessionGuids`
- **uberAgent Compatibility**: Designed to work alongside uberAgent for comprehensive endpoint monitoring
- **Output**: Outputs metrics in key=value format 

## SessionGUID Mapping

The application performs automatic mapping between Citrix SessionIDs and uberAgent SessionGUIDs:

1. **Citrix SessionID**: Retrieved from WMI Citrix endpoint metrics
2. **uberAgent SessionGUID**: Read from Windows registry at `HKEY_LOCAL_MACHINE\SOFTWARE\vast limits\uberAgent\SessionGuids`
3. **Correlation**: The SessionID is used as a key to lookup the corresponding SessionGUID from the uberAgent registry data
4. **Output**: Both SessionID and mapped SessionGUID are included in the output for correlation with uberAgent data

This mapping enables correlation between Citrix endpoint performance data and uberAgent's comprehensive endpoint monitoring metrics.

## Collected Metrics

The application collects the following endpoint metrics:

- **Network Performance**: AvgBeaconLatency, AvgThroughputBytesRcvd, AvgThroughputBytesSent, LinkSpeed, SignalStrength
- **System Performance**: GpuAvgUsage, GpuMaxUsage, RamMaxUsage, MaxProcessorTime, MaxPrivilegedTime, MaxUserTime
- **Endpoint Information**: City, Country, EndpointIP, NetworkInterfaceType
- **Session Data**: ClientTimestamp, SessionGUID (mapped from uberAgent registry)

## Requirements

- **Operating System**: Windows (where Citrix session is running)
- **Go Version**: Go 1.16 or higher
- **Privileges**: Administrative privileges (for WMI and registry access)
- **Citrix Environment**: Properly configured Citrix Virtual Apps/Desktops with supported VDA and Workspace App versions
- **uberAgent**: uberAgent installed and configured (for SessionGUID mapping)

## Dependencies

- `github.com/go-ole/go-ole` - For WMI interaction with Citrix EUEM namespace
- `golang.org/x/sys/windows/registry` - For Windows registry access to map SessionGUIDs

## Installation

1. Clone or download the project files to your target directory:
```bash
git clone https://github.com/andreaz86/uberAgentEndpointMetrics
cd uberAgentEndpointMetrics
```

2. To compile the project into a Windows executable:

```bash
# Build for current architecture
go build -o uaEndpointMetrics.exe main.go

# Build for specific architecture (if needed)
# For 64-bit Windows
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o uaEndpointMetrics.exe main.go

```


## Usage

To integrate this tool with uberAgent for automated collection, add the following configuration to your uberAgent configuration file:

```ini
[Timer]
Name              = EndpointMetrics
Interval          = 60000
Script            = "#exePath#" (example "C:\Scripts\uaEndpointMetrics.exe")
ScriptContext     = Session0AsSystem
```

**Configuration Parameters:**
- **Name**: Identifier for the timer in uberAgent
- **Interval**: Collection interval in milliseconds (60000 = 1 minute)
- **Script**: Full path to the compiled executable
- **ScriptContext**: Execution context (Session0AsSystem for system-level access)

This configuration will automatically execute the Citrix endpoint metrics collector every minute and integrate the data with uberAgent's monitoring system.

## Output Format

The application outputs metrics in key=value format on a single line per metric record:

```
AvgBeaconLatency=197 AvgThroughputBytesRcvd=4096 AvgThroughputBytesSent=4096 City=Milan ClientTimestamp=1748452631 Country=Italy EndpointIP=93.45.42.162 GpuAvgUsage=12 GpuMaxUsage=15 LatencyUnit=ms LinkSpeed=648000000 MaxPrivilegedTime=7 MaxProcessorTime=18 MaxThroughputBytesRcvd=10240 MaxThroughputBytesSent=8192 MaxUserTime=10 NetworkInterfaceType=Wifi RamMaxUsage=53 SessionGUID=00000003-7165-2c71-538b-6ec2f3cfdb01 SignalStrength=61
```

This format is designed for easy parsing and integration with monitoring systems and correlation with uberAgent data.

## Project Structure

```
ua_endpoint_metrics/
├── main.go                    # Main application entry point
├── go.mod                     # Go module definition
├── build.bat                  # Build script for Windows
├── README.md                  # This file
└── pkg/
    ├── wmi/
    │   └── citrix_metrics.go   # WMI data extraction and formatting
    └── registry/
        └── session_guids.go    # Registry access for SessionGUID mapping
```


## License

[GNU General Public License v3.0](LICENSE)
