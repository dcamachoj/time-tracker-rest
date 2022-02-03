package common

import "os"

// GetEnv get Environment Value with a default value if not exists
func GetEnv(key string, value string) string {
	var val, ok = os.LookupEnv(key)
	if !ok {
		return value
	}
	return val
}
