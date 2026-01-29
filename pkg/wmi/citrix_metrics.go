package wmi

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// Helper functions for type conversion
func convertToUint64(val *ole.VARIANT) uint64 {
	if val == nil {
		return 0
	}

	switch v := val.Value().(type) {
	case uint64:
		return v
	case int64:
		if v >= 0 {
			return uint64(v)
		}
		return 0
	case uint32:
		return uint64(v)
	case int32:
		if v >= 0 {
			return uint64(v)
		}
		return 0
	case uint16:
		return uint64(v)
	case int16:
		if v >= 0 {
			return uint64(v)
		}
		return 0
	case uint8:
		return uint64(v)
	case int8:
		if v >= 0 {
			return uint64(v)
		}
		return 0
	case float64:
		if v >= 0 {
			return uint64(v)
		}
		return 0
	case float32:
		if v >= 0 {
			return uint64(v)
		}
		return 0
	case string:
		if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
			return parsed
		}
		return 0
	default:
		return 0
	}
}

func convertToInt64(val *ole.VARIANT) int64 {
	if val == nil {
		return 0
	}

	switch v := val.Value().(type) {
	case int64:
		return v
	case uint64:
		if v <= 9223372036854775807 { // Max int64
			return int64(v)
		}
		return 0
	case int32:
		return int64(v)
	case uint32:
		return int64(v)
	case int16:
		return int64(v)
	case uint16:
		return int64(v)
	case int8:
		return int64(v)
	case uint8:
		return int64(v)
	case float64:
		return int64(v)
	case float32:
		return int64(v)
	case string:
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			return parsed
		}
		return 0
	default:
		return 0
	}
}

func convertToFloat64(val *ole.VARIANT) float64 {
	if val == nil {
		return 0
	}

	switch v := val.Value().(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int64:
		return float64(v)
	case uint64:
		return float64(v)
	case int32:
		return float64(v)
	case uint32:
		return float64(v)
	case int16:
		return float64(v)
	case uint16:
		return float64(v)
	case int8:
		return float64(v)
	case uint8:
		return float64(v)
	case string:
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			return parsed
		}
		return 0
	default:
		return 0
	}
}

// CitrixEndpointMetric represents the WMI data structure for Citrix Endpoint Metrics
type CitrixEndpointMetric struct {
	AvgBeaconLatency       float64
	AvgPrivilegedTime      float64
	AvgProcessorTime       float64
	AvgThroughputBytesRcvd float64
	AvgThroughputBytesSent float64
	AvgUserTime            float64
	City                   string
	ClientTimestamp        uint64
	Country                string
	EndpointIP             string
	GpuAvgUsage            float64
	GpuMaxUsage            uint64
	ISP                    string
	LatencyUnit            string
	LinkSpeed              uint64
	MaxPrivilegedTime      uint64
	MaxProcessorTime       uint64
	MaxThroughputBytesRcvd uint64
	MaxThroughputBytesSent uint64
	MaxUserTime            uint64
	NetworkInterfaceType   string
	RamAvgUsage            float64
	RamMaxUsage            uint64
	SessionID              int64
	SignalStrength         uint64
	SpeedUnit              string
	Timestamp              time.Time
	// Added mapped field
	SessionGUID string
}

// GetCitrixMetrics retrieves all Citrix Endpoint Metrics using WMI
func GetCitrixMetrics() ([]CitrixEndpointMetric, error) {
	var metrics []CitrixEndpointMetric

	// Initialize OLE
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OLE: %v", err)
	}
	defer ole.CoUninitialize()

	// Create WMI service
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil, fmt.Errorf("failed to create WMI locator object: %v", err)
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, fmt.Errorf("failed to get WMI interface: %v", err)
	}
	defer wmi.Release() // Connect to WMI service
	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer", nil, "ROOT\\Citrix\\EUEM")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WMI service: %v", err)
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// Execute WMI query
	resultsRaw, err := oleutil.CallMethod(service, "ExecQuery", "SELECT * FROM Citrix_Euem_EndpointMetrics")
	if err != nil {
		return nil, fmt.Errorf("failed to execute WMI query: %v", err)
	}
	results := resultsRaw.ToIDispatch()
	defer results.Release()

	// Process results
	countVar, err := oleutil.GetProperty(results, "Count")
	if err != nil {
		return nil, fmt.Errorf("failed to get result count: %v", err)
	}
	count := int(countVar.Val)

	for i := 0; i < count; i++ {
		itemRaw, err := oleutil.CallMethod(results, "ItemIndex", i)
		if err != nil {
			return nil, fmt.Errorf("failed to get result item %d: %v", i, err)
		}
		item := itemRaw.ToIDispatch()
		defer item.Release()

		metric := CitrixEndpointMetric{}

		// Extract all properties
		if val, err := oleutil.GetProperty(item, "AvgBeaconLatency"); err == nil {
			metric.AvgBeaconLatency = convertToFloat64(val)
		}
		if val, err := oleutil.GetProperty(item, "AvgPrivilegedTime"); err == nil {
			metric.AvgPrivilegedTime = convertToFloat64(val)
		}
		if val, err := oleutil.GetProperty(item, "AvgProcessorTime"); err == nil {
			metric.AvgProcessorTime = convertToFloat64(val)
		}
		if val, err := oleutil.GetProperty(item, "AvgThroughputBytesRcvd"); err == nil {
			metric.AvgThroughputBytesRcvd = convertToFloat64(val)
		}
		if val, err := oleutil.GetProperty(item, "AvgThroughputBytesSent"); err == nil {
			metric.AvgThroughputBytesSent = convertToFloat64(val)
		}
		if val, err := oleutil.GetProperty(item, "AvgUserTime"); err == nil {
			metric.AvgUserTime = convertToFloat64(val)
		}
		if val, err := oleutil.GetProperty(item, "City"); err == nil {
			metric.City = val.ToString()
		}
		if val, err := oleutil.GetProperty(item, "ClientTimestamp"); err == nil {
			metric.ClientTimestamp = convertToUint64(val)
		}
		if val, err := oleutil.GetProperty(item, "Country"); err == nil {
			metric.Country = val.ToString()
		}
		if val, err := oleutil.GetProperty(item, "EndpointIP"); err == nil {
			metric.EndpointIP = val.ToString()
		}
		if val, err := oleutil.GetProperty(item, "GpuAvgUsage"); err == nil {
			metric.GpuAvgUsage = convertToFloat64(val)
		}
		if val, err := oleutil.GetProperty(item, "GpuMaxUsage"); err == nil {
			metric.GpuMaxUsage = convertToUint64(val)
		}
		if val, err := oleutil.GetProperty(item, "ISP"); err == nil {
			metric.ISP = val.ToString()
		}
		if val, err := oleutil.GetProperty(item, "LatencyUnit"); err == nil {
			metric.LatencyUnit = val.ToString()
		}
		if val, err := oleutil.GetProperty(item, "LinkSpeed"); err == nil {
			metric.LinkSpeed = convertToUint64(val)
		}
		if val, err := oleutil.GetProperty(item, "MaxPrivilegedTime"); err == nil {
			metric.MaxPrivilegedTime = convertToUint64(val)
		}
		if val, err := oleutil.GetProperty(item, "MaxProcessorTime"); err == nil {
			metric.MaxProcessorTime = convertToUint64(val)
		}
		if val, err := oleutil.GetProperty(item, "MaxThroughputBytesRcvd"); err == nil {
			metric.MaxThroughputBytesRcvd = convertToUint64(val)
		}
		if val, err := oleutil.GetProperty(item, "MaxThroughputBytesSent"); err == nil {
			metric.MaxThroughputBytesSent = convertToUint64(val)
		}
		if val, err := oleutil.GetProperty(item, "MaxUserTime"); err == nil {
			metric.MaxUserTime = convertToUint64(val)
		}
		if val, err := oleutil.GetProperty(item, "NetworkInterfaceType"); err == nil {
			metric.NetworkInterfaceType = val.ToString()
		}
		if val, err := oleutil.GetProperty(item, "RamAvgUsage"); err == nil {
			metric.RamAvgUsage = convertToFloat64(val)
		}
		if val, err := oleutil.GetProperty(item, "RamMaxUsage"); err == nil {
			metric.RamMaxUsage = convertToUint64(val)
		}
		if val, err := oleutil.GetProperty(item, "SessionID"); err == nil {
			metric.SessionID = convertToInt64(val)
		}
		if val, err := oleutil.GetProperty(item, "SignalStrength"); err == nil {
			metric.SignalStrength = convertToUint64(val)
		}
		if val, err := oleutil.GetProperty(item, "SpeedUnit"); err == nil {
			metric.SpeedUnit = val.ToString()
		}
		if val, err := oleutil.GetProperty(item, "Timestamp"); err == nil {
			// Convert OLE date to Go time.Time
			if date, ok := val.Value().(time.Time); ok {
				metric.Timestamp = date
			}
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// PrintMetricKeyValue formats and prints a CitrixEndpointMetric in key=value format (same as PowerShell script)
// Restored full output with selected metrics (previously excluded: AvgPrivilegedTime, AvgProcessorTime, AvgUserTime, ISP, RamAvgUsage, SessionID, SpeedUnit, Timestamp)
func PrintMetricKeyValue(metric CitrixEndpointMetric) {
	output := fmt.Sprintf("AvgBeaconLatency=%v AvgThroughputBytesRcvd=%v AvgThroughputBytesSent=%v ClientTimestamp=%v GpuAvgUsage=%v GpuMaxUsage=%v LinkSpeed=%v MaxPrivilegedTime=%v MaxProcessorTime=%v MaxThroughputBytesRcvd=%v MaxThroughputBytesSent=%v MaxUserTime=%v NetworkInterfaceType=%s RamMaxUsage=%v SessionGUID=%s SignalStrength=%v",
		metric.AvgBeaconLatency,
		// metric.AvgPrivilegedTime,     // EXCLUDED - uncomment to re-enable
		// metric.AvgProcessorTime,      // EXCLUDED - uncomment to re-enable
		metric.AvgThroughputBytesRcvd,
		metric.AvgThroughputBytesSent,
		// metric.AvgUserTime,           // EXCLUDED - uncomment to re-enable
		//metric.City,
		metric.ClientTimestamp,
		//metric.Country,
		//metric.EndpointIP,
		metric.GpuAvgUsage,
		metric.GpuMaxUsage,
		// metric.ISP,                   // EXCLUDED - uncomment to re-enable
		//metric.LatencyUnit,
		metric.LinkSpeed,
		metric.MaxPrivilegedTime,
		metric.MaxProcessorTime,
		metric.MaxThroughputBytesRcvd,
		metric.MaxThroughputBytesSent,
		metric.MaxUserTime,
		metric.NetworkInterfaceType,
		// metric.RamAvgUsage,           // EXCLUDED - uncomment to re-enable
		metric.RamMaxUsage,
		// metric.SessionID,             // EXCLUDED - uncomment to re-enable
		metric.SessionGUID,
		metric.SignalStrength,
		// metric.SpeedUnit,             // EXCLUDED - uncomment to re-enable
		// metric.Timestamp,             // EXCLUDED - uncomment to re-enable
	)
	fmt.Println(output)
}
