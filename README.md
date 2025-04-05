# mvp_projects

# backend (solo go con fiber)

PROJECT-SERVICE

go get github.com/gofiber/fiber/v2
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/go-playground/validator/v10
go get github.com/gofiber/jwt/v3
go get github.com/confluentinc/confluent-kafka-go/kaf

PA TODOS LIT:
github.com/gofiber/fiber/v2
github.com/google/uuid
github.com/jackc/pgx/v5
github.com/joho/godotenv
go get github.com/jackc/pgx/v5/pgxpool

GORM

pal jwt:
"github.com/golang-jwt/jwt/v5"

comandos:
go mod init nombre_service

go get

go mod tidy

go run main.go

air

Reiniciar:
rm -rf go.mod go.sum vendor
go mod init ---

//INICIAR ZOOKEPER
bin\windows\zookeeper-server-start.bat config\zookeeper.properties

//INICIAR KAFKA
bin\windows\kafka-server-start.bat config\server.properties

//AHORA KAFKA

bin\windows\kafka-topics.bat --create --topic notifications --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1

EN GO para Kafka:
go get github.com/segmentio/kafka-go

DETENER KAFKA BROKER
bin\windows\kafka-server-stop.bat

DETENER ZOOKEPER
bin\windows\zookeeper-server-stop.bat

# backend_pruebas_go_fiber + ent

Datos:
pgx + ent

go install entgo.io/ent/cmd/ent@latest

go mod init go_fiber

Al trabajar con ent:
se define tu task.go con el comando:
ent new Task

y recordar siempre usar el go mod tidy

y luego usar el comando generate para generar la entidad:
ent generate ./ent/schema
USAR:
go generate ./ent

para q funcione bien migrate de ent
go get ariga.io/atlas/sql/migrate

RECORDAR para las importaciones acordarse de q nombre se pone en go mod init: "go_fiber/db"

Dependecias de driver entre pgx y ent, driver:
go get entgo.io/ent
go get github.com/jackc/pgx/v5

Como la mayoria de ORM, la convencion q utilizan es q cada tabla se llame en plural y minusculas:

tasks -> Correcto
Task -> NO
task -> NO
