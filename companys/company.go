package companys

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	inputDir     = "companys"
	outputDir    = "data"
	groupCount   = 20
	outputPrefix = "company"
)

func Run() {
	groupedData := make(map[int]map[string]string)

	// 初始化每組 map
	for i := 0; i <= groupCount; i++ {
		groupedData[i] = make(map[string]string)
	}

	// 掃描所有 CSV 檔案
	err := filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".csv") {
			fmt.Println("處理:", path)
			if err := parseCSV(path, groupedData); err != nil {
				fmt.Printf("錯誤處理 %s: %v\n", path, err)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// 確保輸出目錄存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}

	// 輸出每個 group 的 JSON 檔案
	for groupID, companyMap := range groupedData {
		if len(companyMap) == 0 {
			continue
		}
		filePath := fmt.Sprintf("%s/%s%d.json", outputDir, outputPrefix, groupID)
		jsonData, err := json.MarshalIndent(companyMap, "", "  ")
		if err != nil {
			fmt.Printf("JSON 轉換失敗 group %d: %v\n", groupID, err)
			continue
		}
		err = os.WriteFile(filePath, jsonData, 0644)
		if err != nil {
			fmt.Printf("寫入失敗: %s: %v\n", filePath, err)
			continue
		}
		fmt.Println("已輸出:", filePath)
	}
}

// 解析 CSV 並加入對應 group
func parseCSV(filePath string, groupedData map[int]map[string]string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	if len(records) < 2 {
		return nil
	}

	const (
		indexID   = 1
		indexName = 2
	)

	for _, row := range records[1:] {
		if len(row) <= indexName || len(row) <= indexID {
			continue
		}
		name := strings.Trim(row[indexName], "\"")
		id := strings.Trim(row[indexID], "\"")

		if name == "" || id == "" {
			continue
		}

		firstChar := string([]rune(name)[0])
		groupID := hashGroup(firstChar, groupCount)

		groupedData[groupID][name] = id
	}
	return nil
}

// hashGroup 將字元分配到 1~groupCount 的組別
func hashGroup(char string, groupCount int) int {
	runes := []rune(char)
	if len(runes) == 0 {
		return 0
	}
	code := int(runes[0])
	return (code % groupCount)
}
