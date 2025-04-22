package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const dataFile = "tokyo.json"

type Data struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	ID    int64             `json:"id"`
	Type  string            `json:"type"`
	Lat   float64           `json:"lat,omitempty"`
	Lon   float64           `json:"lon,omitempty"`
	Tags  map[string]string `json:"tags,omitempty"`
	Nodes []int64           `json:"nodes,omitempty"`
}

/*
seoul, yeo-ui-do
37.50550444160902,126.89277648925783,37.56988888346707,126.97208404541016

seoul, gangnam
37.47322325945579,126.99371337890626,37.54294402158324,127.10494995117189

jeju, seo-gui-po
33.23244172181488,126.5394973754883,33.26718283944606,126.58730506896974

tokyo, imperial palace
35.66594341573466,139.72034454345706,35.70428649806425,139.79175567626956
*/

func GetOSMData() Data {
	query := `
		[out:json];
		way["highway"](35.66594341573466,139.72034454345706,35.70428649806425,139.79175567626956);
		out body;
		>;
		out skel qt;
	`

	resp, err := http.PostForm("https://overpass-api.de/api/interpreter",
		url.Values{"data": {query}})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result Data
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}

	// Print data
	// for _, el := range result.Elements {
	// 	if el.Type == "way" && el.Tags != nil {
	// 		fmt.Printf("Way ID: %d, highway: %s, tags: %v\n", el.ID, el.Tags["highway"], el.Tags)
	// 	}
	// }
	return result
}

func StoreOSMData() {
	file, err := os.Create("data.json")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(GetOSMData()); err != nil {
		log.Fatal("Error encoding JSON to file:", err)
	}

	fmt.Println("data.json created")
}

func LoadOSMData() (*Data, error) {
	//seoguipo.json
	//yeouido.json
	//gangnam.json
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data Data
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
