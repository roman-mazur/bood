The simplest bood example
=========================

From this directory, try the following commands:

#### Install bood
```
$ go get -u github.com/roman-mazur/bood/cmd/bood
```

#### Build the program
```
$ bood
INFO 2020/03/10 23:37:46 Ninja build file is generated at out/build.ninja
INFO 2020/03/10 23:37:46 Starting the build now
[1/1] Build hello as Go binary
```

#### Run the program
```
$ out/bin/hello
Hello!
```

#### Run build again (and see nothing is done)
```
$ bood
INFO 2020/03/10 23:37:50 Ninja build file is generated at out/build.ninja
INFO 2020/03/10 23:37:50 Starting the build now
ninja: no work to do.
```
