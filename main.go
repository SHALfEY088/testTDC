package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

func main() {
	// Подключение к базе данных
	db, err := sql.Open("postgres", "postgres://root:root@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание маршрута для обработки запросов
	http.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {
		// Получение параметров запроса
		clientID := r.FormValue("client_id")
		amount := r.FormValue("amount")

		// Проверка параметров
		if clientID == "" || amount == "" {
			http.Error(w, "Invalid parameters", http.StatusBadRequest)
			return
		}

		// Преобразование параметров в числа
		clientIDInt, err := strconv.Atoi(clientID)
		if err != nil {
			http.Error(w, "Invalid client ID", http.StatusBadRequest)
			return
		}

		amountInt, err := strconv.Atoi(amount)
		if err != nil {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		// Начало транзакции
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Transaction error", http.StatusInternalServerError)
			return
		}

		// Проверка баланса клиента
		var balance int
		err = tx.QueryRow("SELECT balance FROM clients WHERE id = $1", clientIDInt).Scan(&balance)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if balance < amountInt {
			tx.Rollback()
			http.Error(w, "Insufficient funds", http.StatusBadRequest)
			return
		}

		// Обновление баланса клиента
		_, err = tx.Exec("UPDATE clients SET balance = balance - $1 WHERE id = $2", amountInt, clientIDInt)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Запись транзакции в базу данных
		_, err = tx.Exec("INSERT INTO transactions (client_id, amount, result) VALUES ($1, $2, $3)", clientIDInt, amountInt, "success")
		if err != nil {
			tx.Rollback()
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Коммит транзакции
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Transaction error", http.StatusInternalServerError)
			return
		}

		// Отправка ответа
		fmt.Fprintf(w, "Transaction completed successfully")
	})

	// Запуск веб-сервера
	log.Fatal(http.ListenAndServe(":8080", nil))
}
