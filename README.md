# containers

## Some linux concepts
- Namespace
- Chroot
- Cgroups

## Host v/s container
```bash
---Container---
# starting a ubuntu container
$ docker run --rm -it ubuntu /bin/bash
$ hostname
$ ps
$ exit

---Host---
$ hostname
$ ps
```

## Namespaces
It is responsible for what a process can see. It is created using a syscall.
- Unix Timesharing System
- Process IDs
- Mounts
- Networks
- User IDs
- InterProcess Communication

It forms an important part of container, as it restricting the view that process has of things going on host machine.

## Interactive results
```bash
(1) --- after creating CLONE_NEWUTS
# this would start the bash inside a new namespace
go run main.go run /bin/bash
$ hostname
"some hostname inherited from host"
$ hostname maverick
"maverick"

# Now open a new bash and run
$ hostname
"host's original hostname"

(2) --- after setting hostname maverick
go run main.go run /bin/bash
$ root@maverick

(3) --- after creating CLONE_NEWPID
go run main.go run /bin/bash
$ ps # will still return all process id's

(4) --- after setting chroot chdir
go run main.go run /bin/bash
$ ls # must see contents of sample-root
$ ps # won't work

(5) --- after setting mounting proc
go run main.go run /bin/bash
$ ps
$ mount grep | proc

(6) --- after creating CLONE_NEWNS for mount
go run main.go run /bin/bash
$ mount grep | proc
```

## Cgroups
As Namespace can restrict what a process inside container can **see** from host machine, Control groups restrict what a process can **use** from a host machine. It is created using a pseudo file system that emits configuration files defining availaible resource to that process.


