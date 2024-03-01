package botcmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"strings"
)

func cmdIP(cmd string) string {
	ip := strings.TrimSpace(cmd)
	if ip == "" {
		return printIpHelp()
	}

	ip = strings.Split(ip, " ")[0]
	if ip == "" {
		return printIpHelp()
	}

	if _, err := netip.ParseAddr(ip); err != nil {
		return "Invalid IP"
	}

	return fetchIpInfo(ip)
}

func printIpHelp() string {
	return `
	Command: .ip <ip>
	Example: .ip 8.8.8.8
	`
}

func fetchIpInfo(ip string) string {
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return "Error processing your request"
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "Error processing your request"
	}

	return fmt.Sprintf("IP: %s\nCountry: %s", data["query"], data["country"])
}
