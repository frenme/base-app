package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Order service") // ✅ Отправляем строку в ответ
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(":8082", nil)
}
