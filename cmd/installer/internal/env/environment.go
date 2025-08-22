package env

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"claude_analysis/cmd/installer/internal/logger"
)

// EnvironmentConfig represents a domain environment and its associated endpoints
type EnvironmentConfig struct {
	// Domain is the value to send to the login endpoint (e.g., "oa", "swrd")
	Domain string
	// MLOPHosts are the candidate base hosts for GAISF/MLOP gateway
	MLOPHosts []string
	// RegistryHosts are the candidate npm registry mirrors for this domain
	RegistryHosts []string
}

// environmentConfigs defines the available domain environments and their mappings.
// Add new mappings here to support additional domains.
var environmentConfigs = []EnvironmentConfig{
	// Temporary solution for OA, since some computer can access the default registry url.
	{
		Domain:        "oa",
		MLOPHosts:     []string{"https://mlop-azure-gateway.mediatek.inc"},
		RegistryHosts: []string{"https://registry.npmjs.org"},
	},
	{
		Domain:        "oa",
		MLOPHosts:     []string{"https://mlop-azure-gateway.mediatek.inc"},
		RegistryHosts: []string{"https://oa-mirror.mediatek.inc/repository/npm"},
	},
	{
		Domain:        "swrd",
		MLOPHosts:     []string{"https://mlop-azure-rddmz.mediatek.inc"},
		RegistryHosts: []string{"https://swrd-mirror.mediatek.inc/repository/npm"},
	},
}

// Environment is the resolved and connectivity-validated selection
type Environment struct {
	Config      EnvironmentConfig
	MLOPBaseURL string // with scheme
	RegistryURL string // with scheme (or empty to use default)
}

var selectedEnv *Environment

// checkURLReachability performs an HTTP HEAD request to test if URL is reachable
func checkURLReachability(url string, timeout time.Duration) error {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: false},
		},
	}

	// Single lightweight GET request. Many services don't support HEAD reliably.
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "claude-installer/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	io.Copy(io.Discard, io.LimitReader(resp.Body, 512))
	resp.Body.Close()
	// Treat any HTTP response as reachable (even 4xx/5xx indicate server reachable and TLS/DNS ok).
	return nil
}

// SelectAvailableURL resolves and caches the environment selection by testing connectivity
func SelectAvailableURL() *Environment {
	if selectedEnv != nil {
		return selectedEnv
	}

	// Try each configured environment in order; pick the first with a reachable MLOP host
	for _, cfg := range environmentConfigs {
		var chosenMLOP string
		for _, httpsURL := range cfg.MLOPHosts {
			// Use a shorter timeout to avoid long delays when hosts are not reachable.
			if checkURLReachability(httpsURL, 2*time.Second) == nil {
				chosenMLOP = httpsURL
				break
			}
		}
		if chosenMLOP == "" {
			continue
		}
		// Registry is optional; use first reachable, else empty to fall back to default
		var chosenRegistry string
		for _, httpsURL := range cfg.RegistryHosts {
			if checkURLReachability(httpsURL, 2*time.Second) == nil {
				chosenRegistry = httpsURL
				break
			}
		}
		selectedEnv = &Environment{
			Config:      cfg,
			MLOPBaseURL: chosenMLOP,
			RegistryURL: chosenRegistry,
		}
		return selectedEnv
	}

	// As a last resort, fall back to the first environment with HTTPS for MLOP host, without connectivity check
	if len(environmentConfigs) > 0 {
		cfg := environmentConfigs[0]
		mlopHost := cfg.MLOPHosts[0]
		selectedEnv = &Environment{
			Config:      cfg,
			MLOPBaseURL: "https://" + strings.TrimSuffix(strings.TrimPrefix(mlopHost, "https://"), "/"),
			RegistryURL: "", // default registry
		}
		logger.Warning("⚠️ Warning: Falling back to default environment without connectivity check",
			fmt.Sprintf("domain=%s, mlop=%s", cfg.Domain, selectedEnv.MLOPBaseURL))
		return selectedEnv
	}

	// Should not happen; create a stub
	selectedEnv = &Environment{Config: EnvironmentConfig{Domain: "oa"}, MLOPBaseURL: "https://mlop-azure-gateway.mediatek.inc"}
	return selectedEnv
}
