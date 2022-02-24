module github.com/eztrade/kpi

go 1.14

replace github.com/eztrade/login => ../login

require (
github.com/360EntSecGroup-Skylar/excelize/v2 v2.3.2
	github.com/99designs/gqlgen v0.12.2
	github.com/Azure/azure-storage-blob-go v0.10.0
	github.com/eztrade/login v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/gorilla/websocket v1.4.2
	github.com/jackc/pgtype v1.4.2
	github.com/jackc/pgx/v4 v4.8.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.8.0
	github.com/microsoft/ApplicationInsights-Go v0.4.4
	github.com/rs/cors v1.6.0
	github.com/vektah/gqlparser/v2 v2.0.1
	google.golang.org/appengine v1.6.7 // indirect
)
