package ast

type AppSpec struct {
	App         string        `json:"app" yaml:"app"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string        `json:"version,omitempty" yaml:"version,omitempty"`
	Features    []FeatureSpec `json:"features" yaml:"features"`
}

type FeatureSpec struct {
	Name        string       `json:"name" yaml:"name"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string       `json:"version,omitempty" yaml:"version,omitempty"`
	Params      []ParamSpec  `json:"params,omitempty" yaml:"params,omitempty"`
	Actions     []ActionSpec `json:"actions,omitempty" yaml:"actions,omitempty"`
	Line        int          `json:"-" yaml:"-"`
}

type ParamSpec struct {
	Name     string      `json:"name" yaml:"name"`
	Type     string      `json:"type" yaml:"type"`
	Required bool        `json:"required" yaml:"required"`
	Default  interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	Validate []string    `json:"validate,omitempty" yaml:"validate,omitempty"`
	Length   int         `json:"length,omitempty" yaml:"length,omitempty"`
	Line     int         `json:"-" yaml:"-"`
}

type ActionSpec struct {
	Statement string         `json:"statement" yaml:"statement"`
	Condition *ConditionSpec `json:"condition,omitempty" yaml:"condition,omitempty"`
	Line      int            `json:"-" yaml:"-"`
}

type ConditionSpec struct {
	Type       string `json:"type" yaml:"type"`
	Expression string `json:"expression" yaml:"expression"`
	Line       int    `json:"-" yaml:"-"`
}
