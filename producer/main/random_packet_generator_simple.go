package main

import "flow-generator/producer"

func main() {
	sfg := producer.SimpleFlowGenerator{}
	sfg.Init()
	sfg.Run()
}
