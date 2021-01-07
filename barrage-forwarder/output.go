package main

type OutputFactory interface {
	Type() string
	NewOutput(config *settingOutput) Output
}

type Output interface {
	Emit(messageTime int64, data string)
}
