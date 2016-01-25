package main

import (
	"bytes"
	"fmt"
	"time"
	"testing"
	"net/http"
)

func TestCreateAsset(t *testing.T) {
	url := "http://couch1.vagrant:3000/api/assets"
    fmt.Println("URL:>", url)

    var jsonStr = []byte(`{"name":"Test User", "uri": "myorg://users/test-user"}`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    if resp.Status != "201 Created" {
		t.Errorf("Expected response of 201 Created, but it was %v instead.", resp.Status)
	}

	time.Sleep(250 * time.Millisecond) // Creating the asset can lag on VMs
}

func TestAddNoteToAsset(t *testing.T) {
	url := "http://couch1.vagrant:3000/api/assets"
    fmt.Println("URL:>", url)

    var jsonStr = []byte(`{"note":"This is a note for Test User!", "uri": "myorg://users/test-user"}`)
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    if resp.Status != "200 OK" {
		t.Errorf("Expected response of 200 OK, but it was %v instead.", resp.Status)
	}
}

func TestRemoveAsset(t *testing.T) {
	url := "http://couch1.vagrant:3000/api/assets"
    fmt.Println("URL:>", url)

    var jsonStr = []byte(`{"uri": "myorg://users/test-user"}`)
    req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    if resp.Status != "204 No Content" {
		t.Errorf("Expected response of 200 No Content, but it was %v instead.", resp.Status)
	}
}

func TestCreateInvalidAsset(t *testing.T) {
	url := "http://couch1.vagrant:3000/api/assets"
    fmt.Println("URL:>", url)

    var jsonStr = []byte(`{"name":"Bad Test User", "uri": "urn:invalid:org:test-user"}`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    if resp.Status[:3] != "400" {
		t.Errorf("Expected response of 400, but it was %v instead.", resp.Status[:3])
	}
}

func TestAddNoteToDeletedAsset(t *testing.T) {
	url := "http://couch1.vagrant:3000/api/assets"
    fmt.Println("URL:>", url)

    var jsonStr = []byte(`{"note":"This is a note for Nonexistent User!", "uri": "myorg://users/test-user"}`)
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    if resp.Status[:3] != "404" {
		t.Errorf("Expected response of 400, but it was %v instead.", resp.Status[:3])
	}
}
