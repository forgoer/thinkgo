package helper

import (
	"os"
	"strings"
)

func ParseAddr(addrs ...string) string {
	var addr = "0.0.0.0"
	var port = "9011"
	switch len(addrs) {
	case 0:
		if a := os.Getenv("THINKGO_ADDR"); a != "" {
			addr = a
		}
		if p := os.Getenv("THINKGO_PORT"); p != "" {
			port = p
		}
	case 1:
		strs := strings.Split(addrs[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			addr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			port = strs[1]
		}
	default:

	}
	return addr + ":" + port
}
