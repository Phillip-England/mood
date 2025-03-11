package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {

	app := NewMood()

	app.At("build", 0, func(cli *Mood) error {
		fmt.Println("hit the build route!")
		return nil
	})

	app.At("bundle", 0, func(cli *Mood) error {
		fmt.Println("hit the bundle route!")
		return nil
	})

	err := app.Run()
	if err != nil {
		fmt.Println(err.Error())
	}

}

type Action struct {
	Func     ActionFunc
	Priority int
}

type ActionFunc func(app *Mood) error

type Mood struct {
	OriginalArgs []string
	Args         []string
	Destination  string
	Commands     map[string]Command
	Flags        []string
}

func NewMood() *Mood {
	originalArgs := os.Args
	destination := originalArgs[0]
	var args []string
	var flags []string
	for i, arg := range originalArgs {
		if i == 0 {
			continue
		}
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			continue
		}
		args = append(args, arg)
	}
	return &Mood{
		OriginalArgs: originalArgs,
		Destination:  destination,
		Commands:     make(map[string]Command),
		Args:         args,
		Flags:        flags,
	}
}

func (app *Mood) At(arg string, priority int, actionFunc ActionFunc) {
	action := Action{
		Func:     actionFunc,
		Priority: priority,
	}
	command := Command{
		Arg:    arg,
		Action: action,
	}
	app.Commands[arg] = command
}

func (app *Mood) HasFlag(flag string) bool {
	return slices.Contains(app.Flags, flag)
}

func (app *Mood) Run() error {
	var actions []Action

	for _, arg := range app.Args {
		command, exists := app.Commands[arg]
		if !exists {
			return fmt.Errorf(`[%s] is not associated with an operation`, arg)
		}

		if command.Action.Func == nil {
			return fmt.Errorf(`[%s] is not associated with a valid function`, arg)
		}

		actions = append(actions, command.Action)
	}

	slices.SortFunc(actions, func(a, b Action) int {
		return a.Priority - b.Priority
	})

	for _, action := range actions {
		if err := action.Func(app); err != nil {
			return err
		}
	}

	return nil
}

type Command struct {
	Arg    string
	Action Action
}
