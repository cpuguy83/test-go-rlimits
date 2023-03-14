This is attempting to test changes to the go-runtime related to automatically
setting the RLIMIT_NOFILE lower limit. This is a test to see if the changes work
as expected.

## Background

go1.19 made change such that anyone who imports the "os" package will have the
lower limit of RLIMIT_NOFILE changed to be the same limit as the max.

Example (psuedo code):

```
lim = getrlimit(RLIMIT_NOFILE)
if lim.Cur != lim.Max {
  lim.Cur = lim.Max
  setrlimit(RLIMIT_NOFILE, lim)
}
```

This has some potential implications for applications that are not expecting
this change, especially programs that execute other, non-go programs which may
be using the `select` system call. Technically also a problem for go programs
however it would be weird for a go program to ever use this system call (nothing
to do with the `select` keyworrd in go).

## What to test?

From go's HEAD commit, cherry-pick the following commits in order:

- https://go-review.googlesource.com/c/go/+/476096 (moves the code out of the `os` package and into `syscall`)
- https://go-review.googlesource.com/c/go/+/476097 (restore original rlimit on exec)

Set your GOROOT and PATH to point to the go runtime you are testing. Then call `go run .` in this directory.