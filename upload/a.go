package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	fmt.Printf("%x", md5.Sum([]byte("Hell")))
}
