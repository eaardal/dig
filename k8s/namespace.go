package k8s

import (
	"fmt"
	"github.com/eaardal/dig/config"
	"github.com/eaardal/dig/utils"
	"os"
	"os/exec"
	"strings"
)

func ResolveNamespace() (string, error) {
	env, found := os.LookupEnv(config.KubernetesNamespaceEnvVar)
	if found {
		return env, nil
	}

	if !utils.HasExeOnPath("kubens") {
		return "", fmt.Errorf("could not resolve kubernetes namespace using kubens")
	}

	kubensNamespace, err := exec.Command("kubens", "--current").Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute kubens: %w", err)
	}

	if len(kubensNamespace) <= 0 {
		return "", fmt.Errorf("found no namespace by invoking `kubens --current`. Ensure kubens is set up correctly")
	}

	return strings.TrimSuffix(string(kubensNamespace), "\n"), nil
}
