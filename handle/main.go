package handle

import (
	"net/http"

	"github.com/carbocation/interpose"
	"github.com/carbocation/interpose/adaptors"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func BuildHTTPMux() *interpose.Middleware {
	//log.WithFields(log.Fields{
	//	"DBFILE":       settings.DBFILE,
	//	"PORT":         settings.PORT,
	//	"CORS_ORIGINS": settings.CORS_ORIGINS,
	//	"STARTTIME":    settings.STARTTIME,
	//}).Info("starting database server.")

	master := interpose.New()
	router := mux.NewRouter()
	master.UseHandler(router)

	// special database actions -- different context, no CORS
	// matches if request comes with a special header
	actionsMiddle := interpose.New()
	actions := mux.NewRouter()
	actions.HandleFunc("/_destroy", Destroy).Methods("POST")
	actionsMiddle.UseHandler(actions)
	router.Headers("Summadb-Admin", "true").Handler(actionsMiddle)

	// normal requests -- matches everything
	normalMiddle := interpose.New()
	normal := mux.NewRouter()
	normalMiddle.Use(setCommonVariables)
	normalMiddle.Use(adaptors.FromNegroni(cors.New(cors.Options{
		// CORS
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "If-Match"},
		AllowCredentials: true,
	})))
	normal.HandleFunc("/{path:.*}", Get).Methods("GET")
	normal.HandleFunc("/{path:.*}", Put).Methods("PUT")
	normal.HandleFunc("/{path:.*}", Patch).Methods("PATCH")
	normal.HandleFunc("/{path:.*}", Delete).Methods("DELETE")
	normal.HandleFunc("/{path:.*}", Post).Methods("POST")
	normalMiddle.UseHandler(normal)
	router.MatcherFunc(func(_ *http.Request, _ *mux.RouteMatch) bool { return true }).Handler(normalMiddle)

	return master
}
