package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ttacon/chalk"
)

func rgb(i int) (int, int, int) {
	var f = 0.1
	return int(math.Sin(f*float64(i)+0)*127 + 128),
		int(math.Sin(f*float64(i)+2*math.Pi/3)*127 + 128),
		int(math.Sin(f*float64(i)+4*math.Pi/3)*127 + 128)
}

func rainbow(output []rune, offset int) string {
	sb := strings.Builder{}
	for j := 0; j < len(output); j++ {
		r, g, b := rgb(j + offset)
		sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, output[j]))
	}
	return sb.String()
}

func main() {
	var yearly float64 = 100_000
	chStyle := chalk.Bold.NewStyle().WithForeground(chalk.Yellow).WithBackground(chalk.Blue)
	fmt.Printf("\n\t%sBoring Meeting Mode: %s[ACTIVATED]\n\n%s", chalk.Cyan, chStyle, chalk.Reset)
	perMilli := yearly / 52 / 5 / 8 / 60 / 60 / 1000
	var earned float64 = 0
	ticker := time.Tick(time.Millisecond)
	startTime := time.Now()
	for {
		<-ticker
		earned += perMilli
		fmt.Print("\r\tEarned: ", rainbow([]rune(fmt.Sprintf("€%.10f", earned)), int(time.Now().Sub(startTime).Seconds())))
	}
}
