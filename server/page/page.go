package page

import (
	"fmt"
	"net/http"
)

// Page holds data for displaying webpages
type Page struct {
	Body []byte
}

// Serves the web-page contained in the Page struct.
func (p *Page) ServePage(w http.ResponseWriter) {
	fmt.Fprintf(w, "%s", p.Body)
}

func servePage(w http.ResponseWriter, body []byte) {
	fmt.Fprintf(w, "%s", body)
}
