package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// docker         run 	 image          <cmd> <params>
// go run main.go run {some command}    <cmd> <params>

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

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	// proc dir is a directory where all processes metadata is there
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Cloning is what creates the process(container) in which we would be running our command.
		// CLONE_NEWUTS will allow to have our own hostname inside our container.
		// CLONE_NEWPID assigns pids to only process inside the new namspace.
		// CLONE_NEWNS namspace for mount.
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// Mounts in systemd gets recursively shared property.
		// Unshare the recursively shared property for new mount namespace.
		// It prevents sharing of new namespace with the host.
		Unshareflags: syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		log.Fatal("Error: ", err)
	}
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	// sets hostname for newly created namespace
	must(syscall.Sethostname([]byte("maverick")))
	must(syscall.Chroot("/sample-root"))
	must(syscall.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("Error: ", err)
	}

	syscall.Unmount("/proc", 0)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}