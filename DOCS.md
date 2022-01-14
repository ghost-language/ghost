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

maker := CoffeeMaker()
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
    function constructor(recipe := "americano") {
        this.recipe := recipe
    }

    function setRecipe(recipe) {
        this.recipe := recipe

        // Returning 'this' allows you to chain methods
        return this
    }

    function brew() {
        print("Brewing and making your %s.".format(this.recipe))
    }
}

maker := CoffeeMaker()

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
`x := 5; x += 1; // x == 6`
`x := 5; x -= 1; // x == 4`
`x := 5; x *= 2; // x == 10`
`x := 10; x /= 2; // x == 5`

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

## Modularity
Ghost employs a simple module system to split and organize code into self-contained files.

Every ghost file is its own module with its own scope. Importing a file into another does not explicitely merge its scope. For example, two modules can define the same top-level variable with the same name without causing any name collision.

### Shared Imports
Ghost keeps track of every file it imports. Importing a module in multiple locations will not execute or load that module every time. The first encounter of the imported module will be the only time its loaded and evaluated.

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