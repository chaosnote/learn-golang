package main

import "idv/chris/utils"

func main() {
	console := utils.NewConsoleLogger(1)

	console.Debug(utils.LogFields{
		"test": "123",
	})

	w := utils.NewFileLogger("./dist", 1)
	w.Info(map[string]interface{}{
		"name": "this is test",
	})
}
