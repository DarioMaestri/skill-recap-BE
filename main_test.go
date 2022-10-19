package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DarioMaestri/01_skill-recap/skilldomain"
)

func TestSkillAlive(t *testing.T) {

	req, err := http.NewRequest("GET", "/skills/alive", nil)
	if err != nil {
		t.Error(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(skilldomain.Alive)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
