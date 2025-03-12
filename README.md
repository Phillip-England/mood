# mood
simple cli applications in go

## Installation
```bash
go get github.com/Phillip-England/mood
```

## Hello, World!
```go
package main

import "github.com/Phillip-England/mood"

func main() {

	app := mood.New()

	app.At("build", func(app *Mood) error {
		if app.HasFlag("-f") {
      fmt.Println("do something..")
    }
    fmt.Println("building...")
		return nil
	})

	err := app.Run()
	if err != nil {
		panic(err)
	}

}
```