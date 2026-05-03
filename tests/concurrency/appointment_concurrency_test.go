package concurrency

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	baseURL = "http://localhost:8080/api/v1/app/appointment"
)

// Adjust this payload to match your seeded test data.
type createAppointmentRequest struct {
	DealershipID int    `json:"dealership_id"`
	ServiceID    int    `json:"service_id"`
	BayTypeID    int    `json:"bay_type_id"`
	CustomerName string `json:"customer_name"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
}

func TestConcurrentBooking_OnlyOneShouldSucceed(t *testing.T) {
	t.Parallel()

	const concurrentClients = 10

	payload := createAppointmentRequest{
		DealershipID: 1,
		ServiceID:    2,
		BayTypeID:    2,
		CustomerName: "race-test",
		StartTime:    "2026-05-02T08:00:00+07:00",
		EndTime:      "2026-05-02T11:00:00+07:00",
	}

	var (
		successCount int32
		failCount    int32
		wg           sync.WaitGroup
	)

	start := make(chan struct{})

	for i := 0; i < concurrentClients; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			<-start

			ok, err := sendCreateAppointment(i, payload)
			if err != nil {
				t.Logf("[client-%d] request error: %v", i, err)
				atomic.AddInt32(&failCount, 1)
				return
			}

			if ok {
				atomic.AddInt32(&successCount, 1)
				t.Logf("[client-%d] booking success", i)
			} else {
				atomic.AddInt32(&failCount, 1)
				t.Logf("[client-%d] booking rejected (expected for losers)", i)
			}
		}(i)
	}

	close(start)
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("expected exactly 2 successful booking, got %d", successCount)
	}

	if failCount != concurrentClients-1 {
		t.Fatalf("expected %d failed bookings, got %d", concurrentClients-1, failCount)
	}
}

func sendCreateAppointment(i int, payload createAppointmentRequest) (bool, error) {
	reqBody := payload
	reqBody.CustomerName = fmt.Sprintf("%s-%d", payload.CustomerName, i)

	body, err := json.Marshal(reqBody)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return true, nil

	case http.StatusConflict, http.StatusBadRequest, http.StatusUnprocessableEntity:
		return false, nil

	default:
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
