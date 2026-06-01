package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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

func createURL(originalURL string) string{
	shortURL:= generateShortURL(originalURL)
	id:=shortURL
	newURL:= URL{
		ID:id,
		OriginalURL:originalURL,
		ShortURL:shortURL,
		CreationDate:time.Now(),
	}

	urlDB[id]=newURL
	return shortURL
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

func shortenURLHandler(w http.ResponseWriter, r *http.Request){
	var data struct{
		URL string `json:"url"`
	}

	err:=json.NewDecoder(r.Body).Decode(&data)
	if err!=nil{
		http.Error(w,"Invalid request body", http.StatusBadRequest)
		return
	}


	shortURL_:= createURL(data.URL)
	response:=struct{
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func redirectHandler(w http.ResponseWriter, r *http.Request){

	id:=r.URL.Path[len("/redirect/"):]
	url,err:=getURL(id)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusNotFound)
		return
	}

	http.Redirect(w,r,url.OriginalURL,http.StatusFound)




}


func main(){

		//Register the handler function to handle the requests to the root url ("/")
		http.HandleFunc("/",handler)
		http.HandleFunc("/shorten",shortenURLHandler)
		http.HandleFunc("/redirect/",redirectHandler)

	fmt.Println("starting the server on port 3000...")	
	err:=http.ListenAndServe(":3000",nil)
	if err!=nil{
		fmt.Println("Error starting server: ", err)
		return
	}



	




}