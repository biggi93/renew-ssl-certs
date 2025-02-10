package certbot

import (
	"fmt"
	"os/exec"
)

func Run() {
	out, err := exec.Command("certbot/certbot.sh").Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}