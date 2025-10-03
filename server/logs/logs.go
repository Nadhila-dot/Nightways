package logg

import (
	"fmt"
	"log"
)

func Error(message string) {
	fmt.Printf("\033[1;31m[VELA]: %s\033[0m\n", message)
}

func Success(message string) {
	fmt.Printf("\033[1;32m[VELA]: %s\033[0m\n", message)
}

func Info(message string) {
	fmt.Printf("\033[1;34m[VELA]: %s\033[0m\n", message)
}

func Warning(message string) {
	fmt.Printf("\033[1;33m[VELA]: ⛔️%s\033[0m\n", message)
}

func Exit() {
	log.Fatalf("\033[1;31m[VELA]: Exiting... Bye Bye\033[0m")
}
