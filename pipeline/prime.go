package main

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"strconv"
)

func RandIntStream(ctx context.Context) <-chan interface{} {
	return repeatFn(ctx, func() interface{} {
		return rand.Intn(10000000000)
	})
}

func checkPrime(arg interface{}) interface{} {
	switch v := arg.(type) {
	case int:
		cmd := exec.Command("factor", strconv.Itoa(v))
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}

		matched, err := regexp.Match("^\\d+:\\s*\\d+\\s*$", out)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		if matched == true {
			return arg
		}
	}
	return nil
}

func PrimeFinder(ctx context.Context, src <-chan interface{}) <-chan interface{} {
	return pipe(ctx, src, checkPrime)
}
