package mood

import (
	"fmt"
	"os"
)

func main() {

	app := New()

	app.At("build", func(app *Mood) error {
		fmt.Println("building...")
		return nil
	})

	err := app.Run()
	if err != nil {
		panic(err)
	}

}

type MoodArg struct {
	Position int
	Value    string
}

type Mood struct {
	OriginalArgs []string
	Args         []MoodArg
	Flags        []MoodArg
	Actions      map[string]ActionFunc
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
	}
}

func (app *Mood) At(commandName string, fn ActionFunc) {
	app.Actions[commandName] = fn
}

func (app *Mood) Run() error {
	if len(app.Args) == 0 {
		return fmt.Errorf(`welcome to mood, please pass an arg to your cli application`)
	}
	cmd := app.Args[0].Value
	fn := app.Actions[cmd]
	if fn == nil {
		return fmt.Errorf(`[%s] is not a registered command`, cmd)
	}
	err := fn(app)
	if err != nil {
		return err
	}
	return nil
}

func (app *Mood) HasFlag(flag string) bool {
	for _, f := range app.Flags {
		if f.Value == flag {
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

type ActionFunc func(app *Mood) error
