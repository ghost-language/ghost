import Ada from 'ada'
import therapist from 'modules/therapist'
import "ghost:ghost"

version = '0.3'

console.log("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
console.log("  Ada:   v%s".format(version))
console.log("  Ghost: v%s".format(ghost.version))
console.log()
console.log("  Ada is a rudimentary AI based on ELIZA.")
console.log("  Speak with Ada in plain English, and they will reply.")
console.log("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
console.log()

ada = Ada()

ada.load(therapist)
ada.greet()

while (true) {
    text = console.read("> ")

    ada.respond(text)
}