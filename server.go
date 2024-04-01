package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Olá, mundo!") // Escreve "Olá, mundo!" na resposta HTTP
}

func server() {
	http.HandleFunc("/", handler) // Define a função handler como a handler para o caminho "/"
	fmt.Println("Servidor HTTP iniciado. Acesse http://localhost:8080/")
	http.ListenAndServe(":8080", nil) // Inicia o servidor na porta 8080
}
