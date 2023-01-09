package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	// "os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllExpense(t *testing.T) {
	var e []Expense

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(e), 0)
}

func TestCreateNewExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "test",
		"amount": 50,
		"note": "test promotion discount 20 bath",
		"tags": [
		  "food",
		  "beverage"
		]
	}`)
	var e Expense

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "test", e.Title)
	assert.Equal(t, float64(50), e.Amount)
	assert.Equal(t, "test promotion discount 20 bath", e.Note)
	assert.Equal(t, []string([]string{"food", "beverage"}), e.Tags)
}

func TestGetExpenseByID(t *testing.T) {
	id := 2
	var e Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(id)), nil)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, id, e.ID)
	assert.NotEmpty(t, e.Title)
	assert.NotEmpty(t, e.Amount)
	assert.NotEmpty(t, e.Note)
	assert.NotEmpty(t, e.Tags)
}

func TestUpdateExpenseByID(t *testing.T) {
	id := 1
	body := bytes.NewBufferString(`{
		"title": "testUpdate",
		"amount": 30.1,
		"note": "test promotion discount 40 bath",
		"tags": [
		  "food",
		  "beverage"
		]
	}`)

	var e Expense

	res := request(http.MethodPut, uri("expenses", strconv.Itoa(id)), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "testUpdate", e.Title)
	assert.Equal(t, 30.1, e.Amount)
	assert.Equal(t, "test promotion discount 40 bath", e.Note)
	assert.Equal(t, []string([]string{"food", "beverage"}), e.Tags)
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	// req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
