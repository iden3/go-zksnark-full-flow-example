package zkflowexample

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func printT(s string) {
	fmt.Printf("%s "+s+"\n", time.Now().Format("15:04:05"))
}

func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func downloadFiles(path, filesServer string) error {
	// download wasm witness calculation, proving key and verification key.
	printT("- Download witness calculation")
	if err := downloadFile(
		path+"/circuit.wasm",
		filesServer+"/circuit.wasm",
	); err != nil {
		return err
	}
	printT("- Downloading proving_key.json...")
	if err := downloadFile(
		path+"/proving_key.json",
		filesServer+"/proving_key.json",
	); err != nil {
		return err
	}
	printT("- Download verification_key.json")
	if err := downloadFile(
		path+"/verification_key.json",
		filesServer+"/verification_key.json",
	); err != nil {
		return err
	}
	printT("- Download input.json")
	if err := downloadFile(
		path+"/input.json",
		filesServer+"/input.json",
	); err != nil {
		return err
	}
	fmt.Print("=============\n\n\n")

	return nil
}
