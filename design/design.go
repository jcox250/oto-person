// Package design contains the design for the person service
package design

type PersonService interface {
	// Add adds a person
	Add(AddRequest) AddResponse
	// Show shows a person
	Show(ShowRequest) ShowResponse
}

type AddRequest struct {
	Person
}

type AddResponse struct{}

type ShowRequest struct {
	// example: Person{Name: "James", Age: 26}
	Person
}

type ShowResponse struct {
	// example: Person{Name: "James"}
	Person
}

type Person struct {
	// Name is the name of the person to show.
	// example: "James"
	Name string
	// Age is the age of a person
	// example: "26"
	Age int
}
