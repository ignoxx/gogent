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
casualDev := gogent.NewGopher().
    WithRole("Casual JS Developer").
    WithGoal("Rate the jokes").
    WithBackstory("You have been watching memes, fun reels, attenting comedy clubs and more since 20 years. You know what is funny and what not!")

gpt35 := gogent.NewLLM(
    gogent.WithProvider(gogent.ProviderOpenAI),
    gogent.WithModel(gogent.ModelOpenAIGpt35),
    gogent.WithApiKey("YOUR-OPENAI-KEY"),
)

task := gogent.NewTask().
    WithLLM(gpt35).
    WithDescription("Write 10 jokes").
    WithExpectedOutput(`All jokes must be returned in a single valid JSON array like: ["joke1", "joke2", "joke3", ...]`).
    WithGopher(jokesWriter)

ctx := context.Background()
if err := task.Process(ctx); err != nil {
    slog.Error("task failed", slog.String("description", task.Description), slog.String("err", err.Error()))
}

fmt.Println(task.Output)
```

Check `examples/` for more.
