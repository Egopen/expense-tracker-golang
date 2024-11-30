package main

import (
	routing "main/Routing"
)

func main() {
	r := routing.SetupRouter()
	r.Run()
}
