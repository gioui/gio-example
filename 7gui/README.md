# 7 GUIs

This demonstrates several classic GUI problems based on [7GUIs](https://eugenkiss.github.io/7guis/). They show different ways of using Gio framework.

The examples show one way of implementing of things, of course, there are many more.

The examples are over-commented to help understand the structure better, in practice, you don't need that many comments.

## Counter

Counter shows basic usage of Gio and how to write interactions.

It displays a count value that increases when you press a button.

[UI](./counter/main.go)

## Temperature Converter

Temperature conversion shows bidirectional data flow between two editable fields.

It implements a bordered field that can be used to propagate values back to another field without causing update loops.

[UI](./temperature/main.go)


## Timer

Timer shows how to react to external signals.

It implements a timer that is running in a separate goroutine and the UI interacts with it. The same effect can be implemented in shorter ways without goroutines, however it nicely demonstrates how you would interact with information that comes in asynchronously.

The UI shows a slider to change the duration of the timer and there is a button to reset the counter.

[UI](./timer/main.go), [Timer](./timer/timer.go)