package utils

import "os"

func GetEnvDefault(key string, def string) (result string) {
	result = os.Getenv(key)
	if len(result) == 0 {
		result = def
	}
	return
}
