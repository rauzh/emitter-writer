package networking

import (
	"gitlab.com/Skinass/hakaton-2023-1-1/writer/processing"
	"io"
	"log"
	"net/http"
	"os"
)

type Handler struct {
	File *os.File
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	respBody, err := io.ReadAll(r.Body)

	if err == nil {
		w.WriteHeader(http.StatusOK)
		// engine
	} else {
		log.Printf("err body read")
	}

	mes, err := processing.MessageInitTime(respBody)

	if err != nil {
		log.Printf("err init time: %s", err)
	}

	err = processing.MessageWrite(mes, h.File)

	if err != nil {
		log.Printf("message write: %s", err)
	}

	if r.Body.Close() != nil {
		log.Printf("err body close")
	}
}
