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

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid Id",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old Silent Pond",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		}, {
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		}, {
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
