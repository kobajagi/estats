package output

import (
	"encoding/json"
	"fmt"
)

// JsonWriter satisfies OutputWriter interface.
type JsonWriter struct{}

// Write prints Report to stdout as JSON.
// JSON document is minimized.
func (w JsonWriter) Write(data Report) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
