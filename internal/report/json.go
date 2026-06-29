package report

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/boserwuge/k8s-cluster-inspector/internal/model"
)

func WriteJSON(r model.Report, filename string) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Println("json marshal failed:", err)
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Println("write json failed:", err)
		return
	}

	fmt.Println("Generated", filename)
}
