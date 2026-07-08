package simulator

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"butter/pkg/ast"
	"butter/pkg/output"
)

func init() {
	output.Register(simExt{})
}

type simExt struct{}

func (simExt) Name() string          { return "sim" }
func (simExt) FileExtension() string { return ".sim.html" }

func (simExt) Serialize(spec *ast.AppSpec) ([]byte, error) {
	apiKey := os.Getenv("BUTTER_AI_API_KEY")

	fmt.Print("Step 1/2: Parsing spec...                    ")
	specJSON, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal spec: %w", err)
	}
	fmt.Println("Done")

	var provider string
	var generatedCode string
	envProvider := os.Getenv("BUTTER_AI_PROVIDER")

	if apiKey == "" && envProvider == "opencode" {
		apiKey = "opencode"
		provider = "opencode"
	}

	if apiKey == "" && provider == "" {
		provider, apiKey = promptProviderAndKey()
	} else if apiKey != "" && provider == "" {
		provider = detectProvider(apiKey)
		if provider == "ChatGPT" && envProvider == "" {
			provider = promptProvider()
		}
	}

	if provider != "" && apiKey != "" {
		fmt.Printf("Using %s\n", provider)

		fmt.Print("Step 2/2: Generating implementation code...  ")
		codegenResult, err := callAI(provider, apiKey, buildCodePrompt(spec), 4000)
		if err != nil {
			return nil, fmt.Errorf("code generation failed: %w", err)
		}
		generatedCode = extractCode(codegenResult)
		// Validate the AI actually returned JS code, not natural language
		if generatedCode != "" && !isValidAICode(generatedCode) {
			fmt.Println("\n  Warning: AI response doesn't look like JS code — skipping AI integration")
			fmt.Printf("  Response preview: %s...\n", truncate(generatedCode, 100))
			generatedCode = ""
			provider = "built-in"
		} else {
			tokenEst := len(strings.Fields(generatedCode))
			fmt.Printf("Done (~%d tokens)\n", tokenEst)
		}
	} else {
		fmt.Println("Step 2/2: Code generation skipped")
		provider = "built-in"
	}

	var buf bytes.Buffer
	data := map[string]interface{}{
		"Spec":          template.JS(specJSON),
		"GeneratedCode": string(generatedCode),
		"AIExecutable":  template.JS(generatedCode),
		"HasCode":       generatedCode != "",
		"Provider":      provider,
		"AppName":       spec.App,
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("template execution failed: %w", err)
	}
	return buf.Bytes(), nil
}

var providerNames = map[string]string{
	"anthropic": "Anthropic",
	"chatgpt":   "ChatGPT",
	"gemini":    "Gemini",
	"deepseek":  "DeepSeek",
	"opencode":  "opencode",
}

var stdinReader = bufio.NewReader(os.Stdin)

func readLine() string {
	line, err := stdinReader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(line)
}

func promptProviderAndKey() (string, string) {
	fmt.Println("Select AI provider:")
	fmt.Println("  1) Anthropic")
	fmt.Println("  2) ChatGPT (OpenAI)")
	fmt.Println("  3) Gemini (Google)")
	fmt.Println("  4) DeepSeek")
	fmt.Println("  5) opencode (free, no key needed)")
	fmt.Println("  6) none (skip AI generation)")
	fmt.Print("Enter number [5]: ")

	choice := strings.TrimSpace(readLine())
	if choice == "" {
		choice = "5"
	}

	switch choice {
	case "1":
		fmt.Print("Enter Anthropic API key: ")
		key := readLine()
		if key == "" {
			return "", ""
		}
		return "Anthropic", key
	case "2":
		fmt.Print("Enter OpenAI API key: ")
		key := readLine()
		if key == "" {
			return "", ""
		}
		return "ChatGPT", key
	case "3":
		fmt.Print("Enter Gemini API key: ")
		key := readLine()
		if key == "" {
			return "", ""
		}
		return "Gemini", key
	case "4":
		fmt.Print("Enter DeepSeek API key: ")
		key := readLine()
		if key == "" {
			return "", ""
		}
		return "DeepSeek", key
	case "5":
		return "opencode", "opencode"
	case "6":
		return "", ""
	default:
		fmt.Println("Invalid choice, defaulting to opencode")
		return "opencode", "opencode"
	}
}

func promptProvider() string {
	for {
		fmt.Print("Provider? (anthropic/chatgpt/gemini/deepseek/opencode/none) [chatgpt]: ")
		p := strings.ToLower(readLine())
		if p == "" {
			return "ChatGPT"
		}
		if p == "none" {
			return ""
		}
		if name, ok := providerNames[p]; ok {
			return name
		}
	}
}

func detectProvider(key string) string {
	if p := os.Getenv("BUTTER_AI_PROVIDER"); p != "" {
		if name, ok := providerNames[strings.ToLower(p)]; ok {
			return name
		}
	}
	if strings.HasPrefix(key, "sk-ant-") {
		return "Anthropic"
	}
	if strings.HasPrefix(key, "sk-") {
		return "ChatGPT"
	}
	return "Gemini"
}



func callAI(provider, apiKey, prompt string, maxTokens int) (string, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	switch provider {
	case "Anthropic":
		return callAnthropic(client, apiKey, prompt, maxTokens)
	case "ChatGPT":
		return callOpenAI(client, apiKey, prompt, maxTokens)
	case "Gemini":
		return callGemini(client, apiKey, prompt, maxTokens)
	case "DeepSeek":
		return callDeepSeek(client, apiKey, prompt, maxTokens)
	case "opencode":
		return callOpenCode(prompt, maxTokens)
	default:
		return "", fmt.Errorf("unknown provider: %s", provider)
	}
}

func callOpenCode(prompt string, maxTokens int) (string, error) {
	model := os.Getenv("BUTTER_AI_MODEL")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if model != "" {
		cmd = exec.CommandContext(ctx, "opencode", "run", "--model", model, "--format", "json")
	} else {
		cmd = exec.CommandContext(ctx, "opencode", "run", "--format", "json")
	}
	cmd.Stdin = strings.NewReader(prompt)
	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("opencode timed out after 300s")
		}
		if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
			return "", fmt.Errorf("opencode failed: %s", string(ee.Stderr))
		}
		return "", fmt.Errorf("opencode failed: %w", err)
	}

	var parts []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev struct {
			Type string `json:"type"`
			Part struct {
				Text string `json:"text"`
			} `json:"part"`
		}
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		if ev.Type == "text" {
			parts = append(parts, ev.Part.Text)
		}
	}
	return strings.TrimSpace(strings.Join(parts, "\n")), nil
}

func callAnthropic(client *http.Client, apiKey, prompt string, maxTokens int) (string, error) {
	body := map[string]interface{}{
		"model":      "claude-sonnet-4-20250514",
		"max_tokens": maxTokens,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(bodyJSON))
	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")
	return doRequestWithRetry(client, req, func(r io.Reader) (string, error) {
		var res struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.NewDecoder(r).Decode(&res); err != nil {
			return "", err
		}
		if res.Error.Message != "" {
			return "", fmt.Errorf("Anthropic API error: %s", res.Error.Message)
		}
		if len(res.Content) == 0 {
			return "", fmt.Errorf("empty Anthropic response")
		}
		return res.Content[0].Text, nil
	}, 2)
}

func callOpenAI(client *http.Client, apiKey, prompt string, maxTokens int) (string, error) {
	body := map[string]interface{}{
		"model":      "gpt-4o",
		"max_tokens": maxTokens,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(bodyJSON))
	if err != nil {
		return "", err
	}
	req.Header.Set("authorization", "Bearer "+apiKey)
	req.Header.Set("content-type", "application/json")
	return doRequestWithRetry(client, req, func(r io.Reader) (string, error) {
		var res struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.NewDecoder(r).Decode(&res); err != nil {
			return "", err
		}
		if res.Error.Message != "" {
			return "", fmt.Errorf("OpenAI API error: %s", res.Error.Message)
		}
		if len(res.Choices) == 0 {
			return "", fmt.Errorf("empty OpenAI response")
		}
		return res.Choices[0].Message.Content, nil
	}, 2)
}

func callGemini(client *http.Client, apiKey, prompt string, maxTokens int) (string, error) {
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyJSON))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	return doRequestWithRetry(client, req, func(r io.Reader) (string, error) {
		var res struct {
			Candidates []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			} `json:"candidates"`
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.NewDecoder(r).Decode(&res); err != nil {
			return "", err
		}
		if res.Error.Message != "" {
			return "", fmt.Errorf("Gemini API error: %s", res.Error.Message)
		}
		if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
			return "", fmt.Errorf("empty Gemini response")
		}
		return res.Candidates[0].Content.Parts[0].Text, nil
	}, 2)
}

type rateLimitError struct{ msg string }

func (e *rateLimitError) Error() string { return e.msg }

func extractErrorBody(status int, body io.Reader) error {
	data, _ := io.ReadAll(body)
	var apiErr struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if json.Unmarshal(data, &apiErr) == nil && apiErr.Error.Message != "" {
		msg := apiErr.Error.Message
		if len(msg) > 200 {
			msg = msg[:200] + "..."
		}
		err := fmt.Errorf("API returned status %d: %s", status, msg)
		if status == http.StatusTooManyRequests {
			return &rateLimitError{err.Error()}
		}
		return err
	}
	msg := string(data)
	if len(msg) > 200 {
		msg = msg[:200] + "..."
	}
	err := fmt.Errorf("API returned status %d: %s", status, msg)
	if status == http.StatusTooManyRequests {
		return &rateLimitError{err.Error()}
	}
	return err
}

func callDeepSeek(client *http.Client, apiKey, prompt string, maxTokens int) (string, error) {
	body := map[string]interface{}{
		"model":      "deepseek-chat",
		"max_tokens": maxTokens,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.deepseek.com/chat/completions", bytes.NewReader(bodyJSON))
	if err != nil {
		return "", err
	}
	req.Header.Set("authorization", "Bearer "+apiKey)
	req.Header.Set("content-type", "application/json")
	return doRequestWithRetry(client, req, func(r io.Reader) (string, error) {
		var res struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.NewDecoder(r).Decode(&res); err != nil {
			return "", err
		}
		if res.Error.Message != "" {
			return "", fmt.Errorf("DeepSeek API error: %s", res.Error.Message)
		}
		if len(res.Choices) == 0 {
			return "", fmt.Errorf("empty DeepSeek response")
		}
		return res.Choices[0].Message.Content, nil
	}, 2)
}

func doRequestWithRetry(client *http.Client, req *http.Request, parse func(io.Reader) (string, error), maxRetries int) (string, error) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	req.Body.Close()

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			wait := time.Duration(attempt*2) * time.Second
			fmt.Printf("  Rate limited. Retrying in %v (attempt %d/%d)...\n", wait, attempt, maxRetries)
			time.Sleep(wait)
		}
		r, err := http.NewRequest(req.Method, req.URL.String(), bytes.NewReader(bodyBytes))
		if err != nil {
			return "", err
		}
		r.Header = req.Header.Clone()

		resp, err := client.Do(r)
		if err != nil {
			return "", fmt.Errorf("API request failed: %w", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests && attempt < maxRetries {
			lastErr = extractErrorBody(resp.StatusCode, resp.Body)
			resp.Body.Close()
			continue
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", extractErrorBody(resp.StatusCode, resp.Body)
		}
		return parse(resp.Body)
	}
	return "", fmt.Errorf("API request failed after %d retries (last error: %v)", maxRetries, lastErr)
}

func extractCode(text string) string {
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```") {
		lines := strings.SplitN(text, "\n", 2)
		if len(lines) > 1 {
			text = lines[1]
		}
		if idx := strings.LastIndex(text, "```"); idx >= 0 {
			text = text[:idx]
		}
	}
	return strings.TrimSpace(text)
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}

func isValidAICode(code string) bool {
	trimmed := strings.TrimSpace(code)
	// Must start with window.AISim assignment
	if strings.HasPrefix(trimmed, "window.AISim") || strings.HasPrefix(trimmed, "const AISim") || strings.HasPrefix(trimmed, "var AISim") || strings.HasPrefix(trimmed, "let AISim") {
		return true
	}
	// Also accept if it has AISim.run as evidence of a method
	if strings.Contains(trimmed, "AISim.run") || strings.Contains(trimmed, "AISim =") {
		return true
	}
	return false
}



type codeGenFeature struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Params      []codeGenParam   `json:"params,omitempty"`
	Actions     []codeGenAction  `json:"actions,omitempty"`
}

type codeGenParam struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Required bool        `json:"required"`
	Default  interface{} `json:"default,omitempty"`
}

type codeGenAction struct {
	Statement string           `json:"statement"`
	Condition *codeGenCond     `json:"condition,omitempty"`
}

type codeGenCond struct {
	Type       string `json:"type"`
	Expression string `json:"expression"`
}

func buildMinSpec(spec *ast.AppSpec) []byte {
	min := struct {
		Features []codeGenFeature `json:"features"`
	}{
		Features: make([]codeGenFeature, len(spec.Features)),
	}
	for i, f := range spec.Features {
		mf := codeGenFeature{Name: f.Name, Description: f.Description}
		for _, p := range f.Params {
			mf.Params = append(mf.Params, codeGenParam{
				Name: p.Name, Type: p.Type, Required: p.Required, Default: p.Default,
			})
		}
		for _, a := range f.Actions {
			ma := codeGenAction{Statement: a.Statement}
			if a.Condition != nil {
				ma.Condition = &codeGenCond{Type: a.Condition.Type, Expression: a.Condition.Expression}
			}
			mf.Actions = append(mf.Actions, ma)
		}
		min.Features[i] = mf
	}
	b, _ := json.Marshal(min)
	return b
}

func buildCodePrompt(spec *ast.AppSpec) string {
	specJSON := buildMinSpec(spec)
	return `You are a code generator. Your task is to generate ONLY valid JavaScript code.
NO explanation, NO markdown fences, NO comments outside the code.

Generate window.AISim with a run(featureName, params) method that returns an array of result objects, one per action in that feature, in order. Each result has:
  { action: "the statement text", status: "ran"|"skipped", detail: "human-readable explanation" }

Implement each feature's behavior realistically:
- For "ran" actions, include realistic output (generated IDs, timestamps, computed values)
- For "skipped" actions, explain why (condition not met, validation failed, etc.)
- ALWAYS use the params object to evaluate conditions and drive action logic
- Parse natural language conditions using params (e.g. "Title is not empty" → params.Title)
- If a param key doesn't match the spec, try common variants (snake_case, camelCase)

RULE: Return ONLY raw JavaScript. No markdown. No explanation. No text before or after. Start with: window.AISim =

Spec:
` + string(specJSON)
}

var tmpl = template.Must(template.New("sim").Parse(page))

const page = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Butter Sim — {{.AppName}}</title>
<style>
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{--bg:#0a0c10;--surface:#111318;--surface2:#181b20;--border:#22262c;--border2:#2a2e35;--text:#c1c6cc;--text2:#6b7280;--text3:#9ca3af;--accent:#6366f1;--accent2:#818cf8;--accent-glow:rgba(99,102,241,.25);--green:#22c55e;--green-bg:rgba(34,197,94,.1);--green-border:rgba(34,197,94,.25);--red:#ef4444;--red-bg:rgba(239,68,68,.1);--red-border:rgba(239,68,68,.25);--amber:#f59e0b;--amber-bg:rgba(245,158,11,.1);--amber-border:rgba(245,158,11,.25);--radius:10px;--radius-sm:6px;--shadow:0 1px 3px rgba(0,0,0,.3),0 1px 2px rgba(0,0,0,.2)}
body{background:var(--bg);color:var(--text);font-family:system-ui,-apple-system,'Segoe UI',Roboto,sans-serif;min-height:100vh;display:flex;flex-direction:column;line-height:1.5;-webkit-font-smoothing:antialiased}
header{background:linear-gradient(135deg,#1e1b4b,#312e81);padding:14px 28px;display:flex;align-items:center;gap:16px;flex-shrink:0;border-bottom:1px solid var(--border)}
header h1{font-size:17px;font-weight:600;color:#e0e7ff;letter-spacing:-.01em}
header .badge{font-size:11px;font-weight:500;color:#a5b4fc;background:rgba(255,255,255,.08);padding:3px 10px;border-radius:20px;margin-left:auto;border:1px solid rgba(165,180,252,.15)}
.container{display:flex;flex:1;overflow:hidden}
aside{width:220px;background:var(--surface);border-right:1px solid var(--border);padding:20px 0;flex-shrink:0;overflow-y:auto}
aside h2{font-size:10px;font-weight:600;text-transform:uppercase;letter-spacing:.08em;color:var(--text2);padding:0 20px 14px}
aside ul{list-style:none}
aside li{padding:9px 20px;cursor:pointer;font-size:13px;color:var(--text3);border-left:2px solid transparent;transition:all .2s;word-break:break-word;position:relative}
aside li:hover{background:var(--surface2);color:var(--text)}
aside li.active{background:var(--surface2);border-left-color:var(--accent);color:var(--accent2);font-weight:500}
main{flex:1;padding:28px 36px;overflow-y:auto}
#feature-panel{max-width:800px;animation:mainFade .35s ease}
@keyframes mainFade{from{opacity:0;transform:translateY(8px)}to{opacity:1;transform:translateY(0)}}
#feature-title{font-size:22px;font-weight:600;color:#e0e7ff;margin-bottom:6px;letter-spacing:-.02em}
#feature-desc{font-size:14px;color:var(--text3);margin-bottom:24px;line-height:1.6}
.no-params,.no-actions{font-size:13px;color:var(--text2);font-style:italic;padding:16px 0}
.params-card{background:var(--surface);border:1px solid var(--border);border-radius:var(--radius);padding:20px;margin-bottom:20px;box-shadow:var(--shadow)}
.params-card h3{font-size:10px;font-weight:600;text-transform:uppercase;letter-spacing:.08em;color:var(--text2);margin-bottom:16px}
.param-row{display:flex;align-items:center;gap:14px;margin-bottom:12px;padding:6px 0}
.param-row:last-child{margin-bottom:0}
.param-row label{width:140px;font-size:13px;font-weight:500;color:var(--text);flex-shrink:0;word-break:break-word}
.param-row label.required::after{content:" *";color:var(--red);font-weight:600}
.param-row input[type="text"],.param-row input[type="number"],.param-row select{flex:1;padding:8px 12px;background:var(--bg);border:1px solid var(--border2);border-radius:var(--radius-sm);color:var(--text);font-family:inherit;font-size:13px;outline:none;transition:border-color .2s,box-shadow .2s}
.param-row input[type="text"]:focus,.param-row input[type="number"]:focus,.param-row select:focus{border-color:var(--accent);box-shadow:0 0 0 3px var(--accent-glow)}
.param-row select{cursor:pointer;-webkit-appearance:none;appearance:none;background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' fill='%236b7280' viewBox='0 0 16 16'%3E%3Cpath d='M8 11L3 6h10z'/%3E%3C/svg%3E");background-repeat:no-repeat;background-position:right 10px center;padding-right:32px}
.param-row select option{background:var(--surface);color:var(--text)}
.toggle-wrap{display:inline-flex;align-items:center;gap:10px;cursor:pointer}
.toggle-track{position:relative;width:38px;height:22px;background:var(--border2);border-radius:11px;transition:background .25s;flex-shrink:0}
.toggle-track::after{content:'';position:absolute;top:2px;left:2px;width:18px;height:18px;background:#fff;border-radius:50%;transition:transform .25s cubic-bezier(.4,0,.2,1),background .25s;box-shadow:0 1px 3px rgba(0,0,0,.3)}
.toggle-track.on{background:var(--accent)}
.toggle-track.on::after{transform:translateX(16px)}
.toggle-label{font-size:13px;color:var(--text3);user-select:none}
#run-btn{display:inline-flex;align-items:center;gap:8px;padding:10px 28px;background:linear-gradient(135deg,var(--accent),#4f46e5);color:#fff;border:none;border-radius:var(--radius-sm);font-family:inherit;font-size:14px;font-weight:500;cursor:pointer;transition:box-shadow .2s,transform .15s;margin-bottom:24px;box-shadow:0 1px 3px rgba(99,102,241,.3)}
#run-btn:hover{box-shadow:0 4px 12px rgba(99,102,241,.45)}
#run-btn:active{transform:scale(.97)}
#flow-panel{min-height:40px;position:relative;padding-left:28px}
#flow-panel::before{content:'';position:absolute;left:11px;top:8px;bottom:8px;width:2px;background:var(--border2);border-radius:1px}
.action-card{display:flex;align-items:flex-start;gap:12px;padding:12px 16px;margin-bottom:8px;border-radius:var(--radius-sm);opacity:0;transform:translateY(10px);animation:actionSlide .4s cubic-bezier(.16,1,.3,1) forwards;animation-delay:calc(var(--i,0) * .1s);font-size:13px;line-height:1.5;position:relative}
.action-card::before{content:'';position:absolute;left:-23px;top:15px;width:10px;height:10px;border-radius:50%;border:2px solid var(--border2);background:var(--bg);transition:all .3s}
@keyframes actionSlide{to{opacity:1;transform:translateY(0)}}
.action-card.will-run{background:var(--green-bg);border:1px solid var(--green-border)}
.action-card.will-run::before{border-color:var(--green);background:var(--green);box-shadow:0 0 6px rgba(34,197,94,.4)}
.action-card.skipped{background:var(--red-bg);border:1px solid var(--red-border);opacity:0;animation-name:actionSlideSkipped}
.action-card.skipped::before{border-color:var(--red);background:transparent}
@keyframes actionSlideSkipped{to{opacity:.7;transform:translateY(0)}}
.action-card.unknown{background:var(--amber-bg);border:1px solid var(--amber-border)}
.action-card.unknown::before{border-color:var(--amber);background:transparent}
.action-icon{font-size:14px;width:22px;text-align:center;flex-shrink:0;margin-top:1px}
.action-card.will-run .action-icon{animation:iconPop .3s cubic-bezier(.34,1.56,.64,1) calc(var(--i,0)*.1s + .2s) both}
@keyframes iconPop{0%{transform:scale(0)}60%{transform:scale(1.3)}to{transform:scale(1)}}
.action-text{flex:1;word-break:break-word}
.action-reason{font-size:11px;color:var(--text2);flex-shrink:0;max-width:320px;text-align:right;word-break:break-word;font-style:italic}
details#code-panel{margin-top:24px;border:1px solid var(--border);border-radius:var(--radius);background:var(--surface);overflow:hidden;box-shadow:var(--shadow)}
details#code-panel summary{padding:12px 18px;cursor:pointer;font-size:13px;font-weight:500;color:var(--text3);user-select:none;transition:color .2s;list-style:none;-webkit-list-style:none}
details#code-panel summary::-webkit-details-marker{display:none}
details#code-panel summary:hover{color:var(--text)}
details#code-panel pre{background:var(--bg);padding:18px;overflow-x:auto;border-top:1px solid var(--border);max-height:500px;overflow-y:auto;margin:0}
details#code-panel pre{background:var(--bg);padding:18px;overflow-x:auto;border-top:1px solid var(--border);max-height:500px;overflow-y:auto;margin:0}
details#code-panel code{font-family:'JetBrains Mono','Fira Code','Cascadia Code',monospace;font-size:13px;line-height:1.7;color:#e2e8f0;white-space:pre}
::-webkit-scrollbar{width:8px;height:8px}
::-webkit-scrollbar-track{background:transparent}
::-webkit-scrollbar-thumb{background:var(--border2);border-radius:4px}
::-webkit-scrollbar-thumb:hover{background:var(--text2)}
@media(max-width:700px){aside{display:none}main{padding:20px}.param-row{flex-direction:column;align-items:stretch;gap:6px}.param-row label{width:auto}}
</style>
</head>
<body>
<header>
  <h1>Butter Sim — {{.AppName}}</h1>
  <span class="badge">Generated via {{.Provider}}</span>
</header>
<div class="container">
  <aside>
    <h2>Features</h2>
    <ul id="feature-list"></ul>
  </aside>
  <main>
    <div id="feature-panel">
      <h2 id="feature-title"></h2>
      <p id="feature-desc"></p>
      <div id="params-panel"></div>
      <button id="run-btn">&#9654; Run Simulation</button>
      <div id="flow-panel"></div>
      {{if .HasCode}}
      <details id="code-panel">
        <summary><span id="code-arrow">▶</span> AI Source Code</summary>
        <pre><code id="code-content">{{.GeneratedCode}}</code></pre>
      </details>
      <style>
      #code-arrow{display:inline-block;transition:transform .25s;font-size:10px;color:var(--text2);margin-right:8px}
      #code-panel[open] #code-arrow{transform:rotate(90deg)}
      </style>
      {{end}}
    </div>
  </main>
</div>
{{if .HasCode}}<script>{{.AIExecutable}}</script>{{end}}
<script>
const SPEC = {{.Spec}};

function evalExpr(expr, params, paramNames) {
  var js = expr;
  var hasParam = false;
  paramNames.forEach(function(name){
    if (expr.indexOf(name) !== -1) hasParam = true;
    var re = new RegExp('\\b' + name.replace(/[.*+?^${}()|[\]\\]/g,'\\$&') + '\\s+is\\s+not\\s+empty\\b');
    js = js.replace(re, '(p["' + name + '"]!==""&&p["' + name + '"]!=null)');
    var re2 = new RegExp('\\b' + name.replace(/[.*+?^${}()|[\]\\]/g,'\\$&') + '\\s+is\\s+empty\\b');
    js = js.replace(re2, '(p["' + name + '"]===""||p["' + name + '"]==null)');
    var re3 = new RegExp('\\b' + name.replace(/[.*+?^${}()|[\]\\]/g,'\\$&') + '\\s*(==|!=|>=|<=|>|<)\\s*("[^"]*"|\'[^\']*\'|[\\w.]+)', 'g');
    js = js.replace(re3, function(_,op,val){return 'p["' + name + '"]' + op + val;});
    var re4 = new RegExp('\\b' + name.replace(/[.*+?^${}()|[\]\\]/g,'\\$&') + '\\b(?!\\s*(==|!=|>=|<=|>|<|\\s+is\\b))', 'g');
    js = js.replace(re4, 'p["' + name + '"]');
  });
  if (!hasParam) return true;
  try { return new Function('p','return ('+js+')')(params); }
  catch(e) { return null; }
}

(function(){
let currentFeature = null;

function renderFeatures() {
  const list = document.getElementById('feature-list');
  (SPEC.features||[]).forEach(function(f,i){
    const li = document.createElement('li');
    li.textContent = f.name;
    li.dataset.index = i;
    li.addEventListener('click', function(){selectFeature(i);});
    list.appendChild(li);
  });
  if ((SPEC.features||[]).length > 0) selectFeature(0);
}

function selectFeature(index) {
  currentFeature = index;
  const feature = SPEC.features[index];
  document.querySelectorAll('#feature-list li').forEach(function(li,i){
    li.classList.toggle('active', i === index);
  });
  document.getElementById('feature-title').textContent = feature.name;
  document.getElementById('feature-desc').textContent = feature.description || '';
  renderParams(feature);
  var fp = document.getElementById('flow-panel');
  fp.innerHTML = '';
}

function paramId(name){return 'p-'+name;}

function createToggle(name, checked, labelEl) {
  var wrap = document.createElement('span');
  wrap.className = 'toggle-wrap';
  var track = document.createElement('span');
  track.className = 'toggle-track' + (checked?' on':'');
  track.setAttribute('data-param', name);
  var lbl = document.createElement('span');
  lbl.className = 'toggle-label';
  lbl.textContent = checked ? 'On' : 'Off';
  wrap.appendChild(track);
  wrap.appendChild(lbl);
  track.addEventListener('click', function(){
    var isOn = track.classList.toggle('on');
    lbl.textContent = isOn ? 'On' : 'Off';
  });
  labelEl.appendChild(wrap);
}

function renderParams(feature) {
  var panel = document.getElementById('params-panel');
  panel.innerHTML = '';
  if (!feature.params || feature.params.length === 0) {
    panel.innerHTML = '<div class="no-params">No parameters</div>';
    return;
  }
  var card = document.createElement('div');
  card.className = 'params-card';
  var h3 = document.createElement('h3');
  h3.textContent = 'Parameters';
  card.appendChild(h3);
  feature.params.forEach(function(p){
    var row = document.createElement('div');
    row.className = 'param-row';
    var label = document.createElement('label');
    label.textContent = p.name;
    label.htmlFor = paramId(p.name);
    if (p.required) label.classList.add('required');
    var input;
    if (p.type && p.type.startsWith('enum')) {
      input = document.createElement('select');
      input.id = paramId(p.name);
      var m = p.type.match(/\[(.+?)\]/);
      if (m) {
        m[1].split(',').map(function(v){return v.trim().replace(/"/g,'');}).forEach(function(v){
          var opt = document.createElement('option');
          opt.value = v; opt.textContent = v;
          input.appendChild(opt);
        });
      }
      if (p.default !== undefined && p.default !== null) input.value = String(p.default).replace(/^"|"$/g,'');
      row.appendChild(label);
      row.appendChild(input);
    } else if (p.type === 'bool') {
      createToggle(paramId(p.name), p.default === true || p.default === 'true', label);
      row.appendChild(label);
    } else if (p.type === 'int' || p.type === 'float') {
      input = document.createElement('input');
      input.type = 'number'; input.id = paramId(p.name);
      input.step = p.type === 'int' ? '1' : 'any';
      if (p.default !== undefined && p.default !== null) input.value = p.default;
      row.appendChild(label);
      row.appendChild(input);
    } else {
      input = document.createElement('input');
      input.type = 'text'; input.id = paramId(p.name);
      if (p.default !== undefined && p.default !== null) input.value = String(p.default).replace(/^"|"$/g,'');
      row.appendChild(label);
      row.appendChild(input);
    }
    card.appendChild(row);
  });
  panel.appendChild(card);
}

function gatherParams(feature) {
  var params = {};
  (feature.params||[]).forEach(function(p){
    var id = paramId(p.name);
    if (p.type === 'bool') {
      var track = document.querySelector('.toggle-track[data-param="'+id+'"]');
      params[p.name] = track ? track.classList.contains('on') : false;
    } else {
      var el = document.getElementById(id);
      if (!el) return;
      if (p.type === 'int') {
        params[p.name] = el.value !== '' ? parseInt(el.value, 10) : '';
      } else if (p.type === 'float') {
        params[p.name] = el.value !== '' ? parseFloat(el.value) : '';
      } else {
        params[p.name] = el.value;
      }
    }
  });
  return params;
}

function runSimulation() {
  if (currentFeature === null) return;
  var feature = SPEC.features[currentFeature];
  var params = gatherParams(feature);
  var panel = document.getElementById('flow-panel');
  panel.innerHTML = '';
  if (!feature.actions || feature.actions.length === 0) {
    panel.innerHTML = '<div class="no-actions">No actions defined</div>';
    return;
  }

  if (typeof AISim !== 'undefined' && AISim && AISim.run) {
    var aiResults;
    try { aiResults = AISim.run(feature.name, params); } catch(e) { aiResults = []; }
    feature.actions.forEach(function(a,i){
      var card = document.createElement('div');
      card.className = 'action-card';
      card.style.setProperty('--i', i);
      var icon = document.createElement('span');
      icon.className = 'action-icon';
      var text = document.createElement('span');
      text.className = 'action-text';
      text.textContent = a.statement;
      var ar = aiResults[i];
      if (ar && ar.status === 'ran') {
        icon.textContent = '✅';
        card.classList.add('will-run');
        var rsn = document.createElement('span');
        rsn.className = 'action-reason';
        rsn.textContent = '← ' + (ar.detail || 'Done');
        card.appendChild(icon);
        card.appendChild(text);
        card.appendChild(rsn);
      } else if (ar && ar.status === 'skipped') {
        icon.textContent = '⏭️';
        card.classList.add('skipped');
        var rsn = document.createElement('span');
        rsn.className = 'action-reason';
        rsn.textContent = '← ' + (ar.detail || 'Condition not met');
        card.appendChild(icon);
        card.appendChild(text);
        card.appendChild(rsn);
      } else {
        icon.textContent = '❓';
        card.classList.add('unknown');
        card.appendChild(icon);
        card.appendChild(text);
      }
      panel.appendChild(card);
    });
  } else {
    var paramNames = Object.keys(params);
    feature.actions.forEach(function(a,i){
      var card = document.createElement('div');
      card.className = 'action-card';
      card.style.setProperty('--i', i);
      var icon = document.createElement('span');
      icon.className = 'action-icon';
      var text = document.createElement('span');
      text.className = 'action-text';
      text.textContent = a.statement;
      var willRun = true;
      var reason = '';
      if (a.condition) {
        var r = evalExpr(a.condition.expression, params, paramNames);
        if (r === undefined || r === null) {
          icon.textContent = '❓';
          card.classList.add('unknown');
          willRun = false;
        } else if (a.condition.type === 'unless') {
          willRun = !r;
          icon.textContent = willRun ? '✅' : '❌';
          card.classList.add(willRun ? 'will-run' : 'skipped');
          if (!willRun) reason = a.condition.type + ' "' + a.condition.expression + '" → false';
        } else {
          willRun = !!r;
          icon.textContent = willRun ? '✅' : '❌';
          card.classList.add(willRun ? 'will-run' : 'skipped');
          if (!willRun) reason = a.condition.type + ' "' + a.condition.expression + '" → false';
        }
      } else {
        icon.textContent = '✅';
        card.classList.add('will-run');
      }
      card.appendChild(icon);
      card.appendChild(text);
      if (reason) {
        var rsn = document.createElement('span');
        rsn.className = 'action-reason';
        rsn.textContent = '← ' + reason;
        card.appendChild(rsn);
      }
      panel.appendChild(card);
    });
  }
  setTimeout(function(){
    var last = panel.lastElementChild;
    if (last) last.scrollIntoView({behavior:'smooth',block:'nearest'});
  }, feature.actions.length * 100 + 200);
}

document.getElementById('run-btn').addEventListener('click', runSimulation);
renderFeatures();
})();
</script>
</body>
</html>`
