package main

import (
	"context"
	"fmt"
	"testing"
)

func TestPipeline1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	intStream := generator(ctx, 1, 2, 3, 4)
	pipeline := multiply(ctx, add(ctx, multiply(ctx, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}
