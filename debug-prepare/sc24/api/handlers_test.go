package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sc24/store"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateWorkerHandler(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	realStore := store.NewTestStore(db)
	server := &Server{Store: realStore}

	t.Run("Successful Creation", func(t *testing.T) {
		mock.ExpectExec("^INSERT INTO workers").
			WithArgs("worker1", "echo Hello World", 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		req_body := `{
			"name": "worker1",
			"script_content": "echo Hello World",
			"limit": 2
		}`

		req, err := http.NewRequest("PUT", "/workers/", bytes.NewBufferString(req_body))
		rr := httptest.NewRecorder()

		if err != nil {
			t.Fatal(err)
		}
		handler := http.HandlerFunc(server.CreateWorkerHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		if rr.Body.String() != "Created" {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), "Created")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("Missing Body", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/workers/", bytes.NewBufferString(""))
		rr := httptest.NewRecorder()
		if err != nil {
			t.Fatal(err)
		}
		handler := http.HandlerFunc(server.CreateWorkerHandler)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
}
