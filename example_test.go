package securerequest_test

import (
	"fmt"
	"net/http"

	"github.com/pcasaretto/securerequest"
)

func ExampleValidatingHandler() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Super secret thing")
	})

	secrets := map[string]string{"app": "secret"}
	http.Handle("/foo", securerequest.NewValidatingHandler(handler, secrets))
	http.ListenAndServe(":8080", nil)
}
