package main

import (
	"fmt"
	"log"

	"github.com/ua_endpoint_metrics/pkg/registry"
	"github.com/ua_endpoint_metrics/pkg/wmi"
)

func main() {
	// Get Citrix metrics
	metrics, err := wmi.GetCitrixMetrics()
	if err != nil {
		log.Fatalf("Failed to get Citrix metrics: %v", err)
	}

	// Get session GUID mappings from registry
	sessionGUIDs, err := registry.GetSessionGUIDs()
	if err != nil {
		log.Printf("Warning: Failed to get session GUIDs from registry: %v", err)
	} // Associate session GUIDs with metrics
	for i := range metrics {
		sessionIDStr := fmt.Sprintf("%d", metrics[i].SessionID)
		if metrics[i].SessionID != 0 {
			if guid, exists := sessionGUIDs[sessionIDStr]; exists {
				metrics[i].SessionGUID = guid
			} else {
				metrics[i].SessionGUID = "No GUID found for SessionID"
			}
		} else {
			metrics[i].SessionGUID = "No SessionID available"
		}
	}
	// Print metrics in key=value format (one line per metric)
	for _, metric := range metrics {
		wmi.PrintMetricKeyValue(metric)
	}
}
