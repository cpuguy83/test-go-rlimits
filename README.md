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

Get https://go-review.googlesource.com/c/go/+/476097 and build it
Set your GOROOT and PATH to point to the go runtime you are testing.
Then call `go run .` in this directory.

The output should be (when its working correctly according to the above change) something like:

```
parent: rlim_cur: 1048576 rlim_max: 1048576
reexec: rlim_cur: 4096 rlim_max: 4096
 child: rlim_cur: 1024 rlim_max: 4096
```

Where the "parent" has the original ulimits inherited from you, the parent process.
Only go is setting the lower limit to match the max limit.
When the "reexec" process is spawned the test actually sets the lower limit down to 1024 and the upper to 4096.
The "reexec" process should, however, have the lower limit raised back up to 4096 as that is what go is setting it to.
The "child" process is execing C code and should have the lower limit set to 1024 because that's what we set it to before the "reexec".

With go1.19 and go1.20 the output is:

```
parent: rlim_cur: 1048576 rlim_max: 1048576
reexec: rlim_cur: 4096 rlim_max: 4096
 child: rlim_cur: 4096 rlim_max: 4096 
```

Or you can run the test script directly:

```console
$ ./run.sh
remote: Finding sources: 100% (13143/13143)
remote: Total 13143 (delta 1731), reused 8614 (delta 1731)
Receiving objects: 100% (13143/13143), 26.90 MiB | 8.90 MiB/s, done.
Resolving deltas: 100% (1731/1731), done.
From https://go.googlesource.com/go
 * branch              refs/changes/97/476097/3 -> FETCH_HEAD
Warning: you are leaving 1 commit behind, not connected to
any of your branches:

  778f2710 syscall: handle errors.ErrUnsupported in isNotSupported

If you want to keep it by creating a new branch, this may be a good time
to do so with:

 git branch <new-branch-name> 778f2710

HEAD is now at 7483a567 syscall: restore original NOFILE rlimit in child process
Building Go cmd/dist using /home/cpuguy83/dev/go. (devel go1.21-7483a56748 Tue Mar 14 21:19:42 2023 -0700 linux/amd64)
Building Go toolchain1 using /home/cpuguy83/dev/go.
Building Go bootstrap cmd/go (go_bootstrap) using Go toolchain1.
Building Go toolchain2 using go_bootstrap and Go toolchain1.
Building Go toolchain3 using go_bootstrap and Go toolchain2.
Building packages and commands for linux/amd64.
---
Installed Go for linux/amd64 in /home/cpuguy83/dev/test/go-setrlimit-change/go
Installed commands in /home/cpuguy83/dev/test/go-setrlimit-change/go/bin
parent: rlim_cur: 1048576 rlim_max: 1048576
reexec: rlim_cur: 4096 rlim_max: 4096
 child: rlim_cur: 1024 rlim_max: 4096 
```

You can set `USE_SYSTEM_GO=1` to use the system go instead of the one built from source which is useful for comparing versions.