package main

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	//mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)
	ctx, cancelCtx := context.WithCancel(context.Background())
	server1 := &http.Server{Addr: ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {

			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			ctx = context.WithValue(ctx, "uuid", uuid.NewV4().String())
			return ctx
		}}
	server2 := &http.Server{
		Addr:    ":4444",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	go func() {
		err := server1.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Server closed")
		} else if err != nil {
			log.Fatal("Error on server start!")
			//os.Exit(1)
		}
		cancelCtx()
	}()
	go func() {
		err := server2.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Server closed")
		} else if err != nil {
			log.Fatal("Error on server start!")
			//os.Exit(1)
		}
		cancelCtx()
	}()
	<-ctx.Done()
}

const keyServerAddr = "serverAddr"

func getHello(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	fmt.Printf(ctx.Value("uuid").(string))
	fmt.Fprintf(writer, "Hello response")

}

func getRoot(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))
	fmt.Println("Got root request")
}
