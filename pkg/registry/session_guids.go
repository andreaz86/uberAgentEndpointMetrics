package registry

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// SessionGUIDMapping maps SessionIDs to their corresponding SessionGUIDs
type SessionGUIDMapping map[string]string

// GetSessionGUIDs retrieves the mapping between SessionIDs and SessionGUIDs from the registry
func GetSessionGUIDs() (SessionGUIDMapping, error) {
	mapping := make(SessionGUIDMapping)

	// Open the registry key
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\vast limits\uberAgent\SessionGuids`, registry.READ)
	if err != nil {
		return nil, fmt.Errorf("failed to open registry key: %v", err)
	}
	defer k.Close()

	// Read all value names and data
	valueNames, err := k.ReadValueNames(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read registry value names: %v", err)
	}

	// Process each value (key is SessionID, value is SessionGUID)
	for _, name := range valueNames {
		sessionGUID, _, err := k.GetStringValue(name)
		if err != nil {
			fmt.Printf("Warning: Could not read value for %s: %v\n", name, err)
			continue
		}

		// Add to mapping
		mapping[name] = sessionGUID
	}

	return mapping, nil
}

// MapSessionGUIDs associates SessionGUIDs with the corresponding metrics based on SessionID
func MapSessionGUIDs(metrics interface{}, sessionGUIDs SessionGUIDMapping) {
	// This function will be implemented in the main package
	// We'll cast the interface and update the SessionGUID field there
}
