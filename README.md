# pid

Daemonize your go programs by locking the [process id file](https://www.baeldung.com/linux/pid-file).
Simply put the following at the top of your `main`,

```
func main() {
    defer pid.Start("monokuma").Stop()
    // your code...
}
```

This will create a `monokuma.pid` file in `/tmp` with the process id (pid) of the go
program when loaded by your operating system. When your `main` exits, the file will be
cleaned up.

## Motivation

This will prevent any other instances of your program running at the same time (say running
a web server or overloading some resources). Say, when I run

```
λ ./monokuma
```

I can find the PID of it with (no line feed at the end),

```
λ cat /tmp/monokuma.pid
91100⏎
```

If I try to start the same program, I will be greeted by

```
λ ./monokuma
2023/05/17 23:03:03 your app with pid 91100 is already running
```

and exit with code 1 <- opiniated and expected behavior. Some edge-cases are handled below.

## Caveats

There are multiple instances where we have to find edge-cases, say, your program failed
and the cleanup of the pid didn't happen. What then?

`pid` will check for the pid file's existence and if it exists, we will read the pid value
and try to find whether there is a running process with said pid. If there is none, it is
considered a "stale" pid file and deleted. New pid file will be created.

However, if a pid from the file exists, we exit immediately with 1 like above. What if the
running pid *isn't* actually your program but some other program that commandeered said pid?
That [*shouldn't*](https://en.wikipedia.org/wiki/Process_identifier) happen on Unix-like systems,
because each pid is unique until the max value is reached (care to reboot your machine?)

On windows it's a different story, though. But hey, this `pid` library isn't supported on windows,
so this little code is *flawless*. hah.

As of May 17th, 2023, [`os.FindProcess`](https://pkg.go.dev/os#FindProcess) doesn't work on checking
a running pid,

> On Unix systems, FindProcess always succeeds and returns a Process for the given pid,
>regardless of whether the process exists.

Instead, I used fantastic [gopsutil](https://github.com/shirou/gopsutil)'s
[`PidExists`](https://pkg.go.dev/github.com/shirou/gopsutil/v3@v3.23.4/process#PidExists).

## TODO

Should we write to somewhere that is not `/tmp`? For security reasons, `/var/run` is locked behind
special permissions or user id 0, which I don't want to touch.

How do we manage pids on windows? I know there is a way, but I usually use Unix-style machines.

## LICENSE

[Apache licensed](./LICENSE), go wild.
