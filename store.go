package main

import "fmt"

type Store struct {
	data map[string]string
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) Set(k string, v string) {
	fmt.Println("setting", k, v)
	s.data[k] = v
}

func (s *Store) Get(k string) string {
	fmt.Println("retrieving", k)
	return s.data[k]
}
