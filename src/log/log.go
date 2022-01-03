package log

import (
	"fmt"
)

const (
	// ANSI terminal escape codes for color output
	AnsiReset     = "\033[0;0m"
	AnsiBlue      = "\033[34;22m"
	AnsiGreen     = "\033[32;22m"
	AnsiRed       = "\033[31;22m"
	AnsiBlueBold  = "\033[34;1m"
	AnsiGreenBold = "\033[32;1m"
	AnsiRedBold   = "\033[31;1m"
)

func Debug(str string, args ...interface{}) {
	fmt.Println(AnsiBlueBold + "debug: " + AnsiBlue + fmt.Sprintf(str, args...) + AnsiReset)
}

func Info(str string, args ...interface{}) {
	fmt.Println(AnsiGreen + fmt.Sprintf(str, args...) + AnsiReset)
}

func Error(str string, args ...interface{}) {
	fmt.Println(AnsiRed + fmt.Sprintf(str, args...) + AnsiReset)
}
