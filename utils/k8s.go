package utils

import (
	"regexp"
	"strings"
)

func IsValueKubernetesPodID(value string) bool {
	var podIDRegex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9.]*[a-z0-9])?$`)
	return podIDRegex.MatchString(value)
}

func SplitIntoKubernetesPodIDParts(value string) (appName, deploymentID, replicaSetID string) {
	parts := strings.Split(value, "-")
	appName = strings.Join(parts[:len(parts)-2], "-")
	deploymentID = strings.Join(parts[:len(parts)-1], "-")
	replicaSetID = parts[len(parts)-1]
	return
}
