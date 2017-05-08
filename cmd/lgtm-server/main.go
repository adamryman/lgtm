package main

import (
	"fmt"
	. "github.com/y0ssar1an/q"
	"net/http"

	"github.com/StudentRND/lgtm"
)

func main() {
	fmt.Println("You can do anything!")
	Q("Lets debug some shit")
	fmt.Println(lgtm.Start())
	http.ListenAndServe(":8080", nil)
}
