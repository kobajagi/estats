package output

import (
	"encoding/json"
	"fmt"
)

type JsonWriter struct{}

func (w JsonWriter) Write(data Report) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
