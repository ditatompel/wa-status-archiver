package helpers

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/types"
)

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// Copyright (c) 2021 Tulir Asokan
func ParseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			fmt.Printf("Invalid JID %s: %v\n", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			fmt.Printf("Invalid JID %s: no server specified\n", arg)
			return recipient, false
		}
		return recipient, true
	}
}

func StrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	if !strings.ContainsAny(password, "0123456789") {
		return false
	}
	if !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return false
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return false
	}
	if !strings.ContainsAny(password, "!@#$%^&*()_+|~-=`{}[]:;<>?,./") {
		return false
	}
	return true
}

// No special character and unicode for username
func ValidUsername(username string) bool {
	if len(username) < 5 || len(username) > 20 {
		return false
	}

	// reject witespace, tabs, newlines, and other special characters
	if strings.ContainsAny(username, " \t\n") {
		return false
	}
	// reject unicode
	if strings.ContainsAny(username, "^\x00-\x7F") {
		return false
	}
	// reject special characters
	if strings.ContainsAny(username, "!@#$%^&*()_+|~-=`{}[]:;<>?,./ ") { // note last blank space
		return false
	}

	if !strings.ContainsAny(username, "abcdefghijklmnopqrstuvwxyz0123456789") {
		return false
	}
	return true
}
