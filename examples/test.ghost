function expect(value) {
  return {
    toBe: function(expected) {
      if (value != expected) {
        // console.log("Expected failed")
        return false
      } else {
        // console.log("Passed")
        return true
      }
    }
  }
}

function test(description, callback) {
    result = callback()

    if (!result) {
        console.log(description + " ...FAILED")
        return
    }

    console.log(description + " ...PASSED")
}

test("example test A", function() {
  return expect(1).toBe(2)
})

test("example test B", function() {
  return expect(1).toBe(1)
})