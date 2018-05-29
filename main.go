package main

import (
	"fmt"
	"time"
)

func main() {
	unix := time.Now().Unix() & 0xFFFFFFFF
	fmt.Printf("%x\n", unix)
}
