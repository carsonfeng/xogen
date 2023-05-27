# xogen 1.0

## Overview

xogen is a code generator that can convert models from the "github.com/go-xorm" project into dao code. This makes it easy to generate dao crud code that corresponds to the model.

## Feature List

- Converts models from "github.com/go-xorm" into dao code
- Generates dao crud code easily

## Development

### Prerequisites

- Go 1.16 or later

### Installing

```
go get github.com/example/xogen
```

### Building

```
cd $GOPATH/src/github.com/example/xogen
go build
```

## Release

### Version 1.0

First release of xogen.

## Usage

```
xogen [options]
```

### Options

- `-input` specifies the input directory for models
- `-output` specifies the output directory for dao code
- `-pkg` specifies the name of the package for dao code
- `-url` specifies the database connection url
- `-table` specifies the database table name for models