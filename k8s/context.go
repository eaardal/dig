package k8s

import (
	"fmt"
	"github.com/eaardal/dig/config"
	"github.com/eaardal/dig/utils"
	"os"
	"os/exec"
	"strings"
)

func ResolveContext() (string, error) {
	env, found := os.LookupEnv(config.KubernetesContextEnvVar)
	if found {
		return env, nil
	}

	if !utils.HasExeOnPath("kubectx") {
		return "", fmt.Errorf("could not resolve kubernetes context using kubectx")
	}

	kubectxContext, err := exec.Command("kubectx", "--current").Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute kubectx: %w", err)
	}

	if len(kubectxContext) <= 0 {
		return "", fmt.Errorf("found no context by invoking `kubectx --current`. Ensure kubectx is set up correctly")
	}

	return strings.TrimSuffix(string(kubectxContext), "\n"), nil
}
