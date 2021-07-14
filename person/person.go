package person

import (
	"fmt"
	"time"
)

type Person struct {
	name string
	age  int
}

func (p *Person) SetDeatails(name string, age int) {
	p.name = name
	p.age = age
	fmt.Println("SetDetials hong", p)
}

func (p Person) Name() string {
	return p.name
}

func (p Person) SayHello() {
	fmt.Printf("Hello my name is %s and I'm %d\n", p.name, p.korenAge())
}

func (p Person) korenAge() int {
	return time.Now().Year() - p.age + 1
}
