package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"price-fetcher/types"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error
type JSONAPIServer struct {
	listenAddr string
	svc        PriceFetcher
}

func NewJSONAPIServer(listenAdd string, svc PriceFetcher) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAdd,
		svc:        svc,
	}
}
func (s *JSONAPIServer) Run() {
	http.HandleFunc("/", makeHTTPHandlerAPIFunc(s.handFetchPrice))
	http.ListenAndServe(s.listenAddr, nil)
}
func makeHTTPHandlerAPIFunc(apiFn APIFunc) http.HandlerFunc {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestID", rand.Int())
	return func(w http.ResponseWriter, r *http.Request) {
		if err := apiFn(context.Background(), w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
	}
}
func (s *JSONAPIServer) handFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")
	price, err := s.svc.FetchPrice(ctx, ticker)
	if err != nil {
		return err
	}
	priceRes := types.PriceResponse{
		Price:  price,
		Ticker: ticker,
	}
	return writeJSON(w, http.StatusOK, &priceRes)
}
func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
