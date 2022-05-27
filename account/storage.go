package account

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

type TokenStorage struct {
	baseDir       string
	tokenFilePath string
}

func (ts *TokenStorage) ReadToken(a interface{}) error {
	tokenData, err := os.ReadFile(ts.tokenFilePath)
	if err != nil {
		return err
	}
	json.Unmarshal(tokenData, a)
	return nil
}

func (ts *TokenStorage) SaveToken(b interface{}) error {
	stringifiedToken, err := json.Marshal(b)
	if err != nil {
		return err
	}
	os.WriteFile(ts.tokenFilePath, stringifiedToken, os.ModePerm)
	return nil
}

func (ts *TokenStorage) ClearToken() error {
	return os.Remove(ts.tokenFilePath)
}

func NewTokenStorage() *TokenStorage {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panic("Failed to get user home directory", err)
	}
	baseDir := path.Join(homedir, ".config", "musync")
	tokenFilePath := path.Join(baseDir, "musync.json")
	os.MkdirAll(baseDir, os.ModePerm)
	return &TokenStorage{baseDir: baseDir, tokenFilePath: tokenFilePath}
}
