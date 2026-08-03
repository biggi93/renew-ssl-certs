package post

import (
	"fmt"
	"os/exec"
)

func Run() {
	out, err := exec.Command("post/post.sh").Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}