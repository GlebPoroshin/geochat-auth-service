module github.com/GlebPoroshin/geochat-auth-service

go 1.24.0

require (
	github.com/GlebPoroshin/geochat-shared v0.0.0
	github.com/golang-jwt/jwt/v4 v4.5.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.34.0
	gorm.io/gorm v1.25.12
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.22.0 // indirect
)

replace github.com/GlebPoroshin/geochat-shared => ../geochat-shared
