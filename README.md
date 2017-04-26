# fakedata

CLI utility that generates data in various formats.

[Overview](#overview)
  - [Quick start](#quick-start)
  - [Generators](#generators)
  - [Formatters](#formatters)
- [How to install](#how-to-install)
- [How to contribute](#how-to-contribute)
- [Code of conduct](#code-of-counduct)

# Overview

`fakedata` is a small utility that generates data from the command line:

## Quick start

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

Limit the amout of generated rows:

```sh
$ fakedata country.code --limit 5
SH
CF
GQ
PE
FO
```

Choose a different output format:

```sh
$ fakedata product.category product.name --format=csv
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

## Generators

`fakedata` provides a number of generators:

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

Some generators allow you to pass in a range so you can scope their generation
to a subset of values:

```sh
$ fakedata int,1..100 # will generate only integers between 1 and 100
$ fakedata int,50.. # specifying only min number works too
$ fakedata int,50 # also works
```

## Formatters

### SQL formatter

`fakedata` can generate insert statements. By default, it uses the name of the generators 
as column names:

```sh
$ fakedata email domain --format=sql --limit 1
INSERT INTO TABLE (email,domain) values ('yigitpinarbasi@example.org','example.me');
```

You can specify the name of the column using a field with the following format 
`column_name=generator`. For example:

```sh
$ fakedata login=email referral=domain --format=sql --limit 1
INSERT INTO TABLE (login,referral) values ('calebogden@example.com','test.me');
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

Please refer to our [contributing guidelines](/CONTRIBUTING.md).

# Code of Conduct

You are expected to follow our [code of conduct](/CODE_OF_CONDUCT.md) when
interacting with the project via issues, pull requests or in any other form.
Many thanks to the awesome [contributor covenant](http://contributor-covenant.org/) initiative!

# License

[MIT License](/LICENSE) Copyright (c) [2017] [Luca Pette](http://lucapette.me)
