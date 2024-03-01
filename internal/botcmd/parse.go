package botcmd

import (
	"strings"
)

func ParseCmd(message string) string {
	if message == ".help" {
		return help()
	}
	if cut, found := strings.CutPrefix(message, ".ip "); found {
		return cmdIP(cut)
	}

	return ""
}

func help() string {
	return `Available Commands:
.ip <ip> - Get information about an IP address
`
}
