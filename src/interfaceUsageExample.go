package main

import "fmt"

type I interface {
	M()
}

// нет явного указания implements или типа того, как в Java
type T struct {
	S string
}

// просто реализовали метод из интерфейса
func (t T) M() {
	fmt.Println(t.S)
}

func main() {
	var i I = T{"hello"}
	i.M()
}
