package main

import "net/http"

func (h *Handler) DeleteService(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")

	if serviceID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.store.DeleteService(r.Context(), serviceID)
	if err != nil {
		if err == ErrServiceNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
