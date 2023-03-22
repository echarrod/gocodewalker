package main

import (
	"fmt"
	"github.com/boyter/gocodewalker"
	"os"
	"strings"
)

// Proper test designed to confirm that .gitignores work as expected with globs
// Designed to work against https://github.com/svent/gitignore-test
// If you compile and run this it should produce the same output as the following tools
// when run from the directory you check it out into
//
// rg ^foo: | sort
// git grep ^foo: | sort
// gocodewalker | sort
func main() {
	// useful for improving performance, then go tool pprof -http=localhost:8090 profile.pprof
	//f, _ := os.Create("profile.pprof")
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	fileListQueue := make(chan *gocodewalker.File, 10_000)
	fileWalker := gocodewalker.NewFileWalker(".", fileListQueue)

	// handle the errors by printing them out and then ignore
	errorHandler := func(e error) bool {
		fmt.Println("ERR", e.Error())
		return true
	}
	fileWalker.SetErrorHandler(errorHandler)

	go func() {
		err := fileWalker.Start()
		if err != nil {
			fmt.Println("ERR", err.Error())
		}
	}()

	for f := range fileListQueue {
		contents, _ := os.ReadFile(f.Location)
		if len(contents) > 10 {
			contents = contents[:10]
		}
		fmt.Printf("%v:%v\n", f.Location, strings.TrimSpace(string(contents)))
	}
}
