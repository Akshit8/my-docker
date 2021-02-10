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
- Go sdk(Linux)
- Any Linux distribution
<br>
**Linux** is required because containers are practically a wrap around Linux technologies that we'll be looking next.