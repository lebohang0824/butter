package semantic

import "fmt"

type Severity int

const (
	SemError   Severity = iota
	SemWarning
)

func (s Severity) String() string {
	switch s {
	case SemError:
		return "Error"
	case SemWarning:
		return "Warning"
	default:
		return "Unknown"
	}
}

type Diagnostic struct {
	Line     int
	Severity Severity
	Message  string
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("%s: line %d: %s", d.Severity, d.Line, d.Message)
}
