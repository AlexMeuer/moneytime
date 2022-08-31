package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	flag "github.com/spf13/pflag"
	"golang.org/x/term"
)

type RainbowMode string

const (
	RainbowFixed RainbowMode = "fixed"
	RainbowAnim  RainbowMode = "anim"
	RainbowPuke  RainbowMode = "puke"
)

func moneyEarnedSince(start time.Time, moneyPerHour float64) float64 {
	return time.Since(start).Hours() * moneyPerHour
}

func rgb(i int) (int, int, int) {
	var f = 0.1
	return int(math.Sin(f*float64(i)+0)*127 + 128),
		int(math.Sin(f*float64(i)+2*math.Pi/3)*127 + 128),
		int(math.Sin(f*float64(i)+4*math.Pi/3)*127 + 128)
}

func rainbow(output []rune, offset int, mode RainbowMode) string {
	sb := strings.Builder{}
	for j := 0; j < len(output); j++ {
		var transformedIndex int
		switch mode {
		default:
			fallthrough
		case RainbowFixed:
			transformedIndex = j
		case RainbowAnim:
			transformedIndex = j + offset
		case RainbowPuke:
			transformedIndex = j * offset
		}
		r, g, b := rgb(transformedIndex)
		sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, output[j]))
	}
	return sb.String()
}

type model struct {
	moneyPerHour   float64
	earned         float64
	currencyPrefix string
	startTime      time.Time
	frameCount     uint8
	rainbowMode    RainbowMode
}

func (m model) Init() tea.Cmd {
	return nil
}

type tickMsg struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.frameCount = (m.frameCount + 1) % 255
		m.earned = moneyEarnedSince(m.startTime, m.moneyPerHour)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width, height = 20, 10
	}
	var activatedForeground lipgloss.Color
	var activatedBackground lipgloss.Color
	var activatedDecorator string
	if m.frameCount%64 < 32 {
		activatedDecorator = "|"
	} else {
		activatedDecorator = "◊"
	}
	if m.frameCount%128 < 64 {
		activatedForeground = lipgloss.Color("#494d64")
		activatedBackground = lipgloss.Color("#eed49f")
	} else {
		activatedForeground = lipgloss.Color("#eed49f")
		activatedBackground = lipgloss.Color("#494d64")
	}
	container := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#f5bde6")).
		Foreground(lipgloss.Color("#f5a97f")).
		Padding(1, 2).
		Margin(1, 2).
		Bold(true)
	titleLeft := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Italic(true).
		Render("Boring meeting mode:")
	titleRight := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Background(activatedBackground).
		Foreground(activatedForeground).
		Blink(true).
		MarginLeft(4).
		Render(fmt.Sprintf("[ %s ACTIVATED %s ]", activatedDecorator, activatedDecorator))
	title := lipgloss.JoinHorizontal(lipgloss.Center, titleLeft, titleRight)
	subtitle := lipgloss.NewStyle().
		Italic(true).
		MarginTop(1).
		MarginBottom(1).
		Render("Chill, look at how much you're earning just sitting here:")
	value := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f5a97f")).
		Bold(false).
		Render(rainbow([]rune(fmt.Sprintf("%s %.8f", m.currencyPrefix, m.earned)), int(m.frameCount), m.rainbowMode))
	return lipgloss.Place(width, height,
		lipgloss.Center, lipgloss.Center,
		container.Render(lipgloss.JoinVertical(lipgloss.Center, title, subtitle, value)),
		lipgloss.WithWhitespaceChars("お金"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#363a4f")),
	)
}

func main() {
	yearlyPtr := flag.Float64P("yearly", "s", 0, "Yearly salary")
	hourlyPtr := flag.Float64("hourly", 0, "Hourly wage")
	fpsPtr := flag.Int("fps", 60, "Frames per second")
	prettyPtr := flag.String("pretty", "fixed", "off = plain text, fixed = fixed colours, anim = animated colours, puke = awful colours")
	currencyPrefixPtr := flag.StringP("currencyPrefix", "c", "€", "The currency symbol/prefix to use when pretty mode is enabled")
	startTimeOffsetPtr := flag.DurationP("startTimeOffset", "o", 0, "The time offset to use when calculating the money earned")
	flag.Parse()
	fps := *fpsPtr
	yearly := *yearlyPtr
	hourly := *hourlyPtr
	prettyMode := *prettyPtr
	if fps <= 0 {
		fps = 60
	}
	if yearly <= 0 && hourly <= 0 {
		fmt.Println("Please specify a positive yearly salary or hourly wage.")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if yearly > 0 {
		hourly = yearly / 52 / 5 / 8
	}

	tick := time.Tick(time.Second / time.Duration(fps))

	if prettyMode == "off" {
		startTime := time.Now().Add(*startTimeOffsetPtr)
		for range tick {
			fmt.Printf("\r%.10f", moneyEarnedSince(startTime, hourly))
		}
		return
	}

	p := tea.NewProgram(model{
		currencyPrefix: *currencyPrefixPtr,
		startTime:      time.Now().Add(*startTimeOffsetPtr),
		moneyPerHour:   hourly,
		rainbowMode:    RainbowMode(prettyMode),
	})
	go func() {
		for range tick {
			p.Send(tickMsg{})
		}
	}()
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
