package formatter

import (
	"strings"
	"testing"
)

func TestFormatLineSpacing(t *testing.T) {
	input := `app TodoApp
  description "A todo application"
  version "1.0.0"
  rules
    "Users may only access their own tasks"
    "A task cannot be completed until it is assigned to a user"

feature CreateTodo
  description "Creates a todo"
  version "1.0.0"
  params
    title string
    completed boolean
  actions
    "Validate the title is not empty"
    "Create the todo record"

endpoint SaveTodo "/todo"
  method "POST"
  params
    title string
  responses
    TodoResponse
      id integer
  actions
    "Store the todo"
  returns
    201 TodoResponse
    500 "Internal server error"
`

	got, err := Format([]byte(input))
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := `app TodoApp
  description "A todo application"
  version "1.0.0"

  rules
    "Users may only access their own tasks"
    "A task cannot be completed until it is assigned to a user"

feature CreateTodo
  description "Creates a todo"
  version "1.0.0"

  params
    title string
    completed boolean

  actions
    "Validate the title is not empty"
    "Create the todo record"

endpoint SaveTodo "/todo"
  method "POST"

  params
    title string

  responses
    TodoResponse
      id integer

  actions
    "Store the todo"

  returns
    201 TodoResponse
    500 "Internal server error"
`

	if string(got) != want {
		t.Errorf("unexpected formatted output:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestFormatPreservesFinalNewlineAfterKeywordValueLine(t *testing.T) {
	input := `# This is a comment
app MyApp  # inline comments work too
`

	got, err := Format([]byte(input))
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}
	if string(got) != input {
		t.Errorf("expected output to be preserved, got:\n%s", got)
	}
}

func TestFormatRemovesBlankLinesAfterSectionKeywords(t *testing.T) {
	input := `app MyApp
  version "1.0.0"
  description "An app"

  description "Another app"
  version "2.0.0"
`

	got, err := Format([]byte(input))
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	if strings.Contains(string(got), "\n\n") {
		t.Errorf("expected blank lines after section keywords to be removed, got:\n%s", got)
	}
}

func TestFormatNormalizesTabs(t *testing.T) {
	input := "feature Foo\n\t  description \"x\"\n    version \"1.0.0\"\n\t\t  actions\n        \"Do something\"\n"

	got, err := Format([]byte(input))
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	want := "feature Foo\n    description \"x\"\n    version \"1.0.0\"\n\n      actions\n        \"Do something\"\n"
	if string(got) != want {
		t.Errorf("unexpected tab normalization:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestFormatIdempotent(t *testing.T) {
	input := `app TodoApp
  description "A todo application"
  version "1.0.0"

  rules
    "Users may only access their own tasks"

feature CreateTodo
  description "Creates a todo"
  version "1.0.0"

  params
    title string

  actions
    "Validate the title is not empty"
    "Create the todo record"
`

	first, err := Format([]byte(input))
	if err != nil {
		t.Fatalf("first Format returned error: %v", err)
	}
	second, err := Format(first)
	if err != nil {
		t.Fatalf("second Format returned error: %v", err)
	}
	if string(first) != string(second) {
		t.Errorf("Format is not idempotent:\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
}

func TestFormatFileEndingInKeywordValueLineKeepsIdempotency(t *testing.T) {
	input := `# This is a comment
app MyApp  # inline comments work too
`

	first, err := Format([]byte(input))
	if err != nil {
		t.Fatalf("first Format returned error: %v", err)
	}
	second, err := Format(first)
	if err != nil {
		t.Fatalf("second Format returned error: %v", err)
	}
	if string(first) != string(second) {
		t.Errorf("Format is not idempotent for keyword-value endings:\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
}