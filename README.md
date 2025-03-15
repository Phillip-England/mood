# mood
Quick, simple, but *structured*, cli apps in go.

## Installation
```bash
go get github.com/Phillip-England/mood
```

## Hello, World!
A simple cli application made using mood:

```go
package main

import (
  "fmt"
  "github.com/Phillip-England/mood"  
)

func main() {

	app := mood.New()

	app.SetDefault(NewDefaultCmd)
	app.At("help", NewHelpCmd)

	err := app.Run()
	if err != nil {
		panic(err)
	}

}

//======================
// DefaultCmd
//======================

type DefaultCmd struct{}

func NewDefaultCmd(app *mood.App) (mood.Cmd, error) {
	return DefaultCmd{}, nil
}

func (cmd DefaultCmd) Execute(app *mood.App) error {
  fmt.Println("Hello, World!")
	return nil
}

//======================
// HelpCmd
//======================

type HelpCmd struct{}

func NewHelpCmd(app *mood.App) (mood.Cmd, error) {
	return HelpCmd{}, nil
}

func (cmd HelpCmd) Execute(app *mood.App) error {
  fmt.Println("O'Doyle Rules!")
	return nil
}
```

## Cmd
Commands are the center of mood applications and allow us to generate composable, executable blocks of code.

A `Cmd` is any struct which has an execute method. It is defined like so:

```go
type Cmd interface {
	Execute(app *App) error
}
```

Here is a quick copy/paste template for making a `Cmd`:

```go
type DefaultCmd struct{}

func NewDefaultCmd(app *mood.App) (Cmd, error) {
  // creation logic (validation, data-parsing, flag-handling, ect)
	return DefaultCmd{}, nil
}

func (cmd DefaultCmd) Execute(app *App) error {
  // execution logic
	return nil
}
```

## Routing
Routes are based off of the **first** arg passed into a cli program. If no arg is passed, the default handler will be executed, stating:

```bash
Welcome to Mood! No command provided.
📖 Read the docs at https://github.com/Phillip-England/mood
```

You can override the default handler by providing a `Cmd` like so:
```go
app := mood.New()
app.SetDefault(NewDefaultCmd)
```

If a user does provide an arg, the first arg may be associated with a `Cmd`:
```go
app.At("help", NewHelpCmd)
```

With this model, the first arg provided to the cli app becomes our logic-branching mechanism.

## Store
You may want to share data between different execution branches. In such case, the store may be used to get or retrieve values using `app.GetStore()` or `app.SetStore()`.

## Flags
Any arg beginning with a dash `-`, or a double-dash `--` is considered a flag. These flags may be used to alter how an execution branch behaves.

You can view the flags using `*mood.App.Flags` or quickly check if a flag exists using `*mood.App.HasFlag()`

## Helper Methods
In every `Execute` func, an instance of `*mood.App` is piped in. `*mood.App` has some useful methods for handling command line args and flags.
