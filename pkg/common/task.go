package common

type TaskDefinition struct {
	Id        string   `yaml:"id"`
	Name      string   `yaml:"name"`
	Command   string   `yaml:"cmd"`
	Arguments []string `yaml:"args"`
	Script    *string  `yaml:"script"`
}

type TaskDefinitionList []TaskDefinition
