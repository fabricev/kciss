package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// These 2 structures are used to parse trivy json results partially
type TrivyVulnerability struct {
	Severity string `json:"Severity"`
}

type TrivyResult struct {
	Target          string               `json:"Target"`
	Vulnerabilities []TrivyVulnerability `json:"Vulnerabilities"`
}

func TrivyScanImage(trivyserver string, trivypath string, image string) (sum VulnSummary, err error) {
	sum = VulnSummary{}
	err = nil

	// Invoke trivy on the image
	var cmd *exec.Cmd
	if len(trivyserver) > 0 {
		cmd = exec.Command(trivypath, "client", "--remote", trivyserver, "-f", "json", image)
	} else {
		cmd = exec.Command(trivypath, "-f", "json", image)
	}
	stdout, err := cmd.Output()
	if err != nil {
		return
	}

	// Unmarshal the JSON results
	var results []TrivyResult
	err = json.Unmarshal(stdout, &results)
	if err != nil {
		return
	}

	// Parse vulnerabilities and accumulate results in the image vulnerability summary
	for _, result := range results {
		for _, vuln := range result.Vulnerabilities {
			sum = sum.Add(SummaryForSeverity(fmt.Sprintf("%v", vuln.Severity)))
		}
	}
	return
}
