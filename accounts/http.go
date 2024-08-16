package accounts

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type AccountType uint8

const (
	Checking AccountType = iota
	Money
	Saving
	Other
)

type Account struct {
	ID          uint64      `json:"ID"`
	Balance     float64     `json:"balance"`
	Description string      `json:"description"`
	Type        AccountType `json:"type"`
}

type AccountService struct {
	Repository AccountRepository
}

func logError(w http.ResponseWriter, err error) {
	log.Println("Error: ", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (s AccountService) CreateAccount(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logError(w, err)
		return
	}
	defer r.Body.Close()

	c := Account{}

	if err := json.Unmarshal(body, &c); err != nil {
		logError(w, err)
		return
	}

	if err := s.Repository.Add(c); err != nil {
		logError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type ListAccountsResp struct {
	Accounts []Account `json:"accounts"`
}

func (s AccountService) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.Repository.GetAll()
	if err != nil {
		logError(w, err)
		return
	}

	resp := ListAccountsResp{
		Accounts: accounts,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		logError(w, err)
		return
	}

	w.Write(respJSON)
}
