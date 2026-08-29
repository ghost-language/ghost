console.log("starting...")

if (false) {
    console.log(true)
} else {
    console.log("test 1/6: correct")
}

if (true) {
    console.log("test 2/6: correct")
} else {}

if (true) {
    console.log("test 3/6: correct")

    if (1 + 1 == 2) {
        console.log("test 4/6: correct")

        if (true) {
            console.log("test 5/6: correct")
        } else {
            console.log("text 5/6: incorrect")
        }
    } else {
        console.log("text 4/6: incorrect")
    }
}

console.log("test 6/6: correct")
// console.log(1 + 2 == 3)
console.log("completed.")