package report

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/filz0r/go-crawler/internal/parser"
)

func WriteJSONReport(pages map[string]parser.PageData, filename string) error {
	sortedKeys := make([]string, 0, len(pages))
	for key := range pages {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	result := make([]parser.PageData, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		result = append(result, pages[key])
	}
	indented, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, indented, 0644)
	if err != nil {
		return err
	}
	return nil
}
