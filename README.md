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
go get github.com/carsonfeng/xogen
```

### Building

```
cd $GOPATH/src/github.com/carsonfeng/xogen
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

- `-model_path` specifies the path of model. For Example: `testcases/models/models.go`
- `-models` specifies the xorm model struct names, separated by comma ',' . For Example: `DemoTeamModel1,DemoUserModel2`
- `-output_dao` specifies the path of dao output. For Example: `testcases/daos/demo_daos.go`
- `-model_import` specifies imported package of model. For Example: `xogen/testcases/models`

### Integrated with your projects
Step 1: Build a tools project. Follow the examples of `main.go` and `./testcases/templated.go` in this project.