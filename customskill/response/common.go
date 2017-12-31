package response

import "encoding/json"

var jsonMarshal = json.Marshal // Used to enable unit testing

func Bool(bool bool) *bool {
	return &bool
}
