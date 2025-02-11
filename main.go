package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/tris203/lookatmego/parse"
)

type model struct {
	presentation       *parse.Presentation
	CurrentSlide       int
	CurrenSlideSection int
	width              int
	height             int
}

func initialModel(filename string) *model {
	file, err := os.Open(filename)
	if err != nil {
		panic("Error opening file")
	}

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	presentation, err := parse.LoadFromFile(content)
	if err != nil {
		panic("Error loading file")
	}
	return &model{
		presentation:       presentation,
		CurrentSlide:       0,
		CurrenSlideSection: 0,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "h", "j", "down", "left", "enter":
			err := m.Next()
			if err != nil {
				fmt.Print("\033[H\033[2J") // Clear the screen
				return m, tea.Quit
			}
			return m, cmd
		case "l", "k", "up", "right", "backspace":
			_ = m.Prev()
			return m, cmd
		}

	}
	return m, nil
}

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return "loading..."
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap((m.width/3)*2),
	)

	markdown, _ := r.Render(m.GetCurrentSlide(m.CurrentSlide, m.CurrenSlideSection))

	leftFooter := lipgloss.NewStyle().Italic(true).AlignHorizontal(lipgloss.Left).Width(m.width / 2).PaddingLeft(m.width / 10)
	rightFooter := lipgloss.NewStyle().Italic(true).AlignHorizontal(lipgloss.Right).Width(m.width / 2).PaddingRight(m.width / 10)
	headerStyle := lipgloss.NewStyle().Bold(true)

	header := lipgloss.Place(
		m.width,
		1,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			headerStyle.Render(
				fmt.Sprintf("%s\n", m.presentation.Metadata.Title),
			),
		),
	)

	footer := lipgloss.Place(
		m.width,
		1,
		lipgloss.Center,
		lipgloss.Bottom,
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			leftFooter.Render(
				fmt.Sprintf("%s\t%s",
					m.presentation.Metadata.Author,
					m.presentation.Metadata.Date,
				)),
			rightFooter.Render(
				fmt.Sprintf("%d/%d", m.CurrentSlide+1, len(m.presentation.Slides)),
			),
		),
	)

	body := lipgloss.Place(
		m.width,
		m.height-(lipgloss.Height(header)+lipgloss.Height(footer)),
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Left,
			markdown,
		),
	)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		body,
		footer,
	)
}

func (m *model) GetCurrentSlide(slide, section int) string {
	content := m.presentation.Slides[slide].Content
	sections := parse.SplitSections(content)

	shown := strings.Join(sections[0:section+1], "\n")

	return shown
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", args[0])
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(args[1]))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		os.Exit(1)
	}
}
