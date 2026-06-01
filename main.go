package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB= make(map[string]URL)

func generateShortURL(OriginalURL string) string{
	hasher:=md5.New()
	hasher.Write([]byte(OriginalURL)) // converting the orignalURl into small bytes
	// fmt.Println("hasher: ", hasher)
	data:=hasher.Sum(nil)
	hash:=hex.EncodeToString(data)

	return hash[:8]

}

func createURL(originalURL string){
	shortURL:= generateShortURL(originalURL)
	id:=shortURL
	newURL:= URL{
		ID:id,
		OriginalURL:originalURL,
		ShortURL:shortURL,
		CreationDate:time.Now(),
	}

	urlDB[id]=newURL
}

func getURL(id string) (URL,error){
	url,exists:=urlDB[id]
	if !exists{
		return URL{},fmt.Errorf("URL not found")
	}
	return url,nil
	
}

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello, World!")

}


func main(){
	fmt.Println("starting URL Shortener Service...")
	originalURL:= "https://www.google.com"
		shortURL:= generateShortURL(originalURL)
		fmt.Println("shortURL: ", shortURL)
	
	
		//Register the handler function to handle the requests to the root url ("/")
		http.HandleFunc("/",handler)
	fmt.Println("starting the server on port 3000...")	
	err:=http.ListenAndServe(":3000",nil)
	if err!=nil{
		fmt.Println("Error starting server: ", err)
		return
	}



	




}