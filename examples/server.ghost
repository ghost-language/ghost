import "ghost:http"
import "ghost:date"
import "ghost:file"

http.handle("/", function(request) {
    start = date.toUnixNano(date.now())
    view = file.read("views/index.html")

    console.log(view)

    console.log("method:", request["method"])
    console.log("host:", request["host"])
    console.log("content length:", request["contentLength"])
    console.log("body:", request["body"])

    end = date.toUnixNano(date.now())
    total = ((end - start) / 1e6).toString()

    console.log("-->", total, "ms")
})

http.listen(3000, function() {
    console.log("Server started at http://localhost:3000 🌱")
})