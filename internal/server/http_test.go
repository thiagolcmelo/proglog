package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func producerServerAndClient() (*httptest.Server, *http.Client) {
	logServer := newHTTPServer()
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			logServer.handleProduce(rw, r)
		}))
	return server, server.Client()
}

func consumerServerAndClient(records []Record) (*httptest.Server, *http.Client) {
	logServer := newHTTPServer()
	for _, record := range records {
		logServer.Log.Append(record)
	}
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			logServer.handleConsume(rw, r)
		}))
	return server, server.Client()
}

func TestProduce(t *testing.T) {
	server, client := producerServerAndClient()
	defer server.Close()

	produceRequest := ProduceRequest{Record: Record{Value: []byte("iytrfvgs")}}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(produceRequest)
	if err != nil {
		t.Fatal(err)
	}
	body := strings.NewReader(buf.String())

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		server.URL,
		body,
	)

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var result ProduceResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		if err != io.EOF {
			t.Fatal(err)
		}
	}

	expected := ProduceResponse{Offset: 0}
	if diff := cmp.Diff(result, expected); diff != "" {
		t.Fatal(diff)
	}
}

func TestConsume(t *testing.T) {
	records := []Record{
		{Value: []byte("iytrfvgs")},
		{Value: []byte("lihtrefd")},
	}
	server, client := consumerServerAndClient(records)
	defer server.Close()

	consumeRequest := ConsumeRequest{Offset: 1}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(consumeRequest)
	if err != nil {
		t.Fatal(err)
	}
	body := strings.NewReader(buf.String())
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		server.URL,
		body,
	)

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var result ConsumeResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		if err != io.EOF {
			t.Fatal(err)
		}
	}

	expected := ConsumeResponse{Record: Record{Value: []byte("lihtrefd"), Offset: 1}}
	if diff := cmp.Diff(result, expected); diff != "" {
		t.Fatal(diff)
	}
}

func TestProduceInvalidRequestFail(t *testing.T) {
	server, client := producerServerAndClient()
	defer server.Close()

	body := strings.NewReader(`{"something":"off"`)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		server.URL,
		body,
	)

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		body, _ := ioutil.ReadAll(res.Body)
		t.Fatalf("unexpected: %s", body)
	}
}

func TestConsumeInvalidRequestFail(t *testing.T) {
	server, client := consumerServerAndClient(nil)
	defer server.Close()

	body := strings.NewReader(`{"something":"off"`)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		server.URL,
		body,
	)

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		body, _ := ioutil.ReadAll(res.Body)
		t.Fatalf("unexpected: %s", body)
	}
}
