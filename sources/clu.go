package sources

import (
	"fmt"
	"os"
)

func thing() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("pwd is", pwd)
}
