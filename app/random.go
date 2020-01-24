package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randomlyError(input string) {
	if rand.Intn(100) < 5 {
		panic(fmt.Sprintf("an error occurred during the payment for %s bank", input))
	}
}

func randomlyWait(min int) {
	r := rand.Intn(min*10-min) + min
	time.Sleep(time.Duration(r) * time.Millisecond)
}
