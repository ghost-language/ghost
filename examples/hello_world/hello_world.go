package main

import "ghostlang.org/x/ghost/ghost"

func main() {
	ghost.NewScript(`print("hello world!")`)

	ghost.Evaluate()
}
