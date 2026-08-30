// WIP: not complete

import "ghost:random"
import "ghost:ghost"

class Ada {
    knowledge = []

    reflections = {
        "am": "are",
        "was": "were",
        "i": "you",
        "i'd": "you would",
        "i've": "you have",
        "i'll": "you will",
        "my": "your",
        "are": "am",
        "you've": "I have",
        "you'll": "I will",
        "you're": "I am",
        "your": "my",
        "yours": "mine",
        "you": "I",
        "me": "you"
    }

    function load(knowledge) {
        this.knowledge = knowledge

        return this
    }

    function respond(text = "") {
        console.log("responding to: " + text + "...")
        foundMatch = false

        // Loop through our collected knowledge base.
        //
        // NOTE: findAll() now answers every match's full text (§13.7), not a
        // match's capture groups - the {1}-style placeholders below rely on
        // capture groups Ghost's string API no longer exposes any way to
        // read, so reflection substitution here is a known gap pending a
        // capture-group-aware method.
        for (knowledge in this.knowledge) {
            if (text.toLowerCase().matches(knowledge.pattern) and foundMatch == false) {
                foundMatch = true
                matches    = text.toLowerCase().findAll(knowledge.pattern)
                response   = knowledge.responses[random.range(knowledge.responses.length())]

                for (index, match in matches) {
                    response = response.replace("{" + index + "}", this.reflect(match))
                }

                console.log("== " + response)
            }
        }
    }

    function reflect(text) {
        words = text.split(" ")

        for (index, word in words) {
            if (this.reflections[word]) {
                words[index] = this.reflections[word]
            }
        }

        return words.join(" ")
    }
}

version = "0.3"

console.log("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
console.log("  Ada:   v" + version)
console.log("  Ghost: v" + ghost.version)
console.log()
console.log("  Ada is a rudimentary AI based on ELIZA.")
console.log("  Speak with Ada in plain English, and they will reply.")
console.log("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
console.log()

ada = Ada()

ada.load([
    "one",
    "two",
    "three"
]).respond()

// ada.respond()