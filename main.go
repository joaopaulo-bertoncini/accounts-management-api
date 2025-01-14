package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type AccountType string

const (
	Payable    AccountType = "payable"
	Receivable AccountType = "receivable"
)

type Account struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Amount      float64     `json:"amount"`
	AccountType AccountType `json:"account_type"`
}

type Server struct {
	db *sql.DB
}

func (s *Server) CreateAccount(c *gin.Context) {
	var account Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := "INSERT INTO accounts (name, amount, account_type) VALUES (?, ?, ?)"
	result, err := s.db.Exec(query, account.Name, account.Amount, account.AccountType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	id, _ := result.LastInsertId()
	account.ID = int(id)
	c.JSON(http.StatusCreated, account)
}

func (s *Server) ListAccounts(c *gin.Context) {
	rows, err := s.db.Query("SELECT id, name, amount, account_type FROM accounts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accounts"})
		return
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.ID, &account.Name, &account.Amount, &account.AccountType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse accounts"})
			return
		}
		accounts = append(accounts, account)
	}

	c.JSON(http.StatusOK, accounts)
}

func (s *Server) UpdateAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var updatedAccount Account
	if err := c.ShouldBindJSON(&updatedAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := "UPDATE accounts SET name = ?, amount = ?, account_type = ? WHERE id = ?"
	_, err = s.db.Exec(query, updatedAccount.Name, updatedAccount.Amount, updatedAccount.AccountType, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	updatedAccount.ID = id
	c.JSON(http.StatusOK, updatedAccount)
}

func (s *Server) DeleteAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	query := "DELETE FROM accounts WHERE id = ?"
	_, err = s.db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted"})
}

func (s *Server) InitializeDatabase() {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		amount REAL NOT NULL,
		account_type TEXT NOT NULL
	)`

	if _, err := s.db.Exec(query); err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("sqlite3", "accounts.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := &Server{db: db}
	s.InitializeDatabase()

	r := gin.Default()
	r.POST("/accounts", s.CreateAccount)
	r.GET("/accounts", s.ListAccounts)
	r.PUT("/accounts/:id", s.UpdateAccount)
	r.DELETE("/accounts/:id", s.DeleteAccount)

	r.Run(":8080")
}
