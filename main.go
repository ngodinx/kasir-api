package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/repositories"
	"kasir-api/services"
	"kasir-api/handlers"
	"net/http"
	"strings"
	"os"
	"log"
	"time/tzdata"
	"github.com/spf13/viper"
	_ "kasir-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/report/hari-ini", transactionHandler.HandleReportHariIni)
	http.HandleFunc("/api/report", transactionHandler.HandleReport)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	http.HandleFunc("/api/kategori", categoryHandler.HandleCategories)
	http.HandleFunc("/api/kategori/", categoryHandler.HandleCategoryByID)

	// 1) Root redirect ke swagger-ui
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	// 2) Serve UI custom (neobrutal)
	http.Handle("/swagger-ui/",
		http.StripPrefix("/swagger-ui/",
			http.FileServer(http.Dir("./swagger-ui")),
		),
	)

	// 3) Serve swagger spec (doc.json) via swaggo
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/health", HealthCheck)

	fmt.Println("Server running di localhost:"+config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "API Running",
	})
}
