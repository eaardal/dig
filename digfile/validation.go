package digfile

import "fmt"

func ValidateJobNamesAreUnique(jobs []*Job) error {
	jobNames := make(map[string]bool)
	for _, job := range jobs {
		if jobNames[job.Name] {
			return fmt.Errorf("job names must be unique, but job name %s is used more than once", job.Name)
		}
		jobNames[job.Name] = true
	}
	return nil
}

func ValidateKubernetesConfigExistsForJobs(jobs []*Job) error {
	for _, job := range jobs {
		return ValidateKubernetesConfigExistsForJob(job)
	}
	return nil
}

func ValidateKubernetesConfigExistsForJob(job *Job) error {
	if job.Kubernetes == nil {
		return fmt.Errorf("job %s is missing Kubernetes configuration", job.Name)
	}
	if job.Kubernetes.ContextName == "" {
		return fmt.Errorf("job %s is missing Kubernetes context name", job.Name)
	}
	if job.Kubernetes.Namespace == "" {
		return fmt.Errorf("job %s is missing Kubernetes namespace", job.Name)
	}
	if len(job.Kubernetes.DeploymentNames) == 0 {
		return fmt.Errorf("job %s is missing Kubernetes deployment names", job.Name)
	}
	return nil
}

func ValidateOnlyOneDefaultJob(jobs []*Job) error {
	defaultJobCount := 0
	for _, job := range jobs {
		if job.IsDefault {
			defaultJobCount++
		}
	}
	if defaultJobCount > 1 {
		return fmt.Errorf("only one job can be marked as default, but %d jobs are marked as default", defaultJobCount)
	}
	return nil
}
