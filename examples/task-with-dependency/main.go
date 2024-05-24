package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ignoxx/gogent"
)

func main() {
	jokesWriter := gogent.NewGopher().
		WithRole("Professional Jokes Writer").
		WithGoal("Writing exceptional funny and not well-known jokes about programmers").
		WithBackstory("You work at a leading tech think tank. Your expertise lies in identifying trending jokes")

	casualDev := gogent.NewGopher().
		WithRole("Casual JS Developer").
		WithGoal("Rate the jokes").
		WithBackstory("You have been watching memes, fun reels, attenting comedy clubs and more since 20 years. You know what is funny and what not!")

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

	task2 := gogent.NewTask().
		WithLLM(gpt35).
		WithGopher(casualDev).
		WithDependencies(task1).
		WithCreativity(gogent.SlightlyCreative).
		WithDescription("Rate each joke between 1-5 (1=joke is bad, 5=joke is really funny) and tell the reason why you rated it how you rated").
		WithExpectedOutput(`return a valid JSON array, and put each joke in its own object containing the joke itself, rating and the reason like: [{"joke": "...", rating: 1, reason: "..."}, {...}, ..]`)

	ctx := context.Background()
	if err := task2.Process(ctx); err != nil {
		slog.Error("task failed", slog.String("description", task2.Description), slog.String("err", err.Error()))
	}

	// write task2 output to a file
	file, err := os.Create("task_output.json")
	if err != nil {
		slog.Error("failed to create file", slog.String("err", err.Error()))
	}
	defer file.Close()
	file.ReadFrom(task2.Reader())
}
