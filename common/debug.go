package common

import (
	"encoding/json"
	"fmt"
)

func DumpJSON(value any) {
	formattedBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		// TODO: why don't these trigger linting errors?
		fmt.Printf("DEBUG: warning: couldn't dump JSON")
		return
	}

	fmt.Printf("DEBUG: dumped JSON: %v\n", string(formattedBytes))
}
