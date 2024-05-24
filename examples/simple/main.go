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
		WithCreativity(gogent.SlightlyCreative).
		WithDescription("Write 10 jokes").
		WithExpectedOutput(`All jokes must be returned in a single valid JSON array like: ["joke1", "joke2", "joke3", ...]`)

	ctx := context.Background()
	if err := task1.Process(ctx); err != nil {
		slog.Error("task failed", slog.String("description", task1.Description), slog.String("err", err.Error()))
	}

	fmt.Println(task1.Output)
	fmt.Println()

	res, err := task1.OutputToJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	fmt.Println()

	resStruct := []string{}
	if err := task1.OutputToJSONStruct(&resStruct); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", resStruct)
}
