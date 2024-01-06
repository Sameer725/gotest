package main

import (
	"net/http"
	"testing"

	"snippetbox.samkha.com/internal/assert"
)

func TestPing(t *testing.T) {
	// unit test
	// rr := httptest.NewRecorder()

	// r, err := http.NewRequest(http.MethodGet, "/", nil)

	// if err != nil {
	// 	t.Fatal(err)
	// }

	// ping(rr, r)

	// rs := rr.Result()

	// assert.Equal[int](t, rs.StatusCode, http.StatusOK)

	// defer rs.Body.Close()

	// body, err := io.ReadAll(rs.Body)

	// if err != nil {
	// 	t.Fatal(err)
	// }
	// body = bytes.TrimSpace(body)

	// assert.Equal[string](t, string(body), "OK")

	//e2e
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())

	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
