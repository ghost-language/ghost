import "ghost:ghost"

console.log("Ghost Plugin Example")

ghost.extend("./plugins/greeter.so")

// The plugin registers `greet` only once extend() runs, so it can only be
// imported after that point - importing it up front would find nothing yet.
import "ghost:greet"

greet()