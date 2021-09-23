package ast

type Node struct {
	Position int `yaml:"pos"`
}

type Path []Step

type Step struct {
	Node `yaml:",inline"`

	Index    *int    `yaml:"idx,omitempty"`
	Identity *string `yaml:"id,omitempty"`
}

type Argument struct {
	Node `yaml:",inline"`

	Number *int    `yaml:"number,omitempty"`
	String *string `yaml:"string,omitempty"`
	Path   Path    `yaml:"path,omitempty"`
}

type Command struct {
	Node `yaml:",inline"`

	Path   Path       `yaml:"path,omitempty"`
	Equal  []Argument `yaml:"equal,omitempty"`
	Assign []Argument `yaml:"assign,omitempty"`
	Merge  []Argument `yaml:"merge,omitempty"`
	Pipe   *Command   `yaml:"pipe,omitempty"`
}
