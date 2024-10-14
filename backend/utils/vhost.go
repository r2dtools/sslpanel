package utils

import (
	"strings"

	"github.com/r2dtools/agentintegration"
)

// FilterVhosts filter virtual hosts
func FilterVhosts(vhosts []agentintegration.VirtualHost) []agentintegration.VirtualHost {
	var rVhosts []agentintegration.VirtualHost

	for _, vhost := range vhosts {
		serverName := strings.Trim(vhost.ServerName, ".")
		serverNameParts := strings.Split(serverName, ".")

		// skip vhost names like 'domain'
		if len(serverNameParts) > 1 {
			rVhosts = append(rVhosts, vhost)
		}
	}

	return rVhosts
}
