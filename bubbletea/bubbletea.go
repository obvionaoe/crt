package bubbletea

import (
	"fmt"
	"github.com/BigJk/crt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"image/color"
	"os"
	"syscall"
)

func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}

type fakeEnviron struct{}

func (f fakeEnviron) Environ() []string {
	return []string{"TERM", "COLORTERM"}
}

func (f fakeEnviron) Getenv(s string) string {
	switch s {
	case "TERM":
		return "xterm-256color"
	case "COLORTERM":
		return "truecolor"
	}
	return ""
}

func Window(width int, height int, fonts crt.Fonts, model tea.Model, defaultBg color.Color, options ...tea.ProgramOption) (*crt.Window, error) {
	gameInput := crt.NewConcurrentRW()
	gameOutput := crt.NewConcurrentRW()

	go gameInput.Run()
	go gameOutput.Run()

	prog := tea.NewProgram(
		model,
		append([]tea.ProgramOption{
			tea.WithMouseAllMotion(),
			tea.WithInput(gameInput),
			tea.WithOutput(termenv.NewOutput(gameOutput, termenv.WithEnvironment(fakeEnviron{}), termenv.WithTTY(true), termenv.WithProfile(termenv.TrueColor), termenv.WithColorCache(true))),
			tea.WithANSICompressor(),
		}, options...)...,
	)

	go func() {
		if _, err := prog.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}

		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	return crt.NewGame(width, height, fonts, gameOutput, NewBubbleTeaAdapter(prog), defaultBg)
}
