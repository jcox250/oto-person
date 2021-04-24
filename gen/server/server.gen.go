// Code generated by oto; DO NOT EDIT.

package server

import (
	"context"
	"net/http"

	"github.com/pacedotdev/oto/otohttp"
	
)


type PersonService interface {

	// Add adds a person
Add(context.Context, AddRequest) (*AddResponse, error)
	// Show shows a person
Show(context.Context, ShowRequest) (*ShowResponse, error)
}



type personServiceServer struct {
	server *otohttp.Server
	personService PersonService
}

// Register adds the PersonService to the otohttp.Server.
func RegisterPersonService(server *otohttp.Server, personService PersonService) {
	handler := &personServiceServer{
		server: server,
		personService: personService,
	}
	server.Register("PersonService", "Add", handler.handleAdd)
	server.Register("PersonService", "Show", handler.handleShow)
	}

func (s *personServiceServer) handleAdd(w http.ResponseWriter, r *http.Request) {
	var request AddRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.personService.Add(r.Context(), request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}

func (s *personServiceServer) handleShow(w http.ResponseWriter, r *http.Request) {
	var request ShowRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.personService.Show(r.Context(), request)
	if err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
}




type Person struct {
	// Name is the name of the person to show.
Name string `json:"name"`
// Age is the age of a person
Age int `json:"age"`

}

type AddRequest struct {
	Person Person `json:"person"`

}

type AddResponse struct {
	// Error is string explaining what went wrong. Empty if everything was fine.
Error string `json:"error,omitempty"`

}

type ShowRequest struct {
	Person Person `json:"person"`

}

type ShowResponse struct {
	Person Person `json:"person"`
// Error is string explaining what went wrong. Empty if everything was fine.
Error string `json:"error,omitempty"`

}
