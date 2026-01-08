package cluster

import (
	"encoding/json"
	"fmt"
	"time"
)

type InventoryOutput struct {
	Timestamp string   `json:"timestamp"`
	Servers   []Server `json:"servers"`
	Total     int      `json:"total"`
	Online    int      `json:"online"`
	Offline   int      `json:"offline"`
}

func InventoryToJSON(inv Inventory) ([]byte, error) {
	online := 0
	offline := 0
	for _, s := range inv.Servers {
		if s.Status == "online" {
			online++
		} else if s.Status == "offline" {
			offline++
		}
	}
	out := InventoryOutput{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Servers:   inv.Servers,
		Total:     len(inv.Servers),
		Online:    online,
		Offline:   offline,
	}
	return json.MarshalIndent(out, "", "  ")
}

func InventoryToText(inv Inventory) string {
	b := ""
	for _, s := range inv.Servers {
		id := s.Hostname
		if id == "" {
			id = s.IP
		}
		b += fmt.Sprintf("%s\t%s\t%s\n", id, s.IP, s.Status)
	}
	return b
}
