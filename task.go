package gogent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

// Also known as Temperature
type Creativity float32

const (
	NotCreative       Creativity = 0.0
	SlightlyCreative  Creativity = 0.5
	SomewhatCreative  Creativity = 1.0
	VeryCreative      Creativity = 1.5
	ExtremelyCreative Creativity = 2.0
)

type Task struct {
	Description    string
	ExpectedOutput string
	Output         string
	Creativity     float32
	Dependencies   []*Task
	Gopher         *Gopher
	LLM            *LLM
}

func NewTask() *Task {
	return &Task{}
}

func (t *Task) WithLLM(llm *LLM) *Task {
	t.LLM = llm
	return t
}

func (t *Task) WithCreativity(c Creativity) *Task {
	t.Creativity = float32(c)
	return t
}

func (t *Task) WithDescription(description string) *Task {
	t.Description = description
	return t
}

func (t *Task) WithExpectedOutput(expectedOutput string) *Task {
	t.ExpectedOutput = expectedOutput
	return t
}

func (t *Task) WithGopher(gopher *Gopher) *Task {
	t.Gopher = gopher
	return t
}

func (t *Task) WithDependencies(dependencies ...*Task) *Task {
	t.Dependencies = dependencies
	return t
}

func (t *Task) CanProcess() bool {
	for _, dependency := range t.Dependencies {
		if !dependency.IsDone() {
			return false
		}
	}

	return true
}

func (t *Task) IsDone() bool {
	if t.Output == "" {
		return false
	}

	return true
}

func (t *Task) Reader() io.Reader {
	return strings.NewReader(t.Output)
}

// Tries to marshal the output to a json string
func (t *Task) OutputToJSON() (string, error) {
	json, err := json.Marshal(t.Output)
	if err != nil {
		return "", err
	}

	return string(json), nil
}

// `out` must be a struct with json tags matching the ExpectedOutput
func (t *Task) OutputToJSONStruct(out any) error {
	err := json.Unmarshal([]byte(t.Output), out)
	if err != nil {
		return err
	}

	return nil
}

func (t *Task) Process(ctx context.Context) error {

	if len(t.Dependencies) > 0 {
		slog.Info("Processing Task dependencies first", slog.Int("amount", len(t.Dependencies)), slog.String("description", t.Description))
		for _, dep := range t.Dependencies {
			if !dep.IsDone() {
				slog.Info("Processing dependency", slog.String("description", dep.Description))
				err := dep.Process(ctx)
				if err != nil {
					return err
				}
				slog.Info("Dependency processed", slog.String("description", dep.Description))
			}
		}
	}

	slog.Info("Processing Task", slog.String("description", t.Description))

	client := t.LLM.Client()
	client.WithSystemPrompt(t.systemPrompt())
	client.WithUserPrompt(t.userPrompt())
	client.WithTemperature(t.Creativity)

	resp, err := client.Run(ctx)
	if err != nil {
		return err
	}

	t.Output = resp.Content
	t.LLM.TotalTokensUsed += resp.TotalTokens
	t.LLM.PromptTokensUsed += resp.PromptTokens
	t.LLM.CompletionTokensUsed += resp.CompletionTokens

	return nil
}

// Instruction
func (t *Task) systemPrompt() string {
	sb := strings.Builder{}
	sb.WriteString("Your role is '" + t.Gopher.Role + "'\n")
	sb.WriteString("With the goal: '" + t.Gopher.Goal + "'\n")
	sb.WriteString("Your Backstory: '" + t.Gopher.Backstory + "'\n")
	sb.WriteString("This is your expected way to respond, you MUST ONLY respond with this and NOTHING else!!:\n")
	sb.WriteString("'''\n")
	sb.WriteString(t.ExpectedOutput + "\n")
	sb.WriteString("'''\n")

	return sb.String()
}

// Task
func (t *Task) userPrompt() string {
	sb := strings.Builder{}

	if len(t.Dependencies) > 0 {
		sb.WriteString("------------------------\n")
		sb.WriteString("Before we get to your actual task, the following provides you more context which is relevant for you to finish your task, read it carefully!\n")
		sb.WriteString("The following information is a result of previous finished tasks which are dependencies for your task:")
		for i, dep := range t.Dependencies {
			sb.WriteString("# Finished Task Nr. " + fmt.Sprint(i) + ".\n")
			sb.WriteString("## Task Description\n")
			sb.WriteString(dep.Description + "\n")
			sb.WriteString("## Task Outcome\n")
			sb.WriteString("'''\n")
			sb.WriteString(dep.Output + "\n")
			sb.WriteString("'''\n")
			sb.WriteString("\n")
		}
		sb.WriteString("------------------------\n\n")
		sb.WriteString("When you are done carefully reading the previous task outcomes, proceed to your actual task below\n")
		sb.WriteString("\n")
	}

	sb.WriteString("With all the information you have, read and solve YOUR task carefully:\n")
	sb.WriteString("'''\n")
	sb.WriteString(t.Description + "\n")
	sb.WriteString("'''\n")
	return sb.String()
}
