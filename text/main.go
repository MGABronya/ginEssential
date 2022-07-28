package main

type a struct {
	b int
}

func (A a) c() {
	A.b = 10
}

func main() {
	A := a{
		b: 8,
	}
	A.c()
	println(A)
}
