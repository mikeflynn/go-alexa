package skillserver

import (
	"fmt"
	"net/http"
)

func HomePage(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "Home Page!")
}

func AboutPage(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw, "About Page!")
}
