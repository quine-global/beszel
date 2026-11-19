package agent

import (
	"bufio"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/henrygd/beszel/agent/utils"
)

// connectionPorts reads CONNECTION_PORTS (for example "80,443"): the local
// TCP ports whose open connections the agent counts, such as a web server's.
// Unset or empty: connections aren't counted.
func connectionPorts() map[uint16]struct{} {
	value, _ := utils.GetEnv("CONNECTION_PORTS")
	ports := make(map[uint16]struct{})
	for _, field := range strings.Split(value, ",") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		port, err := strconv.ParseUint(field, 10, 16)
		if err != nil {
			slog.Warn("CONNECTION_PORTS: not a port", "value", field)
			continue
		}
		ports[uint16(port)] = struct{}{}
	}
	if len(ports) == 0 {
		return nil
	}
	return ports
}

// countConnections counts established TCP connections (IPv4 and IPv6) whose
// local port is one of ports, from /proc/net/tcp and /proc/net/tcp6. Linux
// only: elsewhere the files don't exist and the count is an error.
func countConnections(ports map[uint16]struct{}) (int, error) {
	total := 0
	found := false
	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		n, err := countConnectionsIn(path, ports)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return 0, err
		}
		found = true
		total += n
	}
	if !found {
		return 0, os.ErrNotExist
	}
	return total, nil
}

// countConnectionsIn reads one /proc/net/tcp-format file: after a header
// line, one socket per line, "sl local_address rem_address st ...", with
// local_address as HEXADDR:HEXPORT and st 01 for ESTABLISHED.
func countConnectionsIn(path string, ports map[uint16]struct{}) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	count := 0
	scanner := bufio.NewScanner(f)
	scanner.Scan() // header
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 || fields[3] != "01" {
			continue
		}
		_, hexPort, ok := strings.Cut(fields[1], ":")
		if !ok {
			continue
		}
		port, err := strconv.ParseUint(hexPort, 16, 16)
		if err != nil {
			continue
		}
		if _, wanted := ports[uint16(port)]; wanted {
			count++
		}
	}
	return count, scanner.Err()
}
