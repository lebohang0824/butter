package ast

type AppSpec struct {
	App         string           `json:"app" yaml:"app"`
	Description string           `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string           `json:"version,omitempty" yaml:"version,omitempty"`
	Features    []FeatureSpec    `json:"features,omitempty" yaml:"features,omitempty"`
	Endpoints   []EndpointSpec   `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Listeners   []ListenerSpec   `json:"listeners,omitempty" yaml:"listeners,omitempty"`
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
	Statement string           `json:"statement" yaml:"statement"`
	Enforce   []string         `json:"enforce,omitempty" yaml:"enforce,omitempty"`
	Condition *ConditionSpec   `json:"condition,omitempty" yaml:"condition,omitempty"`
	Line      int              `json:"-" yaml:"-"`
}

type ConditionSpec struct {
	Type       string `json:"type" yaml:"type"`
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
	StatusCode     int            `json:"status_code" yaml:"status_code"`
	Payload        string         `json:"payload,omitempty" yaml:"payload,omitempty"`
	PayloadIsString bool          `json:"payload_is_string,omitempty" yaml:"payload_is_string,omitempty"`
	Condition      *ConditionSpec `json:"condition,omitempty" yaml:"condition,omitempty"`
	Line           int            `json:"-" yaml:"-"`
}

type ListenerSpec struct {
	Name        string               `json:"name" yaml:"name"`
	Description string               `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string               `json:"version,omitempty" yaml:"version,omitempty"`
	Topic       string               `json:"topic" yaml:"topic"`
	Params      []ParamSpec          `json:"params,omitempty" yaml:"params,omitempty"`
	Actions     []ActionSpec         `json:"actions,omitempty" yaml:"actions,omitempty"`
	Returns     []ListenerReturnSpec `json:"returns" yaml:"returns"`
	Line        int                  `json:"-" yaml:"-"`
}

type ListenerReturnSpec struct {
	State     string         `json:"state" yaml:"state"`
	Condition *ConditionSpec `json:"condition,omitempty" yaml:"condition,omitempty"`
	Line      int            `json:"-" yaml:"-"`
}
