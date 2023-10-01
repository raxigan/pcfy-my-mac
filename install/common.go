package install

import "fmt"

const Yellow = "\033[33m"
const Red = "\033[31m"
const Blue = "\033[34m"
const Green = "\033[32m"
const Reset = "\033[0m"

func getOrDefaultString(launcher string, launcher2 *string) string {
	if launcher2 != nil {
		return *launcher2
	} else {
		return launcher
	}
}

func getOrDefaultSlice(launcher []string, launcher2 *[]string) []string {
	if launcher2 != nil {
		return *launcher2
	} else {
		return launcher
	}
}

func printColored(color, msg string) {
	fmt.Print(Colored(color, msg))
}

func printlnColored(color, msg string) {
	fmt.Println(Colored(color, msg))
}

func Colored(color, msg string) string {
	return color + msg + Reset
}
