# fakedata

`fakedata` is a small command line utility that generates random data.

# Table Of Contents

- [Overview](#overview)
- [Generators](#generators)
- [Formatters](#formatters)
- [How to install](#how-to-install)
- [How to contribute](#how-to-contribute)
- [Code of conduct](#code-of-conduct)

# Overview

Here is a list of examples to get a feeling of how `fakedata` works

```sh
$ fakedata email country
cemshid@example.com Afghanistan
LucasPerdidao@example.me Turkey
arthurholcombe1@test.us Saint Helena
iamgarth@example.us Montenegro
joelcipriano@test.name Croatia
keryilmaz@test.name Vietnam
plbabin@test.org Lithuania
bermonpainter@test.us Haiti
opnsrce@example.name Malaysia
ankitind@test.info Virgin Islands, British
```

Limit the number of rows:

```sh
$ fakedata --limit 5 country.code
SH
CF
GQ
PE
FO
```

Choose a different output format:

```sh
$ fakedata --format=csv product.category product.name
Shoes,Rankfix
Automotive,Namis
Movies,Matquadfax
Tools,Damlight
Computers,Silverlux
Industrial,Matquadfax
Home,Sil-Home
Health,Toughwarm
Shoes,Freetop
Tools,Domnix
```

# Generators

`fakedata` provides a number of generators. You can see the full list running
the following command:

```sh
$ fakedata --generators
color             one word color
country           Full country name
country.code      2-digit country code
date              date
domain            domain
domain.tld        example|test
# ...
# It's a long list :)
```

Some generators allow you to pass in a range to constraint the output to a
subset of values:

```sh
$ fakedata int,1..100 # will generate only integers between 1 and 100
$ fakedata int,50.. # specifying only min number works too
$ fakedata int,50 # also works
```

The `enum` generator allows you to specify a set of values. It comes handy when
you need random data from a small subset of values:

```sh
$ fakedata enum --limit 5
foo
baz
foo
foo
baz
$ fakedata enum,bug..feature..question..duplicate --limit 5
question
duplicate
duplicate
bug
feature
```

# Formatters

### SQL formatter

`fakedata` can generate insert statements. By default, it uses the name of the
generators as column names:

```sh
$ fakedata --format=sql --limit 1 email domain
INSERT INTO TABLE (email,domain) values ('yigitpinar@example.org','example.me');
```

You can specify the name of the column using a field with the following format 
`column_name=generator`:

```sh
$ fakedata --format=sql --limit 1 login=email referral=domain
INSERT INTO TABLE (login,referral) values ('calebogden@example.com','test.me');
```

```sh
fakedata --format=sql --limit=1 --table=users login=email referral=domain
INSERT INTO users (login,referral) VALUES ('mikema@example.com' 'test.us');
```

# How to install

## Homebrew

`fakedata` can be installed through Homebrew:

``` sh
$ brew tap lucapette/tap
$ brew install fakedata
```

## Standalone

`fakedata` can be installed as an executable. Download the latest
[compiled binaries](https://github.com/lucapette/fakedata/releases) and put it
anywhere in your executable path.

## Source

Please refer to our [contributing guidelines](/CONTRIBUTING.md) to build and
install `fakedata` from the source.

# How to contribute

We love every form of contribution! Good entry points to the project are:

- Our [contributing guidelines](/CONTRIBUTING.md) document
- Issues with the tag
  [gardening](https://github.com/lucapette/fakedata/issues?q=is%3Aissue+is%3Aopen+label%3Agardening)
- Issues with the tag [good first
  patch](https://github.com/lucapette/fakedata/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+patch%22)

If you're still not sure where to start, please open a [new
issue](https://github.com/lucapette/fakedata/issues/new) and we'll gladly
help you get started.

# Code of Conduct

You are expected to follow our [code of conduct](/CODE_OF_CONDUCT.md) when
interacting with the project via issues, pull requests or in any other form.
Many thanks to the awesome [contributor covenant](http://contributor-covenant.org/) initiative!

# License

[MIT License](/LICENSE) Copyright (c) [2017] [Luca Pette](http://lucapette.me)
