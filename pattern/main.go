package main

import "pattern/singleton"

func main() {
	ins := singleton.GetInstance()
	ins.Work()

}
