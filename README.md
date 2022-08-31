# Moneytime

A CLI tool to calculate and display money earned in real time.

![A screenshot showing the program running in pretty mode.](/screenshots/pretty.png)

## Installation

You can use the program directly via `go run main.go` or install it via `go install` and then run it as `moneytime`.

## Usage

Display a pretty UI with the money earned (starting from now) with a yearly salary:
`moneytime -s 12345`
or an hourly salary:
`moneytime --hourly 12345`

If your meeting started 10 minutes ago, you can use `-o -10m` to offset the calculation by negative 10 minutes.

`go run main.go -s 12345 --fps 60 -o -30m --pretty anim -c LOL`

The above command shows a pretty terminal UI with the money earned since 30 minutes ago and increasing over time (recalculated at 60 fps) with an animated rainbow and with a `LOL` currency prefix.

```plaintext
  -c, --currencyPrefix string      The currency symbol/prefix to use when pretty mode is enabled (default "€")
      --fps int                    Frames per second (default 60)
      --hourly float               Hourly wage
      --pretty string              off = plain text, fixed = fixed colours, anim = animated colours, puke = awful colours (default "fixed")
  -o, --startTimeOffset duration   The time offset to use when calculating the money earned
  -s, --yearly float               Yearly salary
```

If both yearly and hourly are specified, hourly will take precidence.
