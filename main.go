package main

import (
	"context"
	"github.com/fatih/color"
	"os"
	"os/signal"

	"github.com/iyear/tdl/cmd"
)

func main() {
	//withTimeout, cancelFunc := context.WithTimeout(context.Background(), 120*time.Second)
	//defer cancelFunc()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := cmd.New().ExecuteContext(ctx); err != nil {
		color.Red("%v", err)
		os.Exit(1)
	}
}
