package main

import (
	"fmt"
	"github.com/victoralmeida428/cdnClientGO/cdn"
	"log"
	"net/http"
)

func init() {
	cndCgf := cdn.GetInstanceConfig()
	cndCgf.Create("http://localhost:8080/cdn")
}

func main() {
	handlePost := func(w http.ResponseWriter, req *http.Request, ) {
		
		Cdn := cdn.New()
		fmt.Println(Cdn)
		ok, err := Cdn.View(1793, w)
		if err != nil {
			log.Fatal(err.Error())
		}
		
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}
		
	}
	http.HandleFunc("/", handlePost)
	
	fmt.Println("Listening on port http://localhost:5000")
	err := http.ListenAndServe("localhost:5000", nil)
	if err != nil {
		panic(err)
	}
}
