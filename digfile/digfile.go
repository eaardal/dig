package digfile

import (
	"fmt"
	"github.com/eaardal/dig/ui"
)

type KubernetesJob struct {
	ContextName     string   `yaml:"contextName,omitempty"`
	Namespace       string   `yaml:"namespace,omitempty"`
	DeploymentNames []string `yaml:"deploymentNames,omitempty"`
}

func (k *KubernetesJob) Print(indent string) {
	if k.ContextName != "" {
		ui.Write("%sContext: %s", indent, k.ContextName)
	} else {
		ui.Write("%sContext: <not set>", indent)
	}

	if k.Namespace != "" {
		ui.Write("%sNamespace: %s", indent, k.Namespace)
	} else {
		ui.Write("%sNamespace: <not set>", indent)
	}

	if len(k.DeploymentNames) > 0 {
		ui.Write("%sDeployments:", indent)
		for _, deployment := range k.DeploymentNames {
			ui.Write("%s  %s", indent, deployment)
		}
	} else {
		ui.Write("%sDeployments: <not set>", indent)
	}
}

type Job struct {
	Name       string         `yaml:"name"`
	Kubernetes *KubernetesJob `yaml:"kubernetes,omitempty"`
	IsDefault  bool           `yaml:"isDefault,omitempty"`
}

func (j *Job) Print() {
	ui.Write("Job %s:", j.Name)
	ui.Write("%sIs default: %t", ui.Indent, j.IsDefault)

	if j.Kubernetes != nil {
		ui.Write("%sKubernetes:", ui.Indent)
		j.Kubernetes.Print(ui.Indent2)
	} else {
		ui.Write("%sKubernetes: <not set>", ui.Indent)
	}
}

type Digfile struct {
	Jobs []*Job `yaml:"jobs"`
}

var DefaultDigfile = &Digfile{
	Jobs: []*Job{},
}

func (d *Digfile) SetDefaultJob(jobName string) {
	for _, job := range d.Jobs {
		job.IsDefault = job.Name == jobName
	}
}

func (d *Digfile) SetAllJobsNotDefault() {
	for _, job := range d.Jobs {
		job.IsDefault = false
	}
}

func (d *Digfile) GetDefaultJob() *Job {
	for _, job := range d.Jobs {
		if job.IsDefault {
			return job
		}
	}
	return nil
}

func (d *Digfile) IsAnyJobDefault() bool {
	for _, job := range d.Jobs {
		if job.IsDefault {
			return true
		}
	}
	return false
}

func (d *Digfile) Validate() []error {
	var errs []error

	if err := d.validateJobNamesAreUnique(); err != nil {
		errs = append(errs, err)
	}

	if err := d.validateKubernetesConfigExists(); err != nil {
		errs = append(errs, err)
	}

	if err := d.validateOnlyOneDefaultJob(); err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (d *Digfile) validateJobNamesAreUnique() error {
	jobNames := make(map[string]bool)
	for _, job := range d.Jobs {
		if jobNames[job.Name] {
			return fmt.Errorf("job names must be unique, but job name %s is used more than once", job.Name)
		}
		jobNames[job.Name] = true
	}
	return nil
}

func (d *Digfile) validateKubernetesConfigExists() error {
	for _, job := range d.Jobs {
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
	}
	return nil
}

func (d *Digfile) validateOnlyOneDefaultJob() error {
	defaultJobCount := 0
	for _, job := range d.Jobs {
		if job.IsDefault {
			defaultJobCount++
		}
	}
	if defaultJobCount > 1 {
		return fmt.Errorf("only one job can be marked as default, but %d jobs are marked as default", defaultJobCount)
	}
	return nil
}

func (d *Digfile) HasJob(name string) bool {
	for _, job := range d.Jobs {
		if job.Name == name {
			return true
		}
	}
	return false
}
