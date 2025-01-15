package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	query := `CREATE TABLE accounts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        amount REAL NOT NULL,
        account_type TEXT NOT NULL
    )`
	if _, err := db.Exec(query); err != nil {
		panic(err)
	}

	return db
}

func TestCreateAccount(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	server := &Server{db: db}
	router := gin.Default()
	router.POST("/accounts", server.CreateAccount)

	body := `{
        "name": "Internet Bill",
        "amount": 100.50,
        "account_type": "payable"
    }`

	req := httptest.NewRequest("POST", "/accounts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestListAccounts(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	server := &Server{db: db}
	router := gin.Default()
	router.GET("/accounts", server.ListAccounts)

	db.Exec("INSERT INTO accounts (name, amount, account_type) VALUES (?, ?, ?)", "Electricity Bill", 150.75, "payable")

	req := httptest.NewRequest("GET", "/accounts", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var accounts []Account
	if err := json.Unmarshal(w.Body.Bytes(), &accounts); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(accounts) != 1 || accounts[0].Name != "Electricity Bill" {
		t.Errorf("Unexpected accounts data: %+v", accounts)
	}
}

func TestUpdateAccount(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	server := &Server{db: db}
	router := gin.Default()
	router.PUT("/accounts/:id", server.UpdateAccount)

	db.Exec("INSERT INTO accounts (name, amount, account_type) VALUES (?, ?, ?)", "Water Bill", 80.00, "payable")

	body := `{
        "name": "Updated Water Bill",
        "amount": 90.00,
        "account_type": "receivable"
    }`

	req := httptest.NewRequest("PUT", "/accounts/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response Account
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Name != "Updated Water Bill" || response.Amount != 90.00 || response.AccountType != "receivable" {
		t.Errorf("Unexpected account data: %+v", response)
	}
}

func TestDeleteAccount(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	server := &Server{db: db}
	router := gin.Default()
	router.DELETE("/accounts/:id", server.DeleteAccount)

	db.Exec("INSERT INTO accounts (name, amount, account_type) VALUES (?, ?, ?)", "Phone Bill", 50.00, "payable")

	req := httptest.NewRequest("DELETE", "/accounts/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["message"] != "Account deleted" {
		t.Errorf("Unexpected response: %v", response)
	}
}
