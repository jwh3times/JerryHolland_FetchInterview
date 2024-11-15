package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	ID           uuid.UUID   `json:"id"`
	Retailer     string      `json:"retailer"`
	PurchaseDate ReceiptDate `json:"purchaseDate"`
	PurchaseTime ReceiptTime `json:"purchaseTime"`
	Total        string      `json:"total"`
	Items        []Item      `json:"items"`
}

type ReceiptID struct {
	ID uuid.UUID `json:"id"`
}

type ReceiptPoints struct {
	Points int `json:"points"`
}

type ReceiptTime time.Time
type ReceiptDate time.Time

var (
	receipts = make(map[uuid.UUID]Receipt)
)

func main() {
	// define the api endpoints and their handler functions
	http.HandleFunc("/receipts/process", receiptsHandler)
	http.HandleFunc("/receipts/{id}/points", pointsHandler)

	// informative log message to indicate server is successfully started or failed to start
	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func receiptsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleProcessReceipt(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func pointsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid receipt ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		handleGetPoints(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Process receipt and store in map in memory for use during the run duration
// of this server. For the purpose of this assignment, no db persistance is used
func handleProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var p = Receipt{}
	var pID = ReceiptID{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &p); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	id := uuid.New()
	p.ID, pID.ID = id, id
	receipts[p.ID] = p
	fmt.Println(receipts)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pID)
}

// Calculate and return the point value computation for the requested
// receipt ID
func handleGetPoints(w http.ResponseWriter, _ *http.Request, id uuid.UUID) {
	var rPoints = ReceiptPoints{}
	p, ok := receipts[id]
	if !ok {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	rPoints.Points = calcReceiptPoints(p)
	if rPoints.Points == -1 {
		fmt.Println("Error calculating receipt points")
		http.Error(w, "Error calculating receipt points", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rPoints)
}

func calcReceiptPoints(p Receipt) int {
	points := 0

	// One point for every alphanumeric character in the retailer name.
	var nonAlphaNumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	points += len(nonAlphaNumericRegex.ReplaceAllString(p.Retailer, ""))

	totalDecimal, err := strconv.ParseInt(p.Total[strings.Index(p.Total, ".")+1:], 10, 32)
	if err != nil {
		return -1
	}

	// 50 points if the total is a round dollar amount with no cents.
	if totalDecimal == 0 {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25.
	if totalDecimal%25 == 0 {
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += ((len(p.Items) / 2) * 5)

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for i := 0; i < len(p.Items); i++ {
		l := len(strings.TrimSpace(p.Items[i].ShortDescription))
		if l%3 == 0 {
			price, err := strconv.ParseFloat(p.Items[i].Price, 64)
			if err != nil {
				return -1
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	// 6 points if the day in the purchase date is odd.
	_, _, day := time.Time(p.PurchaseDate).Date()
	if day%2 == 1 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	hour, min, _ := time.Time(p.PurchaseTime).Clock()
	if hour == 14 || hour == 15 {
		if !(hour == 14 && min == 0) {
			points += 10
		}
	}

	return points
}

func (r *ReceiptDate) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return err
	}
	*r = ReceiptDate(t)
	return nil
}

func (r ReceiptDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(r).Format(time.DateOnly) + `"`), nil
}

func (r *ReceiptTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse(time.TimeOnly, value+":00")
	if err != nil {
		return err
	}
	*r = ReceiptTime(t)
	return nil
}

func (r ReceiptTime) MarshalJSON() ([]byte, error) {
	rtime := time.Time(r).Format(time.TimeOnly)
	return []byte(`"` + rtime[:len(rtime)-3] + `"`), nil
}
