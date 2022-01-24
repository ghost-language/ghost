# Documentation

## Error Handling

- Syntax errors
- Runtime errors
- Line and column are reported for every error

## Classes

### Defining

```dart
class CoffeeMaker {
    //
}
```

### Constructor

```dart
class CoffeeMaker {
    function constructor() {
        print("Calibrating your coffee maker.")
    }
}

maker = CoffeeMaker.new()
```

### Methods

```dart
class CoffeeMaker {
    function brew() {
        print("Your coffee is now brewing.")
    }
}
```

### this

```dart
class CoffeeMaker {
    function constructor(recipe = "americano") {
        this.recipe = recipe
    }

    function setRecipe(recipe) {
        this.recipe = recipe

        // Returning 'this' allows you to chain methods
        return this
    }

    function brew() {
        print("Brewing and making your %s.".format(this.recipe))
    }
}

maker = CoffeeMaker.new()

maker.setRecipe("latte").brew()
```

## Values

### Lists

#### Methods

- `first()`
- `join()`
- `last()`
- `length()`
- `pop()`
- `push()`
- `tail()`
- `toString()`

### Numbers

#### Compound operators

`x = 5; x += 1; // x == 6`
`x = 5; x -= 1; // x == 4`
`x = 5; x *= 2; // x == 10`
`x = 10; x /= 2; // x == 5`

### Strings

#### Methods

- `find()`
- `findAll()`
- `format()`
- `endsWith()`
- `length()`
- `matches()`
- `replace()`
- `split()`
- `startsWith()`
- `toLowerCase()`
- `toUpperCase()`
- `toString()`
- `toNumber()`
- `trim()`
- `trimEnd()`
- `trimStart()`

## Switch Statements

The `switch` statement evaluates an expression, matching the expression's value to a `case` clause, and executes statements associated with that `case`, as well as statements in cases that follow the matching `case`.

```typescript
beverage = "coffee"

switch (beverage) {
    case "water" {
        print("Water is $0.75 per bottle.")
    }
    case "juice" {
        print("Juice is $1.25 per bottle.")
    }
    case "coffee", "latte" {
        print("Coffee and lattes are $2.75 per 12oz.")
    }
}
```

## Modularity

Ghost employs a simple module system to split and organize code into self-contained files.

Every ghost file is its own module with its own scope. Importing a file into another does not explicitely merge its scope. For example, two modules can define the same top-level variable with the same name without causing any name collision.

### Shared Imports

Ghost keeps track of every file it imports. Importing a module in multiple locations will not execute or load that module every time. The first encounter of the imported module will be the only time its loaded and evaluated.

### Binding Variables

All top-level variables within a module are exportable. To actually _import_ data, you may specify any number of identifiers in your import statement:

```typescript
import Request, Response from "http"
```

The above will _import_ and _bind_ the values `Request` and `Response` from the `http` module. This will make `Request` and `Response` available in your file.

#### Aliases

You may import a variable under a different name using `as`:

```typescript
import str as isString from "helpers"
```

### Cyclic Imports

Cyclic imports for the most part are "supported" by Ghost, though they should still be considered a code smell if you ever come across them. Because Ghost keeps track of the modules it imports, it's effectively able to short-circuit itself on cyclic imports:

```typescript
// main.ghost
import "a";

// a.ghost
print("start a");
import "b";
print("end a");

// b.ghost
print("start b");
import "c";
print("end b");

// c.ghost
print("start c");
import "a";
print("end c");
```

When running the above program, you'll find that it prints the correct output and does not get stuck in an infinite loop:

```
start a
start b
start c
end c
end b
end a
```

### Importing Imperatively

To import a file imperatively, simply use the `import` statement:

```dart
import "beverages"
```

This will evaluate the module and run it, but it will not bind any new variables.

## Standard Library

### Functions

- `print()`
- `type()`

### Console

#### Methods

- `console.error()`
- `console.info()`
- `console.log()`
- `console.read()`
- `console.warn()`

### Ghost

#### Methods

- `ghost.abort()`
- `ghost.execute()`
- `ghost.extend()`
- `ghost.identifiers()`

#### Properties

- `ghost.version`

### HTTP

#### Methods

- `http.handle()`
- `http.listen()`

### IO

#### Methods

- `io.append()`
- `io.read()`
- `io.write()`

### Math

#### Methods

- `math.abs()`
- `math.cos()`
- `math.isNegative()`
- `math.isPositive()`
- `math.isZero()`
- `math.sin()`
- `math.tan()`

#### Properties

- `math.pi`
- `math.e`
- `math.epsilon`
- `math.tau`

### OS

#### Methods

- `os.args()`
- `os.clock()`
- `os.exit()`

#### Properties

- `os.name`

### Random

#### Methods

- `random.seed()`
- `random.random()`
- `random.range()`

#### Properties

- `random.seed`

### Time

#### Methods

- `time.sleep()`
- `time.now()`

#### Properties

- `time.nanosecond`
- `time.microsecond`
- `time.millisecond`
- `time.second`
- `time.minute`
- `time.hour`
- `time.day`
- `time.week`
- `time.month`
- `time.year`
