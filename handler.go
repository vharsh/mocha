package mocha

import (
	"bufio"
	"log"
	"net/http"

	"github.com/vitorsalgado/mocha/matcher"
)

type Handler struct {
	mocks   MockStore
	parsers []BodyParser
	extras  Extras
}

func newHandler(
	mockstore MockStore,
	parsers []BodyParser,
	extras Extras,
) *Handler {
	return &Handler{mocks: mockstore, parsers: parsers, extras: extras}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parsedbody, err := ParseRequestBody(r, h.parsers)
	if err != nil {
		respondErr(w, err)
		return
	}

	params := matcher.Params{RequestInfo: &matcher.RequestInfo{Request: r, ParsedBody: parsedbody}, Extras: &h.extras}
	result, err := findMockForRequest(h.mocks, params)

	if err != nil {
		respondErr(w, err)
		return
	}

	if !result.Matches {
		noMatch(w, result)
		return
	}

	mock := result.Matched
	res, err := result.Matched.Responder(r, mock)
	if err != nil {
		respondErr(w, err)
		return
	}

	mock.Hits++
	h.mocks.Save(mock)

	for k, v := range res.Headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(res.Status)

	if res.Body == nil {
		return
	}

	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		_, err = w.Write(scanner.Bytes())
		if err != nil {
			respondErr(w, err)
		}
	}

	if scanner.Err() != nil {
		respondErr(w, err)
		return
	}
}

func noMatch(w http.ResponseWriter, result *findMockResult) {
	w.WriteHeader(http.StatusTeapot)
	_, _ = w.Write([]byte("Request was not matched."))

	if result.ClosestMatch != nil {
		_, _ = w.Write([]byte("\n"))
		_, _ = w.Write([]byte("\n"))
	}
}

func respondErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusTeapot)
	_, _ = w.Write([]byte("Request was not matched."))
	_, _ = w.Write([]byte(err.Error()))

	log.Printf("Reason: %v", err)
}
