package blocks

import "encoding/json"

// mustUnmarshal unmarshals JSON data into v, panicking on error.
// Use in tests where unmarshal failure indicates a bug.
func mustUnmarshal(data []byte, v any) {
	if err := json.Unmarshal(data, v); err != nil {
		panic("mustUnmarshal: " + err.Error())
	}
}
