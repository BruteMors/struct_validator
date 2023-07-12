package main

import "fmt"

type Person struct {
	Age  []int  `validate:"min:10|max:30"`
	Name string `validate:"len:3|in:foo,bar"`
}

func main() {
	p := Person{
		Age:  []int{22, 11},
		Name: "foo",
	}
	err, valid := Validate(p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	for _, validationError := range valid {
		fmt.Println(validationError.Field, validationError.Err)
	}
}
