package data

import (
	"encoding/json"
	"fmt"
	"os"
)

func RoadTypes() {
	file, err := os.Open(dataFile)
	if err != nil {
		return
	}
	defer file.Close()

	var data Data
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return
	}

	roadTypes := make(map[string]int)
	for _, el := range data.Elements {
		if el.Type == "way" {
			if road, ok := el.Tags["highway"]; ok {
				roadTypes[road]++
			}
		}

	}
	for rt, count := range roadTypes {
		fmt.Printf("Highway type: %s, Count: %d\n", rt, count)
	}
}
