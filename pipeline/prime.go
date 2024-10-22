package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

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
