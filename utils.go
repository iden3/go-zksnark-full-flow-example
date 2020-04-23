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
