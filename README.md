# The Ghost Programming Language
Ghost is a small, object-oriented, embeddable toy scripting language. While object-oriented, Ghost also supports procedural and functional programming styles as well.

Ghost is dynamically typed, runs by a tree-walking interpreter, and has automatic memory management thanks to its implementation through the Go programming language. Ghost's implementation in Go has no dependencies.

## Status
> Currently in active development. Not feature complete.

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
go install github.com/ghost-language/ghost
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
   Ghost (dev-nightly)
   Press Ctrl + C to exit

>>
```

## CLI
You can execute code written in Ghost in various ways using the CLI.

### REPL
Ghost includes a simple REPL to write and execute Ghost code directly in your terminal. To enter the REPL environment, run `ghost` on its own:

```
$  ghost
   Ghost (dev-nightly)
   Press Ctrl + C to exit

>>
```

### Executing Files
To execute a Ghost source file (`.ghost`), pass either the relative or absolute path of the file to `ghost`. The source file will be executed and then exit back to the terminal.

```
$  ghost examples/fibtc.ghost
   9227465
$
```

### Interactive Mode
Interactive mode allows you to execute a Ghost source file, and then immediately jump into a REPL session with the state of the source file still intact. To enter interactive mode, pass the `-i` flag to `ghost`.

```
$  ghost -i examples/fibtc.ghost
   (executed in: 350.374µs)
>> x
   9227465
>>
```

## Releasing
Ghost is hosted and distributed through GitHub. We utilize [GoReleaser](https://goreleaser.com) to automate the release process. GoReleaser will build all the necessary binaries, publish the release and publish the brew tap formula. The following steps outline the process for maintainers of Ghost:

1. Ensure you have a GitHub token with `repo` access saved to your environment:
  ```
  export GITHUB_TOKEN="YOUR_GH_TOKEN"
  ```
2. Ensure the internal version reference is updated:
   ```go
   // version/version.go
   
   var (
      Version = "x.y.z"
   )
   ```
3. Create a new tag:
  ```
  $ git tag -a vx.y.z -m "Release description"
  $ git push origin vx.y.z
  ```
4. Run GoReleaser:
  ```
  $ goreleaser
  ```

## Credits
- [Crafting Interpreters](https://craftinginterpreters.com/)
- [Writing An Interpreter In Go](https://interpreterbook.com/)

## License
Ghost is open-sourced software licensed under the MIT license. See the [LICENSE](LICENSE) file for complete details.
