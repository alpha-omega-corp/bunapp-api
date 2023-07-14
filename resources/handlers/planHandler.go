package handlers

import (
	"chadgpt-api/app"
	"github.com/uptrace/bunrouter"
	"net/http"
)

type PlanHandler struct {
	app *app.App
}

type CreatePlanRequest struct {
	Diet string `json:"diet"`
}

func NewPlanHandler(app *app.App) *PlanHandler {
	return &PlanHandler{
		app: app,
	}
}

func (h *PlanHandler) Create(w http.ResponseWriter, req bunrouter.Request) error {
	client := h.app.GptClient()
	res, err := client.UserReq(`
		"Summarize the text below as a bullet point list of the most important points."

		Text: """
		 	Go is a new language. Although it borrows ideas from existing languages, it has unusual properties that make effective Go programs different in character from programs written in its relatives. A straightforward translation of a C++ or Java program into Go is unlikely to produce a satisfactory result—Java programs are written in Java, not Go. On the other hand, thinking about the problem from a Go perspective could produce a successful but quite different program. In other words, to write Go well, it's important to understand its properties and idioms. It's also important to know the established conventions for programming in Go, such as naming, formatting, program construction, and so on, so that programs you write will be easy for other Go programmers to understand.
			This document gives tips for writing clear, idiomatic Go code. It augments the language specification, the Tour of Go, and How to Write Go Code, all of which you should read first.
			Note added January, 2022: This document was written for Go's release in 2009, and has not been updated significantly since. Although it is a good guide to understand how to use the language itself, thanks to the stability of the language, it says little about the libraries and nothing about significant changes to the Go ecosystem since it was written, such as the build system, testing, modules, and polymorphism. There are no plans to update it, as so much has happened and a large and growing set of documents, blogs, and books do a fine job of describing modern Go usage. Effective Go continues to be useful, but the reader should understand it is far from a complete guide. See issue 28782 for context. 
		"""
	`)
	if err != nil {
		return err
	}

	return bunrouter.JSON(w, res)
}
