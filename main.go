package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	InfrastructureController "github.com/fergkz/test-haytek-service-go/src/Infrastructure/Controller"
	InfrastructureService "github.com/fergkz/test-haytek-service-go/src/Infrastructure/Service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("PANIC/ERROR interno:", err)
			errStr := fmt.Sprintf("%s", err)
			os.WriteFile("ERROR.txt", []byte(errStr), 0644)
		}
	}()

	config := new(Config)
	config.Load("config.yml")

	serviceAddress := InfrastructureService.NewHaytekAddress()
	serviceBox := InfrastructureService.NewHaytekBox()
	serviceCarrier := InfrastructureService.NewHaytekCarrier()
	serviceOrder := InfrastructureService.NewHaytekOrder()

	router := mux.NewRouter()
	router.StrictSlash(true)

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
	})

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"*"},
	})
	handler := c.Handler(router)

	apiRouter := router.PathPrefix("/v1").Subrouter()
	apiRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Expose-Headers", "*")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	})
	apiRouter.StrictSlash(true)

	// Rota simplificada para verificar saúde do serviço.
	router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, string("ALIVE"))
	})

	controllerDeliveryPack := InfrastructureController.NewDeliveryPack(
		serviceAddress,
		serviceBox,
		serviceCarrier,
		serviceOrder,
	)

	apiRouter.HandleFunc("/delivery-pack", controllerDeliveryPack.Get).Methods("GET")

	log.Printf("Server started at port %s\n", config.Server.Port)

	http.ListenAndServe("127.0.0.1:"+config.Server.Port, handler)
	server := &http.Server{
		Addr:         "127.0.0.1:" + config.Server.Port, // Porta em que o servidor irá ouvir.
		Handler:      handler,                           // Roteador Mux.
		ReadTimeout:  120 * time.Second,                 // Tempo limite de leitura.
		WriteTimeout: 120 * time.Second,                 // Tempo limite de escrita.
	}
	log.Fatal(server.ListenAndServe())
}
