# 7 GUIs

This demonstrates several classic GUI problems based on [7GUIs](https://eugenkiss.github.io/7guis/).

The examples are over-commented to help understand the structure better, in practice, you don't need that many comments.

## Counter

Counter shows basic usage of Gio and how to write interactions.

It displays a count value that increases when you press a button.

[Source](./counter/main.go)

## Temperature Converter

Temperature conversion shows bidirectional data flow between two editable fields.

It implements a bordered field that can be used to propagate values back to another field without causing update loops.

[Source](./temperature/main.go)