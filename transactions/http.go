package transactions

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type TransactionType uint8

const (
	Income TransactionType = iota
	Expense
)

type Period uint8

const (
	Daily Period = iota
	Weekly
	Monthly
	Yearly
)

type Repeat struct {
	Times  uint8  `json:"times"`
	Period Period `json:"period"`
}

const dateLayout = "2006-01-02"

type CustomDate struct {
	time.Time
}

func (c CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Format(dateLayout))
}

func (t *CustomDate) UnmarshalJSON(b []byte) (err error) {
	str := string(b)
	str = str[1 : len(str)-1]
	date, err := time.Parse(dateLayout, str)
	if err != nil {
		return err
	}
	t.Time = date
	return
}

type Transaction struct {
	ID          uint64          `json:"ID"`
	Value       float64         `json:"value"`
	Paid        bool            `json:"paid"`
	Date        CustomDate      `json:"date"`
	Description string          `json:"description"`
	CategoryID  uint64          `json:"category_id"`
	AccountID   uint64          `json:"account_id"`
	TagID       uint64          `json:"tag_id"`
	Type        TransactionType `json:"type"`
	Fixed       bool            `json:"fixed"`
	Repeat      Repeat          `json:"repeat"`
}

type TransactionService struct {
	Repo TransactionRepository
}

func logError(w http.ResponseWriter, err error) {
	log.Println("Error: ", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (s TransactionService) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logError(w, err)
		return
	}
	defer r.Body.Close()

	t := Transaction{}

	if err := json.Unmarshal(body, &t); err != nil {
		logError(w, err)
		return
	}

	if t.Fixed {
		if err := s.Repo.AddFixed(t); err != nil {
			logError(w, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}

	if err := s.Repo.Add(t); err != nil {
		logError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type ListTransactionsResp struct {
	Transactions []Transaction `json:"transactions"`
}

func (s TransactionService) ListTransactions(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	fromStr := values.Get("from")
	toStr := values.Get("to")

	from, err := time.Parse(dateLayout, fromStr)
	if err != nil {
		logError(w, err)
		return
	}

	to, err := time.Parse(dateLayout, toStr)
	if err != nil {
		logError(w, err)
		return
	}

	transactions, err := s.Repo.List(from, to)
	if err != nil {
		log.Println("Failed to fetch tags: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := ListTransactionsResp{
		Transactions: transactions,
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error trying to marshal JSON: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respJSON)
}
