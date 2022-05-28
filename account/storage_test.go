package account

import (
	"os"
	"strings"
	"testing"
)

type mockStorage struct {
	Value string
}

func TestStorage(t *testing.T) {
	storage := NewTokenStorage()

	t.Run("should have the needed paths", func(t *testing.T) {
		if len(storage.baseDir) == 0 {
			t.Fail()
		}

		if len(storage.tokenFilePath) == 0 {
			t.Fail()
		}

		fileInfo, err := os.Stat(storage.baseDir)
		if err != nil {
			t.Fail()
		}

		if !fileInfo.IsDir() {
			t.Fail()
		}
	})

	t.Run("should save and read the same value", func(t *testing.T) {
		persist := &mockStorage{Value: "something"}
		persistClone := &mockStorage{}

		storage.SaveToken(persist)
		storage.ReadToken(persistClone)

		if persist.Value != persistClone.Value {
			t.Fail()
		}
	})

	t.Run("should remove the path after clearing", func(t *testing.T) {
		storage.ClearToken()

		_, err := os.Stat(storage.tokenFilePath)
		if !strings.Contains(err.Error(), "no such file or directory") {
			t.Fail()
		}
	})

}
