package database

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const fn = "cockroach-public-root.crt"

func loadCockroachRootCert(ctx context.Context) (string, error) {
	_, err := os.Stat(fn)
	if err == nil {
		return fn, nil
	}

	// assume we need to get the file
	log.Println("loading cluster root certificate")

	clusterID := os.Getenv("DB_CLUSTER_ID")
	if clusterID == "" {
		return "", fmt.Errorf("DB_CLUSTER_ID is not set")
	}

	url := "https://cockroachlabs.cloud/clusters/" + url.PathEscape(clusterID) + "/cert"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fn)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return fn, err
}

func removeCert(ctx context.Context) error {
	_, err := os.Stat(fn)
	if err != nil {
		return nil
	}

	log.Println("removing cluster root certificate")
	return os.Remove(fn)
}
