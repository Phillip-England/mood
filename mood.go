package mood

import (
	"fmt"
	"os"
)

type ActionFunc func(app *Mood) error

type MoodArg struct {
	Position int
	Value    string
}

type Mood struct {
	OriginalArgs []string
	Args         []MoodArg
	Flags        []MoodArg
	Actions      map[string]ActionFunc
	Default      ActionFunc
}

func New() Mood {
	ogArgs := os.Args
	var flags []MoodArg
	var args []MoodArg
	for i, arg := range ogArgs {
		if i == 0 {
			continue
		}
		if len(arg) > 1 && string(arg[0]) == "-" {
			flags = append(flags, MoodArg{
				Position: i,
				Value:    arg,
			})
			continue
		}
		args = append(args, MoodArg{
			Position: i,
			Value:    arg,
		})
	}
	return Mood{
		OriginalArgs: ogArgs,
		Args:         args,
		Flags:        flags,
		Actions:      make(map[string]ActionFunc),
		Default:      defaultWelcome,
	}
}

func (app *Mood) At(commandName string, fn ActionFunc) {
	app.Actions[commandName] = fn
}

func (app *Mood) SetDefault(fn ActionFunc) {
	app.Default = fn
}

func (app *Mood) Run() error {
	if len(app.Args) == 0 {
		return app.Default(app)
	}
	cmd := app.Args[0].Value
	fn, exists := app.Actions[cmd]
	if !exists {
		return app.Default(app)
	}
	return fn(app)
}

func (app *Mood) HasFlag(flag string) bool {
	for _, f := range app.Flags {
		if f.Value == flag {
			return true
		}
	}
	return false
}

func (app *Mood) HasArg(arg string) bool {
	for _, a := range app.Args {
		if a.Value == arg {
			return true
		}
	}
	return false
}

func (cli *Mood) EnforceArg(position int, expectedValues ...string) error {
	if position < 1 || position > len(cli.Args) {
		return fmt.Errorf("error: expected argument at position %d, but not enough arguments were provided", position)
	}
	actualValue := cli.Args[position-1].Value
	for _, expectedValue := range expectedValues {
		if actualValue == expectedValue {
			return nil
		}
	}
	return fmt.Errorf("error: expected one of %v at position %d, but got [%s]", expectedValues, position, actualValue)
}

func defaultWelcome(app *Mood) error {
	fmt.Println("Welcome to Mood! No command provided. Use --help to see available commands.")
	return nil
}

func (app *Mood) GetArg(position int) string {
	if position < 1 || position > len(app.Args) {
		return ""
	}
	return app.Args[position-1].Value
}

func (app *Mood) GetArgOr(position int, defaultValue string) string {
	if position < 1 || position > len(app.Args) {
		return defaultValue
	}
	return app.Args[position-1].Value
}
