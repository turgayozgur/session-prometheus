package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randomlyError(input string) error {
	if rand.Intn(100) < 20 {
		return fmt.Errorf("an error occurred during the payment for %s bank", input)
	}
	return nil
}

func randomlyWait(min int) {
	r := rand.Intn(min * 10 - min) + min
	time.Sleep(time.Duration(r) * time.Millisecond)
}