// The math module works on one number at a time, and on whole lists of them.

import "ghost:math"

// Scalars.
console.log(math.sqrt(16))
console.log(math.pow(2, 10))
console.log(math.round(2.567, 2))
console.log(math.log(8, 2))
console.log(math.degrees(math.pi))
console.log(math.clamp(15, 0, 10))
console.log(math.isClose(0.1 + 0.2, 0.3))

// Arithmetic on lists is elementwise, and the operators say so more clearly
// than a method call does.
console.log([1, 2, 3] + 10)
console.log([1, 2, 3] * 2)
console.log([1, 2, 3] + [10, 20, 30])
console.log([1, 2] == [1, 2])

// Shapes line up from the right, so a row spreads down a matrix. This is what
// lets one bias row apply to a whole batch of samples.
console.log([[1, 2], [3, 4]] + [10, 20])

// The methods are the same operation under a name, for where a name reads
// better - and every elementwise method broadcasts the same way.
console.log(math.add([1, 2, 3], 10))
console.log(math.sqrt([1, 4, 9]))
console.log(math.sqrt([[1, 4], [9, 16]]))

// Joining is a different operation, so it has a method rather than an operator.
console.log([1, 2].concat([3, 4]))

// Statistics read their values however they are held: spread across the call,
// collected in a list, or arranged as a matrix.
scores = [72, 85, 90, 90, 61, 78]

console.log(math.mean(scores))
console.log(math.median(scores))
console.log(math.mode(scores))
console.log(math.standardDeviation(scores))
console.log(math.percentile(scores, 90))
console.log(math.min(scores), math.max(scores))
console.log(math.sort(scores))

// Building lists of numbers.
console.log(math.arange(0, 1, 0.25))
console.log(math.linspace(0, 1, 5))
console.log(math.zeros(3))
console.log(math.identity(3))
console.log(math.reshape(math.arange(6), 2, 3))

// Vectors and matrices are lists and lists of lists.
console.log(math.dot([1, 2, 3], [4, 5, 6]))
console.log(math.norm([3, 4]))
console.log(math.normalize([3, 4]))
console.log(math.cross([1, 0, 0], [0, 1, 0]))
console.log(math.outer([1, 2, 3], [10, 20]))
console.log(math.distance(0, 0, 3, 4))

matrix = [[4, 7], [2, 6]]

console.log(math.transpose(matrix))
console.log(math.determinant(matrix))
console.log(math.inverse(matrix))

// 2x + y = 5, x + 3y = 10
console.log(math.solve([[2, 1], [1, 3]], [5, 10]))

// Whole-number mathematics.
console.log(math.gcd(12, 18))
console.log(math.factorial(10))
console.log(math.isPrime(97))
console.log(math.combinations(52, 5))

// Random numbers share a seed with the random module, so a seeded run replays
// exactly. Noise drifts where random numbers jitter.
math.randomSeed(42)

console.log(math.randomInt(1, 6))
console.log(math.noise(0.5, 1.5))
