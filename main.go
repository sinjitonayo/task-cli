package main

import (
	"os"

	"github.com/sinjitonayo/task-cli-go/internal/cli"
	"github.com/sinjitonayo/task-cli-go/internal/storage"
)

func main() {
	store := storage.NewJSONStore("tasks.json")
	handler := cli.NewHandler(store)

	// os.Args: [programName, command, args...]
	handler.Run(os.Args[1:])
}
