package common

import "fmt"

func MakeEventID(id string) string {
	return fmt.Sprintf("by-id.%s", id)
}
