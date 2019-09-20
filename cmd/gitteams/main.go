package main

import (
	"github.com/NoUseFreak/gitteams/internal/app/gitteams"
)

func main() {
	gitteams.CreateDynamicCommands()
	gitteams.Execute()
}
