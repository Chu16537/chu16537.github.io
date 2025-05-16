package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// Step 1: 讀取對照表 company.json
	companyMap := map[string]int{}
	companyFile, err := os.ReadFile("data/company.json")
	if err != nil {
		log.Fatalf("read company.json error: %v", err)
	}
	if err := json.Unmarshal(companyFile, &companyMap); err != nil {
		log.Fatalf("unmarshal company.json error: %v", err)
	}
	// Step 2: 讀取 csv
	csvFile, err := os.Open("data.csv")
	if err != nil {
		log.Fatalf("open data.csv error: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("read csv error: %v", err)
	}
	if len(records) < 2 {
		log.Fatalf("csv rows insufficient")
	}
	headers := records[0]

	fmt.Println("companyMap", companyMap)
	// Step 3: 將每一筆資料分類到對應 json 檔案
	for _, row := range records[1:] {
		if len(row) == 0 || row[0] == "" {
			continue
		}

		entry := map[string]string{}
		for i := 0; i < len(headers) && i < len(row); i++ {
			entry[headers[i]] = row[i]
		}

		companyName := entry["公司名稱"]
		id, ok := companyMap[companyName]
		if !ok {
			fmt.Printf("公司名稱找不到對應: %s\n", companyName)
			continue
		}

		fileName := fmt.Sprintf("data/%d.json", id)

		var data []map[string]string
		// 如果檔案存在，先讀進來
		if _, err := os.Stat(fileName); err == nil {
			content, err := os.ReadFile(fileName)
			if err == nil {
				json.Unmarshal(content, &data)
			}
		}

		// append 新資料並寫回檔案
		data = append(data, entry)
		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Fatalf("marshal error: %v", err)
		}

		if err := ioutil.WriteFile(fileName, out, 0644); err != nil {
			log.Fatalf("write file %s error: %v", fileName, err)
		}
	}

	fmt.Println("資料已分類寫入完成。")
}
