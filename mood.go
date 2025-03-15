package mood

import (
	"fmt"
	"os"
)

type AppArg struct {
	Position int
	Value    string
}

type CommandFactory func(app *App) (Cmd, error)

type App struct {
	Source         string
	Args           map[string]AppArg
	Flags          map[string]AppArg
	Commands       map[string]CommandFactory
	DefaultFactory CommandFactory
	Default        Cmd
	Store          map[string]any
}

type Cmd interface {
	Execute(app *App) error
}

func New() App {
	osArgs := os.Args
	flags := make(map[string]AppArg)
	args := make(map[string]AppArg)

	source := ""
	if len(osArgs) > 0 {
		source = osArgs[0]
	}

	for i, arg := range osArgs {
		if len(arg) > 1 && i > 0 && (arg[0] == '-' || (len(arg) > 2 && arg[:2] == "--")) {
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
		Source:         source,
		Args:           args,
		Flags:          flags,
		Commands:       make(map[string]CommandFactory),
		DefaultFactory: nil,
		Default:        defaultCmd{},
		Store:          make(map[string]any),
	}
}

func (app *App) At(commandName string, factory CommandFactory) {
	app.Commands[commandName] = factory
}

func (app *App) SetDefault(factory CommandFactory) {
	app.DefaultFactory = factory
	cmd, err := factory(app)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error setting default command: %v\n", err)
		return
	}
	app.Default = cmd
}

func (app *App) Run() error {
	firstArgPosition := 1
	var firstArg string

	for _, arg := range app.Args {
		if arg.Position == firstArgPosition {
			firstArg = arg.Value
			break
		}
	}

	if firstArg == "" {
		return app.Default.Execute(app)
	}

	if factory, exists := app.Commands[firstArg]; exists {
		cmd, err := factory(app)
		if err != nil {
			return fmt.Errorf("error initializing command '%s': %w", firstArg, err)
		}
		return cmd.Execute(app)
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
