package main

import (
	"net/http"
)

// handlerReset handles the reset endpoint and responds with an HTML page indicating that the Chirpy server has been reset.
// Also reset the user table and cascading all the others.
// It can be reseted onlz in development mode
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment"))
		return
	}

	err := cfg.db.Reset(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't clear tables", err)
		return
	}

	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been reseted to 0</p>
</body>

</html>`))
}

