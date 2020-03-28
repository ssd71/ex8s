package main

import (
	"os"

	"github.com/ssd71/ex8s/sheetutil"
	"github.com/ssd71/ex8s/updatelistener"
)

func main() {
	updatelistener.StartListener(func(data []string) {
		sheetutil.Init(os.Getenv("SHEETID"))
		sheetutil.UpdateOrInsert(data)
	})
}
