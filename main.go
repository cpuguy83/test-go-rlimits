package main

/*
#cgo CFLAGS: -Wall
extern void ulimit();
void __attribute__((constructor)) init(void) {
	ulimit();
}
*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	if len(os.Args) == 1 {
		lim := getRlimit()
		fmt.Println("parent:", "rlim_cur:", lim.Cur, "rlim_max:", lim.Max)
		if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{Cur: 1024, Max: 4096}); err != nil {
			panic(err)
		}
		syscall.Exec(os.Args[0], []string{os.Args[0], "reexec"}, os.Environ())
		return
	}

	if os.Args[1] == "reexec" {
		lim := getRlimit()
		fmt.Println("reexec:", "rlim_cur:", lim.Cur, "rlim_max:", lim.Max)
		env := os.Environ()
		env = append(env, "_ULIMIT=1")
		syscall.Exec(os.Args[0], []string{os.Args[0], "child"}, env)
	}
}

func getRlimit() syscall.Rlimit {
	var lim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim); err != nil {
		panic(err)
	}
	return lim
}
