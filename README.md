# fakedata

CLI utility that generates data in various format.

# Overview

`fakedata` is a small utility that generates data from the command line:

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

List the available generators:

```sh
$ fakedata --generators
color
country
country.code
domain
domain.name
domain.tld
double
email
event.action
http.method
id
ipv4
ipv6
latitude
longitude
mac.address
name
name.first
name.last
product.category
product.name
state
state.code
timezone
unixtime
username
```

## SQL formatter

`fakedata` can generate insert statements. By default, it uses the name of the generators 
as column names:

```sh
$ fakedata email domain --format=sql --limit 1
INSERT INTO TABLE (email,domain) values ('yigitpinarbasi@example.org','example.me');
```

You can specify the name of the column using a field with the following format `column_name=generator`.
For example:

```sh
$ fakedata login=email referral=domain --format=sql --limit 1
INSERT INTO TABLE (login,referral) values ('calebogden@example.com','test.me');
```

# Installation guide

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
