package common

type TaskResultType = string

const (
	TaskResultTypeJson  TaskResultType = "json"
	TaskResultTypeYaml  TaskResultType = "yaml"
	TaskResultTypePlain TaskResultType = "plain"
)

type TaskDefinition struct {
	Id         string          `yaml:"id"`
	Name       string          `yaml:"name"`
	Command    string          `yaml:"cmd"`
	Arguments  []string        `yaml:"args"`
	Script     *string         `yaml:"script"`
	ResultType *TaskResultType `yaml:"result"`
}

type TaskDefinitionList []TaskDefinition
