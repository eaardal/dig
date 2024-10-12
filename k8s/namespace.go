package k8s

import (
	"fmt"
	"github.com/eaardal/dig/config"
	"log"
	"os"
	"os/exec"
	"strings"
)

func findNamespace() (string, error) {
	env, found := os.LookupEnv(config.KubernetesNamespaceEnvVar)
	if found {
		return env, nil
	}

	if !hasExe("kubens") {
		log.Fatalln("Could not resolve the Kubernetes namespace to use because the first two methods of looking up the namespace didn't give any results and it appears kubens is not available on your PATH. Use one of these three methods to provide your namespace: 1) Specify the namespace along with the name to search for: <namespace>/<appName>. 2) Set the PODID_NAMESPACE environment variable. 3) Use kubens (https://github.com/ahmetb/kubectx) to set the current namespace. This is the recommended method. The namespace is resolved in the order these alternatives are listed.")
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

func hasExe(exe string) bool {
	path, err := exec.LookPath(exe)
	return err == nil && len(path) > 0
}
