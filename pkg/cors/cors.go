package cors

import "github.com/rs/cors"

func New() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://83.166.253.130:3000",

			"http://localhost:8003",
			"http://83.166.253.130:8003",
			"http://127.0.0.1:8003",

			"http://localhost:8088",
			"http://83.166.253.130:8088",
			"http://127.0.0.1:8088",

			"http://localhost:8089",
			"http://83.166.253.130:8089",
			"http://127.0.0.1:8089",

			"http://212.233.75.64",
			"http://212.233.75.64:80",
			"http://212.233.75.64:3000",
		},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
			"Accept",
			"Origin",
			"X-Requested-With",
		},
		ExposedHeaders: []string{"Content-Length"},
		MaxAge:         300,
	})
}
