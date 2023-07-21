migrateup:
	migrate -path internal/db/postgres/migration -database "postgresql://postgres:Qwerty123!@localhost:5432/ShortURLDataDB?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/db/postgres/migration -database "postgresql://postgres:Qwerty123!@localhost:5432/ShortURLDataDB?sslmode=disable" -verbose down

migrateup-test:
	migrate -path internal/db/postgres/migration -database "postgresql://postgres:Qwerty123!@localhost:5432/TestShortURLDataDB?sslmode=disable" -verbose up

migratedown-test:
	migrate -path internal/db/postgres/migration -database "postgresql://postgres:Qwerty123!@localhost:5432/TestShortURLDataDB?sslmode=disable" -verbose down


.PHONY: migrateup migratedown migrateup-test migratedown-test
