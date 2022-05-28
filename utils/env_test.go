package utils

import (
	"os"
	"testing"
)

func TestGetEnvDefault(t *testing.T) {
	t.Run("get key that isn't set", func(t *testing.T) {
		result := GetEnvDefault("MUSYNC_UTILS_TESTING", "foo")
		if result != "foo" {
			t.Fail()
		}
	})

	t.Run("get key after setting it", func(t *testing.T) {
		os.Setenv("MUSYNC_UTILS_TESTING", "bar")
		result := GetEnvDefault("MUSYNC_UTILS_TESTING", "foo")
		if result != "bar" {
			t.Fail()
		}
	})
}
