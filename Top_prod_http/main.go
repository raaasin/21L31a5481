package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"

    //"github.com/gorilla/mux"
)

type Product struct {
    ProductName  string  `json:"productName"`
    Price        float64 `json:"price"`
    Rating       float64 `json:"rating"`
    Discount     int     `json:"discount"`
    Availability string  `json:"availability"`
}

const bearerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzIxOTc1NTQ2LCJpYXQiOjE3MjE5NzUyNDYsImlzcyI6IkFmZm9yZG1lZCIsImp0aSI6IjM1NjVlYzIxLWNiMTAtNGM1MC05MmUyLTI2YTg0MjE4ZDUwOSIsInN1YiI6Im5pc2FydnNrcEBnbWFpbC5jb20ifSwiY29tcGFueU5hbWUiOiJWaWduYW4iLCJjbGllbnRJRCI6IjM1NjVlYzIxLWNiMTAtNGM1MC05MmUyLTI2YTg0MjE4ZDUwOSIsImNsaWVudFNlY3JldCI6IkFWblJNcFFMdVdVc2tIbWIiLCJvd25lck5hbWUiOiJOSVNBUiBBSE1FRCBNT0hBTU1FRCIsIm93bmVyRW1haWwiOiJuaXNhcnZza3BAZ21haWwuY29tIiwicm9sbE5vIjoiMjFMMzFBNTQ4MSJ9.JpMdOWE1MfoaEWe1KV2XxRgx1A01fxlIJBW8N0dfNW8"

func fetchProductsFromAPI(company, category string, top int, minPrice, maxPrice float64) ([]Product, error) {
    url := "http://20.244.56.144/test/companies/" + company + "/categories/" + category + "/products?top=" + strconv.Itoa(top) + "&minPrice=" + strconv.FormatFloat(minPrice, 'f', -1, 64) + "&maxPrice=" + strconv.FormatFloat(maxPrice, 'f', -1, 64)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Add("Authorization", "Bearer "+bearerToken)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var products []Product
    err = json.Unmarshal(body, &products)
    if err != nil {
        return nil, err
    }

    return products, nil
}



func main() {
    r := mux.NewRouter()
    r.HandleFunc("/local/companies/{companyname}/categories/{categoryname}/products", getProductsHandler).Methods("GET")

    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
