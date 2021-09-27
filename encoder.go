package kleos

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// RFC3339Ms is the time format with up to three decimal places of ms, used for parsing
	// times.  Used for ingesting times from JSON, SQL.
	RFC3339Ms = "2006-01-02T15:04:05.999Z07:00"
)

// Encode a field value for output.
func encode(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 5, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', 5, 64)
	case bool:
		return strconv.FormatBool(v)
	case string:
		if strings.Contains(v, " ") {
			return strconv.Quote(v)
		}
		return v
	case error:
		return strconv.Quote(v.Error())
	case uuid.UUID:
		if BlankUUID(v) {
			return ""
		}
		return v.String()
	case time.Time:
		if v.IsZero() {
			return ""
		}
		return v.Format(RFC3339Ms)
	case fmt.Stringer:
		res := v.String()
		if strings.Contains(res, " ") {
			res = strconv.Quote(res)
		}
		res = Cleanup(res)
		return res
	default:
		res := fmt.Sprintf("%v", value)
		if strings.Contains(res, " ") {
			res = strconv.Quote(res)
		}
		res = Cleanup(res)
		return res
	}
}

// BlankUUID checks to see if the uuid.UUID value is blank.
func BlankUUID(id uuid.UUID) bool {
	for _, b := range id {
		if b != 0 {
			return false
		}
	}

	return true
}

// Cleanup removes carriage return characters from the string.
func Cleanup(v string) string {
	v = strings.ReplaceAll(v, `\n`, "")
	v = strings.ReplaceAll(v, `\r`, "")
	return v
}
