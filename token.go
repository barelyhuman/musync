package main

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type Token struct {
	AccessToken  string `yaml:"access_token"`
	TokenType    string `yaml:"type"`
	RefreshToken string `yaml:"refresh_token"`
	Expiry       string `yaml:"expiry"`
}

func readToken(path string) (*oauth2.Token, error) {
	token := &oauth2.Token{}
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return token, err
	}

	configToken := &Token{}
	yaml.Unmarshal(data, configToken)

	token.AccessToken = configToken.AccessToken
	token.RefreshToken = configToken.RefreshToken
	token.TokenType = configToken.TokenType
	token.Expiry, err = time.Parse(time.RFC1123, configToken.Expiry)
	return token, err
}

func saveToken(token *oauth2.Token, path string) {
	var sb strings.Builder
	sb.WriteString("refresh_token: " + token.RefreshToken + "\n")
	sb.WriteString("type: " + token.TokenType + "\n")
	sb.WriteString("access_token: " + token.AccessToken + "\n")
	sb.WriteString("expiry: " + token.Expiry.Format(time.RFC1123) + "\n")

	ioutil.WriteFile(path, []byte(sb.String()), os.ModePerm)
}
