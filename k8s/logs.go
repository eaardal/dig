package k8s

import (
	"context"
	"fmt"
	"github.com/eaardal/dig/ui"
	"github.com/eaardal/dig/utils"
	"io"
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type LogMsg struct {
	Origin   string
	LogChunk []byte
}

// ReadLogs fetches logs from all pods in the specified Kubernetes deployments in parallel.
func ReadLogs(ctx context.Context, client *kubernetes.Clientset, namespace string, deploymentNames []string, sinkCh chan<- *LogMsg) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(deploymentNames))

	for _, deploymentName := range deploymentNames {
		wg.Add(1)

		go func(deployment string) {
			defer wg.Done()

			pods, err := GetPodsForDeployment(ctx, client, namespace, deployment)
			if err != nil {
				errCh <- fmt.Errorf("failed to get pods for deployment %s: %w", deployment, err)
				return
			}

			for _, pod := range pods {
				wg.Add(1)
				go func(pod v1.Pod) {
					defer wg.Done()

					if err := getPodLogs(ctx, client, namespace, pod.Name, sinkCh); err != nil {
						errCh <- fmt.Errorf("failed to fetch logs for pod %s: %w", pod.Name, err)
					}
				}(pod)
			}
		}(deploymentName)
	}

	wg.Wait()
	close(sinkCh)

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func getPodLogs(ctx context.Context, client *kubernetes.Clientset, namespace, podName string, sinkCh chan<- *LogMsg) error {
	ui.Write("Fetching logs for pod %s...", podName)

	logsRequest := client.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{})
	logStream, err := logsRequest.Stream(ctx)
	if err != nil {
		return fmt.Errorf("failed to open log stream for pod %s: %w", podName, err)
	}
	defer utils.CloseOrPanic(logStream.Close)

	buf := make([]byte, 2000) // Read chunks of 2KB

	for {
		n, err := logStream.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading log stream: %w", err)
		}
		if n == 0 {
			break
		}

		// Create a new slice with just the data read to avoid overwriting the buffer in the next iteration of the for-loop.
		logChunk := make([]byte, n)
		copy(logChunk, buf[:n])

		logMsg := &LogMsg{
			Origin:   podName,
			LogChunk: logChunk,
		}
		sinkCh <- logMsg
	}

	return nil
}
