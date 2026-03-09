package views

import "github.com/nixbpe/ari/internal/checker"

// ANSI Color Codes for Cyberpunk Theme
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"

	// Foreground Colors
	Cyan    = "\033[36m"
	Magenta = "\033[35m"
	Yellow  = "\033[33m"
	Green   = "\033[32m"
	Red     = "\033[31m"
	White   = "\033[37m"
	Black   = "\033[30m"

	BrightCyan    = "\033[96m"
	BrightMagenta = "\033[95m"
	BrightYellow  = "\033[93m"
	BrightGreen   = "\033[92m"
	BrightRed     = "\033[91m"
	BrightWhite   = "\033[97m"

	// Background Colors
	BgCyan    = "\033[46m"
	BgMagenta = "\033[45m"
	BgYellow  = "\033[43m"
	BgBlack   = "\033[40m"
)

// LevelColor returns the specific neon color for a maturity level
func LevelColor(lvl checker.Level) string {
	switch lvl {
	case checker.LevelFunctional:
		return BrightRed
	case checker.LevelDocumented:
		return BrightYellow
	case checker.LevelStandardized:
		return BrightCyan
	case checker.LevelOptimized:
		return BrightGreen
	case checker.LevelAutonomous:
		return BrightMagenta
	default:
		return Dim + White
	}
}

const CyberHeader = BrightCyan + `
    ___    ____  ____
   /   |  / __ \/  _/
  / /| | / /_/ // /  
 / ___ |/ _, _// /   
/_/  |_/_/ |_/___/   ` + Reset + BrightMagenta + `
:: AGENT READINESS INDEX ::
` + Reset
