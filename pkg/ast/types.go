package ast

type AppSpec struct {
	App         string        `json:"app"`
	Description string        `json:"description,omitempty"`
	Features    []FeatureSpec `json:"features"`
}

type FeatureSpec struct {
	Name    string       `json:"name"`
	Params  []ParamSpec  `json:"params,omitempty"`
	Actions []ActionSpec `json:"actions,omitempty"`
}

type ParamSpec struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Required bool        `json:"required"`
	Default  interface{} `json:"default,omitempty"`
}

type ActionSpec struct {
	Statement string         `json:"statement"`
	Condition *ConditionSpec `json:"condition,omitempty"`
}

type ConditionSpec struct {
	Type       string `json:"type"`
	Expression string `json:"expression"`
}
