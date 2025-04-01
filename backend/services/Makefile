run:
	go run auth-service/main.go &
	go run project-service/main.go &
	go run task-service/main.go &
	go run api-gateway/main.go &

stop:
	@echo "Deteniendo microservicios..."
	@pkill -f "project-service"
	@pkill -f "task-service"
	@pkill -f "notifications-service"
	@pkill -f "api-gateway"
	@echo "Deteniendo Kafka y Zookeeper..."
	@cd /d C:\kafka
	@bin\windows\kafka-server-stop.bat
	@bin\windows\zookeeper-server-stop.bat
	@echo "Todos los servicios han sido detenidos."
