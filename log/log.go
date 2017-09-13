package log

import (
	"fmt"
	"encoding/json"
)

func PrintJSON(data interface{}) {
	if c, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Println(string(c))
	}
}
