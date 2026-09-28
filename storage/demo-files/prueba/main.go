package main

import (
	"fmt"
	"strings"
)

func greet(name string) string {
	return fmt.Sprintf("Hola, %s", strings.TrimSpace(name))
}

func main() {
	languages := []string{"Go", "PHP", "JavaScript"}

	for _, language := range languages {
		fmt.Println(greet(language))
	}
}
