// Contains some json encoding utilities.

package dashboard

import (
	"bytes"
	"encoding/json"
	"regexp"
)

// jsonFieldStartRegexp is a regular expression to match the starting characters
// of field names in a json blob.
var jsonFieldStartRegexp = regexp.MustCompile("[,{]\"([A-Z])")

// jsonMarshalLowercase is a wrapper around json.Marshal, which after converting
// the given value to json will iterate over all te field names and force them to
// lowercase first character. This is needed to handle flattening objects that are
// not under our control.
func jsonMarshalLowercase(v interface{}) ([]byte, error) {
	// Flatten the object into a json string first
	blob, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	// Iterate over all field names and lowercase them
	for _, match := range jsonFieldStartRegexp.FindAllSubmatchIndex(blob, -1) {
		blob[match[2]] = bytes.ToLower(blob[match[2] : match[2]+1])[0]
	}
	return blob, nil
}
