# The Ghost Programming Language

[![Test](https://github.com/ghost-language/ghost/actions/workflows/test.yml/badge.svg)](https://github.com/ghost-language/ghost/actions/workflows/test.yml)

Ghost is a small, object-oriented, embeddable toy scripting language. While object-oriented, Ghost also supports procedural and functional programming styles as well.

Ghost is dynamically typed, runs by a tree-walking interpreter, and has automatic memory management thanks to its implementation through the Go programming language.

## Status

> Currently in beta, vetting out the language and seeing how it feels writing/running. Major changes are still possible at this stage.

## Documentation

You will find robust, user friendly, and updated documentation on our [website](https://ghostlang.org/docs).

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/ghost-language/ghost/tags).

## Installation

### Brew

If you're on mac, you may use `homebrew`:

```
$ brew tap ghost-language/ghost
$ brew install ghost-language/ghost/ghost
```

### Go Install

If you have Go installed, you may use `go install`:

```
go install ghostlang.org/x/ghost
```

### Direct Download

You may download the compiled binaries for your platform from our GitHub [releases](https://github.com/ghost-language/ghost/releases) page.

## Development

- To build and execute Ghost, run `make`.
- To build Ghost, run `make build`.
- To execute tests, run `make test`.

```
$  git clone git@github.com:ghost-language/ghost.git
$  cd ghost
$  make
   Ghost (x.x)
   Press Ctrl + C to exit

>>
```

## CLI

You can execute code written in Ghost in various ways using the CLI.

### REPL

Ghost includes a simple REPL to write and execute Ghost code directly in your terminal. To enter the REPL environment, run `ghost` on its own:

```
$  ghost
   Ghost (x.x)
   Press Ctrl + C to exit

>>
```

### Executing Files

To execute a Ghost source file (`.gs`), pass either the relative or absolute path of the file to `ghost`. The source file will be executed and then exit back to the terminal.

```
$  ghost examples/fib.gs
   34
$
```

### Interactive Mode

Interactive mode allows you to execute a Ghost source file, and then immediately jump into a REPL session with the state of the source file still intact. To enter interactive mode, pass the `-i` flag to `ghost`.

```
$  ghost -i examples/fib.gs
   34
>> fibThree(9)
   34
>>
```

## Releasing

Ghost is hosted and distributed through GitHub. Publishing a release on GitHub
triggers the `Release` workflow, which runs [GoReleaser](https://goreleaser.com)
to build every binary, attach the archives to the release, and update the
Homebrew cask. Maintainers do not run GoReleaser by hand.

1. Update the internal version reference so it matches the tag you are about to
   create. The release workflow refuses to publish if the two disagree:

   ```go
   // version/version.go

   const Version = "x.y.z"
   ```

2. Merge that change, then create and push the tag:

   ```
   $ git tag -a vx.y.z -m "Release description"
   $ git push origin vx.y.z
   ```

3. Publish a release for the tag on GitHub, writing whatever notes you want.
   The generated changelog is appended to them. Publishing is what starts the
   workflow - a draft release does not.

To exercise the whole pipeline without cutting a tag, run the `Release`
workflow manually from the Actions tab. It builds and archives everything,
publishes nothing, and attaches the archives to the run.

### Homebrew

The cask is pushed to [ghost-language/homebrew-ghost](https://github.com/ghost-language/homebrew-ghost),
a separate repository that the workflow's automatic `GITHUB_TOKEN` cannot write
to. It needs a personal access token with `contents:write` on that repository,
stored as the `HOMEBREW_TAP_TOKEN` secret. Without it the release still
publishes; only the cask update is skipped.
