# gows

Go Workspace / Environment Manager, to easily manage the Go Workspace during development.


## The idea

Work in **isolated (development) environment** when you're working on your Go projects.
**No cross-project dependency version missmatch**, no more packages left out from `vendor/`.

No need for initializing a go workspace either, **your project can be located anywhere**,
not just in a predefined `$GOPATH` workspace.
`gows` will take care about crearing the (per-project isolated) workspace directory
structure, no matter where your project is located.

`gows` **works perfectly with other Go tools**, all it does is it ensures
that every project gets it's own, isolated Go workspace and sets `$GOPATH`
accordingly.


## Install

### Install & Prepare Go

Install & configure `Go` - [official guide](https://golang.org/doc/install).

Make sure you have `$GOPATH/bin` in your `$PATH`, e.g. by adding

```
export PATH="$PATH:$GOPATH/bin"
```

to your `~/.bash_profile` or `~/.bashrc` (don't forget to `source` it or to open a new Terminal window/tab
if you just added this line to the profile file).

*This makes sure that the Go projects you build/install will be available in any Terminal / Command Line
window, without the need to type in the full path of the binary.*


**Allow GOPATH to be overwritten by `gows` for "shell jump in" (see the "Jump into a prepared Shell" section)**:
To allow `gows` to overwrite the GOPATH for shells initialized **by/through** `gows` you should
change your `GOPATH` entry in your `~/.bash_profile` / `~/.bashrc` (or wherever you did set GOPATH for your
shell). For Bash (`~/.bash_profile` / `~/.bashrc`) you can use this form:

Instead of:

```
export GOPATH="/my/go/path"
```

Use something like this:

```
if [ -z "$GOPATH" ] ; then
  export GOPATH="/my/go/path"
fi
```

This means that your (Bash) shell will only set the `GOPATH` environment if it's not set to a value already, if it's
empty.

This is not required if you use `gows` only in a "single command" style, it's only required if you want to initialize
the shell and jump into the initialized shell through `gows`. In general it's safe to initialize the
environment variable this way even if you don't plan to initialize any shell through `gows`,
as this will always initialize `GOPATH` *unless* it's already initialized (e.g. by an outer shell).


### Install `gows`

```
go get -u github.com/bitrise-tools/gows
```

That's all. If you have a "properly" configured Go environment (see the previous Install section)
then you should be able to run `gows -version` now, to be able to run `gows` in any directory.


## Usage

There are two ways to use `gows`, it's up to you to choose the one which
fits your work style the most.


### Jump into a prepared Shell

To start a prepared Bash shell:

```
gows bash
```

To start a prepared Fish shell:

```
gows fish
```

To start a prepared ... You got the idea ;)


### Single command / command prefix mode

Just prefix your commands (any command) with `gows`.

Example:

Instead of `go get ./...` use `gows go get ./...`. That's all :)


## TODO

- [ ] Setup the base code (generate the template project, e.g. create a new Xcode project or `rails new`)
  - [ ] commit & push
- [ ] Add linter tools
  - go:
    - [ ] `go test`
    - [ ] `go vet`
    - [ ] [errcheck](github.com/kisielk/errcheck)
    - [ ] [golint](github.com/golang/lint/golint)
    - [ ] __WEB__ [safesql](github.com/stripe/safesql)
- [ ] Write tests & base functionality, BDD/TDD preferred
- [ ] Setup continuous integration (testing) on [bitrise.io](https://www.bitrise.io)
- [ ] Setup continuous deployment for the project - just add it to the existing [bitrise.io](https://www.bitrise.io) config
- [ ] Use [releaseman](https://github.com/bitrise-tools/releaseman) to automate the release and CHANGELOG generation
- [ ] Iterate on the project (and on the automation), test the automatic deployment
