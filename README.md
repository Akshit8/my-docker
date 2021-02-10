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

- **command switch**
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
The command switch picks the command-line argument passed, and runs the function mapped to that argument.

- **run function**
```go
func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	// proc dir is a directory where all processes metadata is there
    // our temporary binary will also be present here
    // below line executes child function inside the newly created container
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
    // attatching os-std process to our cmd-std process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

    // setting some system process attributes
    // below line of code is responsible for creating a new isolated process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Cloning is what creates the process(container) in which we would be running our command.
		// CLONE_NEWUTS will allow to have our own hostname inside our container by creating a new unix timesharing system.
		// CLONE_NEWPID assigns pids to only process inside the new namspace.
		// CLONE_NEWNS new namespace for mount.
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// Mounts in systemd gets recursively shared property.
		// Unshare the recursively shared property for new mount namespace.
		// It prevents sharing of new namespace with the host.
		Unshareflags: syscall.CLONE_NEWNS,
	}

    // running the command and catching error
	if err := cmd.Run(); err != nil {
		log.Fatal("Error: ", err)
	}
}
```
The `run()` function is responsible for creating an isolated process(container) and then execute itself inside this isolated process, this time invoking `child()` function

- **child function**
```go
func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
    // below are some system calls that set some container properties
	// sets hostname for newly created namespace
	must(syscall.Sethostname([]byte("maverick")))
    // setting root director for the container
	must(syscall.Chroot("/"))
    // making "/" as default dir
	must(syscall.Chdir("/"))
    // mounting proc dir to see the process running inside container
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

    // below line finally executes the user-command inside our own created container!
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
    // attatching os-std process to our cmd-std process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

    // running the command and catching error
	if err := cmd.Run(); err != nil {
		log.Fatal("Error: ", err)
	}
    // unmount the proc after command is finished
	syscall.Unmount("/proc", 0)
}
```
The `child()` function is invoked by run function as a child process inside container created by `run()`. It is responsible for some system calls for setting some container properties and finally execute the command dispatched from user.

- **must function**
```go
func must(err error) {
	if err != nil {
		panic(err)
	}
}
```
The `must()` is a simple error wrapper that panics if any system call invoked inside child function fails.

## Let's create some containers
Now we have a mini docker program that can actually create isolated and independent containers on your host machine. Let's fire up our powerful code.<br>
<br>
The command that we would be running inside our container is `/bin/bash`, that will start a new bash program inside our container.<br>
<br>
Run following snippet to create the container.
```bash
go run main.go run /bin/bash
```
<img src="assets/execute.png">

