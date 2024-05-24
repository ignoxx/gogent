# Gogent - Simple Go(pher) Agents

Define your Gophers (agents) once and re-use them for various AI tasks!

## Features
- Define custom Gophers (agents) with specific roles and goals
- Integrate with popular AI models like OpenAI's GPT-3.5 or GPT-4o
- Create tasks with descriptions and expected outputs
- Handle task dependencies and process tasks in a context

## Install
```bash
go get github.com/ignoxx/gogent
```

## Use
```go
package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ignoxx/gogent"
)

func main() {
	jokesWriter := gogent.NewGopher().
		WithRole("Professional Jokes Writer").
		WithGoal("Writing exceptional funny and not well-known jokes about programmers").
		WithBackstory("You work at a leading tech think tank. Your expertise lies in identifying trending jokes")

	gpt35 := gogent.NewLLM(
		gogent.WithProvider(gogent.ProviderOpenAI),
		gogent.WithModel(gogent.ModelOpenAIGpt35),
		gogent.WithApiKey("YOUR-API-KEY"),
	)

	task1 := gogent.NewTask().
		WithLLM(gpt35).
		WithGopher(jokesWriter).
		WithCreativity(gogent.ExtremelyCreative).
		WithDescription("Write 10 jokes").
		WithExpectedOutput(`All jokes must be returned in a single valid JSON array like: ["joke1", "joke2", "joke3", ...]`)

	ctx := context.Background()
	if err := task1.Process(ctx); err != nil {
		slog.Error("task failed", slog.String("description", task1.Description), slog.String("err", err.Error()))
	}

	fmt.Println(task1.Output)
}
```

Check `examples/` for more.
