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


### Install `gows`

```
go get -u github.com/bitrise-tools/gows
```

That's all. If you have a "properly" configured Go environment (see the previous Install section)
then you should be able to run `gows -version` now, to be able to run `gows` in any directory.


## Usage

Just prefix your commands (any command) with `gows`.

Example:

Instead of `go get ./...` use `gows go get ./...`. That's all :)


### Alternative usage option: jump into a prepared Shell

*This solution works for most shells, but there are exceptions, like `fish`.
The reason is: `gows` creates a symlink between your project and the isolated workspace.
In Bash and most shells if you `cd` into a symlink directory (e.g. `cd my/symlink-dir`)
your `pwd` will point to the symlink path (`.../my/symlink-dir`),
but a few shells (`fish` for example) change the path to the symlink target path instead,
which means that when `go` gets the current path that will point to your project's original
path, instead of to the symlink inside the isolated workspace. Which, at that point,
is outside of GOPATH.*

In shells which keep the working directory path to point to the symlink's path, instead
of it's target (e.g. Bash) you can run:

```
gows bash
```

or

```
gows bash -l
```

which will start a shell (in this example Bash) with prepared `GOPATH` and your
working directory will be set to the symlink inside the isolated workspace.

If you want to use this mode you'll have to change how you initialize
your `GOPATH`, to allow it to be overwritten by `gows` for "shell jump in".
To allow `gows` to overwrite the `GOPATH` for shells initialized **by/through** `gows` you should
change your `GOPATH` init entry in your `~/.bash_profile` / `~/.bashrc` (or
wherever you did set GOPATH for your shell).
For Bash (`~/.bash_profile` / `~/.bashrc`) you can use this form:

```
if [ -z "$GOPATH" ] ; then
  export GOPATH="/my/go/path"
fi
```

instead of this one:

```
export GOPATH="/my/go/path"
```

This means that your (Bash) shell will only set the `GOPATH` environment if it's not set to a value already.

This is not required if you use `gows` only in a "single command / prefix" style,
it's only required if you want to initialize
the shell and jump into the initialized shell through `gows`. In general it's safe to initialize the
environment variable this way even if you don't plan to initialize any shell through `gows`,
as this will always initialize `GOPATH` *unless* it's already initialized (e.g. by an outer shell).


## TODO

- [ ] Setup the base code (generate the template project, e.g. create a new Xcode project or `rails new`)
  - [ ] commit & push
- [ ] Add linter tools
  - go:
    - [ ] `go test`
    - [ ] `go vet`
    - [ ] [errcheck](github.com/kisielk/errcheck)
    - [ ] [golint](github.com/golang/lint/golint)
- [ ] Write tests & base functionality, BDD/TDD preferred
- [ ] Setup continuous integration (testing) on [bitrise.io](https://www.bitrise.io)
- [ ] Setup continuous deployment for the project - just add it to the existing [bitrise.io](https://www.bitrise.io) config
- [ ] Use [releaseman](https://github.com/bitrise-tools/releaseman) to automate the release and CHANGELOG generation
- [ ] Iterate on the project (and on the automation), test the automatic deployment
