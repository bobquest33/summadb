package handle

import (
	"net/http"

	db "github.com/fiatjaf/summadb/database"
)

func Destroy(w http.ResponseWriter, r *http.Request) {
	db.End()
	db.Erase()
	db.Start()
	w.WriteHeader(200)
}
