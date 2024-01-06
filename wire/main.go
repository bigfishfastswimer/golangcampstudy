package wire

import "fmt"

func UseRepository() {
	repo := InitRepository()
	fmt.Print(repo)
}
