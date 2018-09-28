# dexec [![GoDoc](https://godoc.org/github.com/ahmetalpbalkan/dexec?status.png)][godoc]


`dexec` is a small Go library allowing you to run processes inside 
Docker containers as if you are running them locally using [`os/exec`][osexec] package.
Read documentation at [godoc.org][godoc] or [see examples](examples).

Using `dexec`, you can do stream processing, off-load computationally
expensive parts of your program to a remote fleet of Docker engines.
Its interface is [strikingly similar][godoc] to [`os/exec`][osexec].

[osexec]: https://godoc.org/os/exec
[godoc]: https://godoc.org/github.com/ahmetalpbalkan/dexec

### Examples

Check out the following [examples](examples):

- [Hello, world inside container →](examples/100-hello)
- [Connect to container’s `STDIN`/`STDOUT` →](examples/200-stdin-stdout)
- [Stream processing with pipes →](examples/300-pipes)
- [Check exit code of a remote process →](examples/400-exit-code)
- [Audio extraction from YouTube videos →](examples/500-video-processing)
- [Parallel computation on Swarm →](examples/600-parallel-compute)

### A Short Example 

It takes only a **4-line code change** to convert a piece of code
using `os/exec` to use `dexec` to start running your stuff inside containers.

Here is a minimal Go program that runs `echo` in a container:

```go
package main

import (
	"log"

	"github.com/ahmetalpbalkan/dexec"
	"github.com/fsouza/go-dockerclient"
)

func main(){
	cl, _ := docker.NewClient("unix:///var/run/docker.sock")
	d := dexec.Docker{cl}

	m, _ := dexec.ByCreatingContainer(docker.CreateContainerOptions{
	Config: &docker.Config{Image: "busybox"}})

	cmd := d.Command(m, "echo", `I am running inside a container!`)
	b, err := cmd.Output()
	if err != nil { log.Fatal(err) }
	log.Printf("%s", b)
}
```

Output: `I am running inside a container!`

### Use Cases

This library is intended for providing an execution model that looks and feels
like [`os/exec`][osexec] package.

You might want to use this library when:

- You want to execute a process, but run it in a container with extra security
  and finer control over resource usage with Docker –and change your code
  minimally.

- You want to execute a piece of work on a remote machine (or even better, a pool
  of machines or a cluster) through Docker. Especially useful to distribute
  computationally expensive workloads.

For such cases, this library abstracts out the details of executing the process
in a container and gives you a cleaner interface you are already familiar with.

[Check out more examples →](examples)

[Read more →](https://ahmetalpbalkan.com/blog/dexec/)

![Analytics](https://ga-beacon.appspot.com/UA-45321252-5/welcome-page)
