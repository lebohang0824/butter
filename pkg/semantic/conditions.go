package semantic

import (
	"strings"
	"unicode"
)

func extractParamRefs(expr string) []string {
	tokens := strings.Fields(expr)
	var refs []string
	for _, tok := range tokens {
		tok = strings.TrimRight(tok, ",;:.!?")
		if len(tok) > 0 && unicode.IsUpper(rune(tok[0])) {
			refs = append(refs, tok)
		}
	}
	return refs
}

func (a *Analyzer) checkConditionRefs() {
	for _, f := range a.app.Features {
		paramNames := make(map[string]bool)
		for _, p := range f.Params {
			paramNames[p.Name] = true
		}
		for _, act := range f.Actions {
			if act.Condition == nil {
				continue
			}
			refs := extractParamRefs(act.Condition.Expression)
			for _, ref := range refs {
				if !paramNames[ref] {
					a.addError(act.Condition.Line, "undefined parameter %q referenced in condition of action in feature %q", ref, f.Name)
				}
			}
		}
	}
}
