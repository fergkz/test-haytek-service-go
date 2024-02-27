package InfrastructureController

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type test struct {
	str string
}

func NewTest(
	str string,
) *test {
	controller := new(test)
	controller.str = str
	return controller
}

func (controller *test) Get(w http.ResponseWriter, r *http.Request) {

	log.Println("Teste")

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	dataQuery := mux.Vars(r)["dataQuery"]

	html := []byte(`LOADED` + "\n" + dataQuery)

	io.WriteString(w, string(html))
}
