package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandlerList_Order(t *testing.T) {
	h := HandlerList{
		http.HandlerFunc(handlerA),
		http.HandlerFunc(handlerB),
		http.HandlerFunc(handlerC),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, nil)
	if !reflect.DeepEqual(rr.HeaderMap["Handler"], []string{"a", "b", "c"}) {
		t.Fatalf("Header list is incompete or out of order: %#v", rr.HeaderMap["Handler"])
	}
}

func TestHandlerList_Fail(t *testing.T) {
	h := HandlerList{
		http.HandlerFunc(handlerA),
		http.HandlerFunc(handlerB),
		http.HandlerFunc(failHandler),
		http.HandlerFunc(handlerC),
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, nil)
	if !reflect.DeepEqual(rr.HeaderMap["Handler"], []string{"a", "b"}) {
		t.Fatalf("Header list is incompete or out of order: %#v", rr.HeaderMap["Handler"])
	}
}

func handlerA(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Handler", "a")
}

func handlerB(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Handler", "b")
}

func handlerC(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Handler", "c")
}

func failHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Error", http.StatusInternalServerError)
}
