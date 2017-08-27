package uaa

import (
	"fmt"
	"net/http"
)

func logResponse(response *http.Response, bytes []byte) {
	fmt.Printf("< %v\n", response.Status)
	logHeaders("<", response.Header)
	fmt.Printf("< %v\n", string(bytes[:]))
	fmt.Println()
}

func logRequest(request *http.Request) {
	fmt.Printf("> %v %v\n", request.Method, request.URL.String())
	logHeaders(">", request.Header)
	if request.Body != nil {
		fmt.Printf("> %v\n", request.Body)
	}
	fmt.Println()
}

func logHeaders(prefix string, headers map[string][]string) {
	for header, values := range headers {
		for _, value := range values {
			fmt.Printf("%v %v: %v\n", prefix, header, value)
		}
	}
}
