# CLI_Expense_Tracker

DB scheme:
CREATE TABLE public.Expenses (
  Id  SERIAL PRIMARY KEY,
  Date TIMESTAMP NOT null,
  Description VARCHAR(30) NOT NULL,
  Amount INT NOT NULL
);

LINTER SETTINGS:
# This file contains all available configuration options
# with their default values.

# options for analysis running
run:
  timeout: 15m
  tests: true

# output configuration options
output:
  formats:
    - format: tab

# all available settings of specific linters
linters-settings:
  unused:
    field-writes-are-uses: false
  exhaustive:
    default-signifies-exhaustive: true

linters:
  enable:
    # mandatory linters
    - govet
    - revive

    # some default golangci-lint linters
    - errcheck
    - gosimple
    - ineffassign
    - staticcheck
    - typecheck
    - unused

    # extra linters
    - exhaustive
    - godot
    - gofmt
    - whitespace
    - goimports
  disable-all: true
  fast: false

issues:
  include:
    - EXC0002 # should have a comment
    - EXC0003 # test/Test ... consider calling this
    - EXC0004 # govet
    - EXC0005 # C-style breaks


