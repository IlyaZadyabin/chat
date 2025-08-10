package prettier

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	PlaceholderDollar = "$"
)

// Pretty formats SQL query, replacing placeholders with actual values for readability
func Pretty(query, placeholder string, args ...interface{}) string {
	for i, arg := range args {
		var value string
		switch v := arg.(type) {
		case string:
			value = fmt.Sprintf("'%s'", v)
		case []byte:
			value = fmt.Sprintf("'%s'", string(v))
		case int, int8, int16, int32, int64:
			value = fmt.Sprintf("%d", v)
		case uint, uint8, uint16, uint32, uint64:
			value = fmt.Sprintf("%d", v)
		case float32, float64:
			value = fmt.Sprintf("%f", v)
		case bool:
			value = strconv.FormatBool(v)
		case nil:
			value = "NULL"
		default:
			value = fmt.Sprintf("'%v'", v)
		}

		placeholder := fmt.Sprintf("%s%d", placeholder, i+1)
		query = strings.Replace(query, placeholder, value, 1)
	}

	return query
}
