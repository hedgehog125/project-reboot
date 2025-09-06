package common

import (
	"encoding/json"
	"fmt"
)

func DumpJSON(value any) {
	formattedBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		//nolint
		fmt.Printf("DEBUG: warning: couldn't dump JSON")
		return
	}

	//nolint
	fmt.Printf("DEBUG: dumped JSON: %v\n", string(formattedBytes))
}
