package main

import (
	"io"
	"log"
	"net/http"
	"tool/getDecryptedPaper"
	"tool/savePaper"
)

func main() {
	// 目标根URL
	url := "http://121.43.151.190:8000/"
	// 发送 GET 请求,返回的结果还需要进行处理才能得到你需要的结果
	paperres, err := http.Get(url + "paper")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer paperres.Body.Close()

	paperbody, err := io.ReadAll(paperres.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	secretres, err := http.Get(url + "secret")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer secretres.Body.Close()

	secretbody, err := io.ReadAll(secretres.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	secret := string(secretbody)
	decryptedpaper, err := getDecryptedPaper.GetDecryptedPaper(string(paperbody), string(secret))
	if err != nil {
		log.Fatalf("Failed to decode paper: %v", err)
	}

	savepath := "C:\\Users\\JRH\\Desktop\\muxi-backend\\paper\\Academician Sun's papers.txt"
	err = savePaper.SavePaper(savepath, decryptedpaper)
	if err != nil {
		log.Fatalf("Failed to save decrypted paper: %v", err)
	}

}
