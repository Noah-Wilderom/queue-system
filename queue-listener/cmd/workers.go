package main

type Pool struct {
	Workers map[*Worker]bool
}

type Worker struct {
	Id string
}

func NewPool() *Pool {
	return &Pool{
		Workers: make(map[*Worker]bool),
	}
}
