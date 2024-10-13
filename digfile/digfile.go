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

func (d *Digfile) HasAnyDefaultJob() bool {
	for _, job := range d.Jobs {
		if job.IsDefault {
			return true
		}
	}
	return false
}

func (d *Digfile) HasJob(name string) bool {
	for _, job := range d.Jobs {
		if job.Name == name {
			return true
		}
	}
	return false
}

func (d *Digfile) GetDefaultJob() *Job {
	for _, job := range d.Jobs {
		if job.IsDefault {
			return job
		}
	}
	return nil
}

func (d *Digfile) GetJobByName(job string) *Job {
	for _, j := range d.Jobs {
		if j.Name == job {
			return j
		}
	}
	return nil
}

func (d *Digfile) GetJobByIndex(index int) (*Job, error) {
	if index < 0 || index >= len(d.Jobs) {
		return nil, fmt.Errorf("job index out of bounds")
	}
	return d.Jobs[index], nil
}

func (d *Digfile) GetJob(jobName *string, jobIndex *int) (*Job, error) {
	if jobName == nil && jobIndex == nil {
		return d.GetDefaultJob(), nil
	} else if jobIndex != nil {
		return d.GetJobByIndex(*jobIndex)
	} else {
		return d.GetJobByName(*jobName), nil
	}
}

func (d *Digfile) Validate() []error {
	var errs []error

	if err := ValidateJobNamesAreUnique(d.Jobs); err != nil {
		errs = append(errs, err)
	}

	if err := ValidateKubernetesConfigExistsForJobs(d.Jobs); err != nil {
		errs = append(errs, err)
	}

	if err := ValidateOnlyOneDefaultJob(d.Jobs); err != nil {
		errs = append(errs, err)
	}

	return errs
}
