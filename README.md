# Migrate up
migrate -database 'postgres://kawaii:123456@0.0.0.0:4444/kawaii_db_test?sslmode=disable' -source file://D:/path-to-migrate -verbose up

# Migrate down
migrate -database 'postgres://kawaii:123456@0.0.0.0:4444/kawaii_db_test?sslmode=disable' -source file://D:/path-to-migrate -verbose down