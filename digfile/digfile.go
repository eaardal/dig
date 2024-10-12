package digfile

type KubernetesJob struct {
	ClusterName     string   `json:"clusterName,omitempty"`
	Namespace       string   `yaml:"namespace,omitempty"`
	DeploymentNames []string `yaml:"deploymentNames,omitempty"`
}

type Job struct {
	Name       string         `yaml:"name"`
	Kubernetes *KubernetesJob `yaml:"kubernetes,omitempty"`
}

type Digfile struct {
	Jobs []*Job `yaml:"jobs"`
}

var DefaultDigfile = &Digfile{
	Jobs: []*Job{},
}
