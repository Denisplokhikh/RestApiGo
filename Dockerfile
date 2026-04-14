FROM golang:1.26.2
 
WORKDIR /app
 
# Установим air для hot reload
RUN go install github.com/air-verse/air@latest
 
# Копируем go.mod/go.sum и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download
 
# Копируем весь код
COPY . .
 
EXPOSE 5001
 
# Air будет отслеживать изменения файлов и перекомпилировать
ENTRYPOINT ["air"]
