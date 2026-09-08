package ast

type AppSpec struct {
	App         string        `json:"app" yaml:"app"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string        `json:"version,omitempty" yaml:"version,omitempty"`
	Rules       []RuleSpec    `json:"rules,omitempty" yaml:"rules,omitempty"`
	Features    []FeatureSpec `json:"features,omitempty" yaml:"features,omitempty"`
	Endpoints   []EndpointSpec `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
}

type RuleSpec struct {
	Statement string `json:"statement" yaml:"statement"`
	Line      int    `json:"-" yaml:"-"`
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
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
	Line int    `json:"-" yaml:"-"`
}

type ActionSpec struct {
	Statement string        `json:"statement" yaml:"statement"`
	Enforce   []EnforceSpec `json:"enforce,omitempty" yaml:"enforce,omitempty"`
	Line      int           `json:"-" yaml:"-"`
}

type EnforceSpec struct {
	Expression string `json:"expression" yaml:"expression"`
	Line       int    `json:"-" yaml:"-"`
}

type EndpointSpec struct {
	Name        string         `json:"name" yaml:"name"`
	Description string         `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string         `json:"version,omitempty" yaml:"version,omitempty"`
	Route       string         `json:"route" yaml:"route"`
	Method      string         `json:"method" yaml:"method"`
	Params      []ParamSpec    `json:"params,omitempty" yaml:"params,omitempty"`
	Responses   []ResponseSpec `json:"responses,omitempty" yaml:"responses,omitempty"`
	Actions     []ActionSpec   `json:"actions,omitempty" yaml:"actions,omitempty"`
	Returns     []ReturnSpec   `json:"returns" yaml:"returns"`
	Line        int            `json:"-" yaml:"-"`
}

type ResponseSpec struct {
	Name   string      `json:"name" yaml:"name"`
	Fields []FieldSpec `json:"fields" yaml:"fields"`
	Line   int         `json:"-" yaml:"-"`
}

type FieldSpec struct {
	Name     string      `json:"name" yaml:"name"`
	Type     string      `json:"type" yaml:"type"`
	SubFields []FieldSpec `json:"sub_fields,omitempty" yaml:"sub_fields,omitempty"`
	Line     int         `json:"-" yaml:"-"`
}

type ReturnSpec struct {
	StatusCode     int    `json:"status_code" yaml:"status_code"`
	Payload        string `json:"payload,omitempty" yaml:"payload,omitempty"`
	PayloadIsString bool  `json:"payload_is_string,omitempty" yaml:"payload_is_string,omitempty"`
	Line           int    `json:"-" yaml:"-"`
}
