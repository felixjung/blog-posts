# JS to Go: Writing Basic Tests in Go

> I have been working with Go in a service-oriented architecture over the past year and a half. It has been a fun ride after years of working with and enjoying JavaScript. This post is part of a series in which I attempt to make my learnings accessible to other engineers with a JavaScript background. Huge thanks go to my dear colleague [Tom Arrell](https://tomarrell.com "Tom's website") who has taught me most of the things I now know about Go.

Tests are one of the most important things when shipping software users depend on. When making fast incremental changes in a large project, tests and good test coverage allow us to trust we did not break anything. This is especially true when working in an unfamiliar codebase where tests also help us understand how a certain part of the code works.

TypeScript and Go each have a static type system which protects us from passing incompatible parameters to functions or accessing object properties in an unsafe way. Having that extra safety compared to JavaScript is great, but we will still want to cover our application or library code in tests. In this blog post we will explore how to do that in Go and we will look at some patterns for different test scenarios and contexts. I assume you test your JavaScript code with [Jest](https://jestjs.io "Jest Testing Framework Website").

## Writing a Basic Function Test

Let’s start by testing a simple `sum` function which takes two numbers and returns their sum.[^1] The JavaScript version of that function looks like as follows.

```js
function sum(a, b) {
  return a + b
}
```

In our test we set up the two inputs `a` and `b`, write down the expected return value value (`want`) and run our assertion using Jest’s `expect` [function](https://jestjs.io/docs/en/expect "Jest documentation for the expect function") and the `toBe` [matcher](https://jestjs.io/docs/en/using-matchers "Jest documentation for common matchers"). The test file, `sum.spec.js`, is colocated with the module, `sum.js`.

```js
import { sum } from './sum';

test('sum', () => {
  const a = 1;
  const b = 2;

  const want = 3;

  const got = sum(a, b)

  expect(got).toBe(want);
});
```

The Go equivalent to the above example looks like this.

```go
package math

func Sum(a, b int) int {
	return a + b
}
```

In a `math` package we export a `Sum` function from a `sum.go` file. Right next to that file we place our test file, `sum_test.go` with the following test code.

```go
package math

import "testing"

func TestSum(t *testing.T) {
	a := 1
	b := 2

	want := 3

	got := Sum(a, b)

	if want != got {
		t.Fatalf("Sum: wanted %d, got %d", want, got)
	}
}
```

Our testing approach is the same as in the JavaScript example. So let us look at the differences in [how tests work in go](https://golang.org/pkg/testing/ "Go documentation on testing"). The biggest difference is that Go has out of the box-support for testing. You do not need a testing framework like Jest.

> Testing frameworks for Go do exist. If you look for them you will most likely stumble upon [Ginkgo](https://onsi.github.io/ginkgo/ "Ginkgo testing framework website") or [GoConvey](http://goconvey.co "GoConvey testing framework website"). Both these examples aim at providing a Behavior Driven Development (BDD) testing experience. They  provide abstractions for writing tests and test runners with features like watch mode.
> I would recommend _against_ using such testing frameworks. As with other aspects of writing Go I feel that sticking close to the language and its tooling keeps things simple. If you want to write BDD style tests, there are ways of doing that with Go’s testing support.

Tests are run with the `go test` command and package `testing` provides all the basics you need to test your code. Additionally, the example demonstrates the following.

* Tests live in a file with the `_test` suffix in the filename. When you run `go test` the test tool will look for files with that suffix.
* Every test is defined with a function using the `Test` prefix in the function name, i.e. `TestSum`. It is important that the first letter after `Test` is capitalized. Functions located in a test file and using this name pattern are identified as test routine.
* A test routine is called with a pointer to an instance of [the `T` struct from the `testing` package](https://pkg.go.dev/testing#T "Documentation for the T type of package testing") (i.e. `*testing.T`). `T` manages the test state and provides functions for failing a test or writing logs.
* Assertions can be written with plain Go. You make your assertion in a simple `if` statement. You fail the test by calling `t.Fail`, `t.FailNow`, `t.Fatal`, `t.Fatalf`, or other methods on the `T` struct. This becomes a bit tedious for more complicated assertions and will look at assertion libraries in the next section.

## Powerful Assertions With Testify

Making assertions with `if` statements can get a bit harder to read and may become quite involved when you start handling more complicated assertions. For example, you may want to test whether a slice contains a certain value or that an HTTP handler returns a certain status code. [Testify](https://github.com/stretchr/testify "Testify GitHub repository") is a very popular library for mocking and assertions. It is designed to work well with the standard library. I would recommend to use the official [GoMock framework](https://github.com/golang/mock "GoMock GitHub repository") over using Testify’s mocking feature.

Let us rewrite the assertion in `TestSum` with Testify.

```go
package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSumWithTestify(t *testing.T) {
	a := 1
	b := 2

	want := 3

	got := Sum(a, b)

	require.Equal(t, want, got)
}
```

Here we are using Testify’s `require` package and the `Equal` matcher. Testify also has an `assert` package with the same API. However, the matchers in `require` will fail your test immediately while `assert`’s matchers will allow your test to continue. I prefer my tests to fail immediately because I find easier to identify the problem causing my test to fail. The first argument to a matcher is the `testing.T` pointer. Matchers comparing an expected and an actual value take the expected value as second parameter and the actual value as third parameter. This is important because you can get easily confused by a failing test’s output when you flipped the two.

If you have a test with several assertions, you can create an instance of the `require` package `Assertions` and use that directly.

```go
package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSumWithAssertionsInstance(t *testing.T) {
	require := require.New(t)

	a := 1
	b := 2

	want := 3

	got := Sum(a, b)

	require.Equal(want, got)
}
```

Some linter rules may complain about shadowing the `require` package inside of your test function. Do not bother with those violations, this is common practice.

## Behavior Driven Development With Subtests

Earlier in this post I mentioned how some Go testing frameworks are made for [Behavior Driven Development (BDD)](https://en.wikipedia.org/wiki/Behavior-driven_development "BDD on Wikipedia"). In BDD you develop your software and tests against human-readable requirements. These requirements can for example be written by product managers using a product domain’s [ubiquitous language](https://martinfowler.com/bliki/UbiquitousLanguage.html "Martin Fowler definition of Ubiquitous Language"). 

To make the example a bit more interesting, let us consider a `sumOrMax` function, which returns the sum of parameters `a` and `b` as long as it stays below a `max` value, which is provided as a third parameter.

```js
export function sumOrMax(a, b, max) {
  const sum = a + b;

  return sum < max ? sum : max;
}
```

Jest supports BDD tests out of the box with the `describe` and `it` functions.

```js
import { sumOrMax } from './sum-or-max';

describe('sumOrMax', () => {
  describe('when the sum is below the maximum', () => {
    it('should return the sum', () => {
      const a = 1;
      const b = 2;
      const max = 4;

      const want = 3;

      const have = sumOrMax(a, b, max);

      expect(have).toBe(want);
    });
  });

  describe('when the sum is at or above the maximum', () => {
    it('should return the maximum', () => {
      const a = 5;
      const b = 10;
      const max = 4;

      const want = max;

      const have = sumOrMax(a, b, max);

      expect(have).toBe(want);
    });
  });
});
```

In Go we can achieve something equivalent with subtests. The Go version of `sumOrMax` looks as follows.

```go
package math

func SumOrMax(a, b, max int) int {
	sum := a + b

	if sum < max {
		return sum
	}

	return max
}
```

To write our BDD-style test, we take advantage of the `Run` method exposed by the `T` struct. It allows us to spawn a new test, a subtest. The behavior is similar to how `it` behaves in Jest.

```go
package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSumOrMax(t *testing.T) {
	t.Run("when the sum is below the maximum", func(t *testing.T) {
		a := 1
		b := 2
		max := 4

		want := 3

		have := SumOrMax(a, b, max)

		require.Equal(t, want, have, "should return the sum")
	})

	t.Run("when the sum is at or above the maximum", func(t *testing.T) {
		a := 5
		b := 10
		max := 4

		want := max

		have := SumOrMax(a, b, max)

		require.Equal(t, want, have, "should return the maximum")
	})
}
```

For every `describe` block in the JavaScript test suite, we introduce a call to `t.Run`. For the description of the `it` blocks we pass a description of the expected behavior to the assertion call. Most of Testify’s assertions support that extra `msgAndArgs` parameter.

## Testing Many Scenarios With Table Driven Tests

Especially in unit tests we want to test a given function for many different inputs. In such cases table driven tests are an efficient approach to test all scenarios in a single test.

Let us try table driven tests with a `fib` function, which returns the nth number in the [Fibonacci series](https://en.wikipedia.org/wiki/Fibonacci_number "Wikipedia page for the Fibonacci number").[^2] Our JavaScript implementation looks as follows.

```js
export function fib(n) {
  if (n < 2) {
    return n;
  }

  return fib(n - 1) + fib(n - 2);
}
```

In [Jest we implement a table driven](https://jestjs.io/docs/en/api#testeachtablename-fn-timeout "Jest documentation on table driven tests") test as follows.

```js
import { fib } from './fib';

describe('fib', () => {
  it.each([
    [0, 0],
    [1, 1],
    [2, 1],
    [3, 2],
    [4, 3],
    [5, 5],
    [6, 8],
    [7, 13],
    [8, 21],
  ])('fib(%i) should be %i', (n, want) => {
    const have = fib(n);

    expect(have).toBe(want);
  });
});
```

In case you are not familiar with how table tests in Jest work, I will walk you through it.

1. In our call to `each()` we pass the _test table_. This is a nested array where every element is a test case. The test case consists of function parameters and the expected value. You may order these as you like as the values are just forwarded in order.
2. In the chained call to the return value of `each()` we pass a description of the test and a function to run the test. The description is a template string, where the test case values can be rendered using `printf` format. The test function receives the values of our test case as arguments.
3. In our test function we run assertions just like we would in a normal `it()` test.

Every test case will be run by Jest and we get a nice report of the successes and failures.

The Go implementation, `Fib` looks as follows.

```go
package math

func Fib(n int) int {
	if n < 2 {
		return n
	}

	return Fib(n-1) + Fib(n-2)
}
```

To create the same table driven test as we did for the JavaScript version we will use `t.Run`, just like we did in our subtests.

```go
package math

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFib(t *testing.T) {
	require := require.New(t)

	tc := []struct {
		n    int
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
		{8, 21},
	}

	for _, tt := range tc {
		name := fmt.Sprintf("fib(%d) should be %d", tt.n, tt.want)

		t.Run(name, func(t *testing.T) {
			require.Equal(tt.want, Fib(tt.n))
		})
	}
}
```

There is no test framework magic involved here.

1. We create a table for our test cases, `tc`, as a slice using an inline struct type for the values.
2. We iterate over the test cases in a [`for` loop using the range expression](https://gobyexample.com/range "Go by Example documentation page for range").
3. We create a name for the subtest and run the test using testify for assertions.

And that is it. You have successfully tested your code.

## Conclusion

In this post we have learned how to write tests for simple Go functions. These functions were simple in the sense that they were pure: the same inputs lead to the same output with no other side effects.

Testing a function is straight forward with Go’s standard library — you just write Go, no frameworks needed. Libraries like testify help a lot by reducing the effort you have to put into your assertions. Subtests allow you to structure your test cases in a BDD style manner. Finally, table driven tests make it easy to handle a lot of test cases with very little code.

Things will get a bit more difficult when you need to deal with dependencies, such as databases. For example, if your function queries a database you will not necessarily want to have that database running for a unit test. In a future post we will look at how interfaces and mocks can help us with this problem.

[^1]:	Incidentally, this is the same example the Jest folks [use on their website](https://jestjs.io/docs/en/getting-started "Go to the Jest Getting Started guide"). 🤷‍♂️

[^2]:	This example is shamelessly taken from [Dave Cheney’s very helpful and always interesting blog](https://dave.cheney.net/2013/06/09/writing-table-driven-tests-in-go "Dave Cheney — Writing table driven tests in Go").