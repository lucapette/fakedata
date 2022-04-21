# Contributing

We love every form of contribution. By participating to this project, you agree
to abide to the `fakedata` [code of conduct](/CODE_OF_CONDUCT.md).

## Setup your machine

Our [Makefile](/Makefile) is the entry point for most of the activities you will
run into as a contributor. To get a basic understanding of what you can do with
it, you can run:

```sh
$ make help
```

Which shows all the documented targets. `fakedata` is written in
[Go](https://golang.org/). Here is a list of prerequisites to build and test the
code:

- `make`
- [Go 1.16+](http://golang.org/doc/install)

Clone `fakedata` from source:

```sh
$ git clone https://github.com/lucapette/fakedata.git
$ cd fakedata
```

A good way of making sure everything is all right is running the test suite:

```sh
$ make test
```

Please open an [issue](https://github.com/lucapette/fakedata/issues/new) if you
run into any problem.

## Building and running fakedata

You can build the entire application by running `make` without arguments:

```sh
make
```

since `build` is the default target.

You can run `fakedata` following the steps:

```sh
$ make
$ ./fakedata username
```

## Testing

We try to cover as much as we can with testing. The goal is having each single
feature covered by one or more tests. Adding more tests is a great way of
contributing to the project!

### Running the tests

Once you are [setup](#setup-your-machine), you can run the test suite with one
command:

```sh
$ make test
```

You can run only a subset of the tests using the `TEST_PATTERN` variable:

```sh
$ make test TEST_PATTERN=TheAnswerIsFortyTwo
```

You can use pass options to `go test` through the `TEST_OPTIONS` variable:

```sh
$ make test TEST_OPTIONS=-v
```

You can combine the two options which is very helpful if you are working on a
specific feature and want immediate feedback. Like so:

```sh
$ make test TEST_OPTIONS=-v TEST_PATTERN=TheAnswerIsFortyTwo
```

## Test your change

You can create a branch for your changes and try to build from the source as
you go:

```sh
$ make build
```

When you are satisfied with the changes, we suggest running:

```sh
$ make lint
```

This command runs the linters and the tests the same way we run them in our CI
system.

## Submit a pull request

Push your branch to your `fakedata` fork and open a pull request against the
main branch.
