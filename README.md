# Creating your own Docker with Go
Docker, and the containers it makes has completely changed the way of packaging applications and deploying them. It helps injecting our source code with mobility to run at scale. I use docker for practically everthing on my laptop, creating local development environments, testing environments and using disposable containers for doing some crazy stuff and deleting them right away if things get messy😅. I am sure that everyone who has ever used docker must be amazed that how a single command lets us create isolated and independent machines(called containers) within a few seconds. <br>
Let's understand this magic while implementing our own **docker** with few lines of some awesome **go** code.

## What we'll do?
At the end of this post we'll be having a go binary that would be capable of running any valid linux command inside a isolated process(practically known as container).
```bash
docker         run 	 image          <cmd> <params>
go run main.go run   {some command}    <cmd> <params>
```

## Requirements
- **Go sdk**(Linux)
- Any **Linux** distribution
- **Docker** for linux<br><br>
**Linux** is required because containers are practically a wrap around Linux technologies that we'll be exploring next.

## Some Linux technologies
- **Namespaces** - what an isolated process can see is defined and controlled by namespaces. It creates isolation by providing each process it's own pseudo environment.
- **Chroots** - it control root filesystem for each process.
- **Cgroups** - what an isloted process can use as resource from host machine is enforced by cgroups.

## Container v/s Host
Enough with theory and definitions, now let's see how a container is different from host machine.<br><br>
We'll create a ubuntu docker container passing `/bin/bash` as entrypoint. Use the following snippet.
```bash
docker run -it --rm ubuntu /bin/bash
```
We'll run a few commands inside both our container(ubuntu 20.04) and host machine(ubuntu 20.04) and observer their behaviour inside both environments:

- **hostname** - return name of the host inside which bash is running.
- **ps** - return list of active process running inside the environment.

**container**
<img src="assets/container.png">

**host**
<img src="assets/host.png">

We can see when we run same commands in docker container and our host machine we get different results.
- Container is assigned a hostname from docker(container ID), while our system have a completely different hostname.
- Lots of process are running inside our host but our container is only aware of process running inside it, thus providing isolation.

## Let's dive deep
we have got a taste of a how containers function. It's time to open our editor and write some go code to do something similar that docker does.<br>
Create a `main.go` file with main package.

- **creating command switch**
```go
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("invalid command")
	}
}
```
the command switch picks the command-line argument passed, and runs the function mapped to that argument


