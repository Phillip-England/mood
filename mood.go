package mood

import (
	"fmt"
	"os"
)

type AppArg struct {
	Position int
	Value    string
}

type App struct {
	OriginalArgs []string
	Source       string
	Args         map[string]AppArg
	Flags        map[string]AppArg
	Commands     map[string]Cmd
	Default      Cmd
	Store        map[string]any
}

type Cmd interface {
	Execute(app *App) error
}

func New() App {
	ogArgs := os.Args
	flags := make(map[string]AppArg)
	args := make(map[string]AppArg)

	source := ""
	if len(ogArgs) > 0 {
		source = ogArgs[0]
	}

	for i, arg := range ogArgs {
		if i == 0 {
			continue
		}
		if len(arg) > 1 && (arg[0] == '-' || (len(arg) > 2 && arg[:2] == "--")) {
			flags[arg] = AppArg{
				Position: i,
				Value:    arg,
			}
			continue
		}
		args[arg] = AppArg{
			Position: i,
			Value:    arg,
		}
	}

	return App{
		OriginalArgs: ogArgs,
		Source:       source,
		Args:         args,
		Flags:        flags,
		Commands:     make(map[string]Cmd),
		Default:      defaultCmd{},
		Store:        make(map[string]any),
	}
}

// Now takes app as a parameter in the function signature
func (app *App) At(commandName string, fn func(app *App) (Cmd, error)) {
	cmd, err := fn(app)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error registering command '%s': %v\n", commandName, err)
		return
	}
	app.Commands[commandName] = cmd
}

// Now takes app as a parameter in the function signature
func (app *App) SetDefault(fn func(app *App) (Cmd, error)) {
	cmd, err := fn(app)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error setting default command: %v\n", err)
		return
	}
	app.Default = cmd
}

func (app *App) Run() error {
	if len(app.Args) == 0 {
		return app.Default.Execute(app)
	}

	for _, arg := range app.Args {
		cmd, exists := app.Commands[arg.Value]
		if exists {
			if err := cmd.Execute(app); err != nil {
				return fmt.Errorf("error executing command '%s': %w", arg.Value, err)
			}
			return nil
		}
	}

	return app.Default.Execute(app)
}

func (app *App) SetStore(key string, value any) error {
	if _, exists := app.Store[key]; exists {
		return fmt.Errorf("error: key '%s' already exists in the store", key)
	}
	app.Store[key] = value
	return nil
}

func (app *App) GetStore(key string) (any, error) {
	value, exists := app.Store[key]
	if !exists {
		return nil, fmt.Errorf("error: key '%s' not found in the store", key)
	}
	return value, nil
}

func (app *App) HasFlag(flag string) bool {
	_, exists := app.Flags[flag]
	return exists
}

func (app *App) HasArg(arg string) bool {
	_, exists := app.Args[arg]
	return exists
}

func (app *App) GetArg(arg string) (AppArg, bool) {
	val, exists := app.Args[arg]
	return val, exists
}

func (app *App) GetArgOr(arg string, defaultValue string) string {
	if val, exists := app.Args[arg]; exists {
		return val.Value
	}
	return defaultValue
}

func (app *App) GetArgByPosition(position int) (AppArg, error) {
	for _, arg := range app.Args {
		if arg.Position == position {
			return arg, nil
		}
	}
	return AppArg{}, fmt.Errorf("error: argument at position %d not found", position)
}

type defaultCmd struct{}

func (d defaultCmd) Execute(app *App) error {
	fmt.Println("Welcome to Mood! No command provided.")
	fmt.Println("📖 Read the docs at https://github.com/Phillip-England/mood")
	return nil
}
