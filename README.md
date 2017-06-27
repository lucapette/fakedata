# fakedata

`fakedata` is a small program that generates random data on the command line.

# Table Of Contents

- [Overview](#overview)
  - [Quick Start](#quick-start)
  - [Why another random data generator?](#why-another-random-data-generator)
- [Generators](#generators)
  -[Constraints](#constraints)
- [Templates](#templates)
- [How to install](#how-to-install)
- [How to contribute](#how-to-contribute)
- [Code of conduct](#code-of-conduct)

# Overview

## Quick Start

`fakedata` helps you generate random data in various ways. You can generate data
by specifying on the command line the kind of data you need:

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
Be default, `fakedata` generates data using a space as a separator. You can
choose a different output format like CSV:

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

or SQL insert statements:

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

If you need some control over the output, you can use templates:

```sh
$ echo '{{Email}}--{{Int}}--{{Color}}' | fakedata -l5
antonyzotov@test.george--967--azure
Skyhartman@test.xn--mgbaam7a8h--238--cyan
syropian@example.jlc--566--black
catadeleon@example.sohu--60--white
kennyadr@test.best--899--red
```

## Why another random data generator?

`fakedata` focuses on a simple UI (if you think it could be simpler, please [let
us know!](https://github.com/lucapette/fakedata/issues/new) We :heart:
feedback!) and the ability to fully control both the output format (using
[templates](#templates)) and the set of values a generator will pick from. We
call this feature "generators' constraints" and it's explained in detail
[here](#constraints).

# Generators

`fakedata` provides a number of generators. You can see the full list running
the following command:

```sh
$ fakedata --generators # or -G
color             one word color
country           Full country name
country.code      2-digit country code
date              date
domain            domain
domain.tld        example|test
# ...
# It's a long list :)
```

You can use the `-g` (or `--generator`) option to see an example:

```sh
$ fakedata -g sentence
Description: sentence

Example:

Jerk the dart from the cork target.
Drop the ashes on the worn old rug.
The sense of smell is better than that of touch.
Tin cans are absent from store shelves.
Shut the hatch before the waves push it in.
```

## Constraints

Some generators allow you to pass in a range to constraint the output to a
subset of values. To find out which generators support constraints:

```sh
$ fakedata --constraints

```

```sh
$ fakedata int:1,100 # will generate only integers between 1 and 100
$ fakedata int:50, # specifying only min number works too
$ fakedata int:50 # also works
```

The `enum` generator allows you to specify a set of values. It comes handy when
you need random data from a small set of values:

```sh
$ fakedata --limit 5 enum
foo
baz
foo
foo
baz
$ fakedata --limit 5 enum:bug,feature,question,duplicate
question
duplicate
duplicate
bug
feature
```

When passing a single value `enum` can be used to repeat a value in every line:

```sh
$ fakedata --limit 5 enum:one,two enum,repeat
two repeat
one repeat
two repeat
one repeat
one repeat
```

The `file` generator can be use to read custom values from a file:

```sh
$ printf "one\ntwo\nthree" > values.txt
$ fakedata -l5 file:values.txt
three
two
two
one
two
```

# Templates

`fakedata` supports parsing and executing template files for generating
customized output formats. `fakedata` executes the provided template a number of
times based on the limit flag (`-l`, `--limit`) and writes the output to
`stdout`, exactly like using inline generators.

`fakedata` can read templates from disk:

```sh
$ echo "{{Email}}--{{Int}}" > /tmp/template.tmpl
$ fakedata --template /tmp/template.tmpl
ademilter@test.school--214
Silveredge9@example.anquan--379
plbabin@example.here--902
silvanmuhlemann@test.aero--412
ivanfilipovbg@test.bmw--517
robbschiller@example.feedback--471
rickdt@example.vista--963
rmlewisuk@test.info--101
linux29@example.archi--453
g3d@test.pl--921
```

Or you can pipe the template into `fakedata`:

```sh
$ echo "#{{ Int 0 100}} {{ Name }} <{{ Email }}>" | fakedata
#56 Dannie Martin <bassamology@test.th>
#89 Moshe Walsh <baires@example.autos>
#48 Buck Reid <syropian@test.cg>
#55 Rico Powell <findingjenny@example.pohl>
#92 Luise Wood <91bilal@example.link>
#30 Isreal Henderson <thierrykoblentz@test.scb>
#96 Josphine Patton <abelcabans@test.wtf>
#95 Jetta Blair <tgerken@example.jewelry>
#10 Clorinda Parsons <roybarberuk@test.gives>
#0 Dionna Bates <jefffis@test.flights>
```

### Generators

All the generators listed under `fakedata -g` are available as functions into
the templates. If the generator name is a single word, then it's available as a
function with the same name capitalized (example: `int` becomes `Int`). If the
generator name is composed by multiple words joined by dots, then the function
name is again capitalized by the first letter of the word and joined together
(example: `product.name` becomes `Product.Name`).

`fakedata` provides the following custom generators:

- `Enum` 
- `File`
- `Int`
- `Date`

### `Enum`

Enum takes one or more strings and returns a random string on each run. Strings
are passed to Enum like so:

```html
{{ Enum "feature" "bug" "documentation" }}
```

This Enum will return either the string `feature`, `bug`, or `documentation` for
each run.

### `File`

File reads a file from disk and returns a random line on each run. It takes one
parameter which is the path to the file on disk.

```
{{ File "/var/data/dummy/dummy.txt" }}
```

### `Int`

Int takes one or two integer values and returns a number within this range. By
default it returns a number between `0` and `1000`.

```
echo "{{ Int 15 20 }}" | fakedata -l5
15
20
15
15
17
```

### `Date`

Date takes one or two dates and returns a date within this range. By default, it
returns a date between one year ago and today.

### Helpers

Beside the generator functions, the `fakedata` template implementation provides
a set of helper functions:

- `Loop`
- `Odd`
- `Even`

When using a custom loop make sure to use `--limit 1`, otherwise the loop will
run multiple times! Running a template with `{{ range Loop 5}}` and `--limit 5`
will execute 25 times.

If you need to create your own loop for advanced templates you can use the `{{
Loop }}` function. This function takes a single integer as parameter which is
the number of iterations. `Loop` has to be used with `range` e.g.

```html
{{ range Loop 10 }}
  I am printed 10 times!
{{ end }}
```

In combination with `Loop` and `range` you can use `Odd` and `Even` to determine
if the current iteration is odd or even. This is especially helpful when
creating HTML tables:

```html
{{ range $i, $j := Loop 5 }}
<tr>
  {{ if Odd $i -}}
  <td class="odd">
    {{- else -}}
  <td class="even">
    {{- end -}}
    {{ Name }}
  </td>
</tr>
{{ end }}
```

By using `Odd` we can create tables with a class name of  `odd` and `even` when
generating our HTML. Odd takes an integer as parameter which is why we need to
assign the return values of `Loop 5` to the variables `$i` and `$j`.

Beside the helper function `Loop`, `Odd`, and `Even` templates also support
manipulation with `printf`. By using `printf` we can create a custom output, for
example to display a full name in the format `Lastname Firstname` instead of
`Firstname Lastname`.

```html
{{ printf "%s %s" NameLast NameFirst }}
```

# How to install

## Homebrew

`fakedata` can be installed through Homebrew:

``` sh
$ brew install lucapette/tap/fakedata
```

## Standalone

`fakedata` can be installed as an executable. Download the latest [compiled
binary](https://github.com/lucapette/fakedata/releases) and put it anywhere in
your executable path.

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
- Join the [#fakedata](https://gophers.slack.com/messages/fakedata/) channel in
the [Gophers Slack](https://invite.slack.golangbridge.org/).

If you're not sure where to start, please open a [new
issue](https://github.com/lucapette/fakedata/issues/new) and we'll gladly help
you get started.

# Code of Conduct

You are expected to follow our [code of conduct](/CODE_OF_CONDUCT.md) when
interacting with the project via issues, pull requests or in any other form.
Many thanks to the awesome [contributor covenant](http://contributor-covenant.org/) initiative!

# License

[MIT License](/LICENSE) Copyright (c) [2017] [Luca Pette](http://lucapette.me)
