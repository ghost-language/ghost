import "ghost:random"

class Ada {
    knowledge = []
    reflections = {}

    function constructor() {
        this.reflections = {
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
    }

    function greet() {
        this.respond('hello')
    }

    function load(module) {
        this.knowledge = module

        return this
    }

    function respond(input) {
        foundMatch = false

        // Loop through the knowledge base
        //
        // NOTE: findAll() now answers every match's full text (§13.7), not a
        // match's capture groups - the {1}-style placeholders below rely on
        // capture groups Ghost's string API no longer exposes any way to
        // read, so reflection substitution here is a known gap pending a
        // capture-group-aware method.
        for (knowledge in this.knowledge) {
            if (input.toLowerCase().matches(knowledge.pattern) and foundMatch == false) {
                foundMatch = true
                matches = input.toLowerCase().findAll(knowledge.pattern)
                response = knowledge.responses[random.range(knowledge.responses.length())]

                for (index, match in matches) {
                    response = response.replace("{%s}".format(index), this.reflect(match))
                }

                console.log("== %s".format(response))
            }
        }
    }

    /**
     * Reflect on the user's input.
     */
    function reflect(input) {
        words = input.split(" ")

        for (index, word in words) {
            if (this.reflections[word]) {
                words[index] = this.reflections[word]
            }
        }

        return words.join(" ")
    }
}