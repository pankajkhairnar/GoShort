package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

var baseUrl = "http://localhost:8080/" // Replace this url with your server goShort server url
var boltDBPath = "shortURL.db"
var shortUrlBkt = []byte("shortUrlBkt")
var seedChars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var seedCharsLen = len(seedChars)
var aChar byte = 97
var dbConn *bolt.DB

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Url    string `json:"url"`
}

func main() {
	var err error
	dbConn, err = bolt.Open(boltDBPath, 0644, nil)
	if err != nil {
		log.Println(err)
	}

	//defer dbConn.Close()
	router := httprouter.New()
	router.GET("/:code", Redirect)
	router.GET("/:code/json", GetOriginalURL)
	router.POST("/create/", Create)
	fmt.Println("Server started on :8080 port")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	urlStr := r.FormValue("url")
	urlStr = strings.Trim(urlStr, " ")
	if valid.IsURL(urlStr) == false {
		resp := &Response{Status: http.StatusBadRequest, Msg: "Invalid input URL", Url: ""}
		respJson, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respJson))
		return
	}

	newCode, err := GetNextCode()
	if err != nil {
		resp := Response{Status: http.StatusInternalServerError, Msg: "Some error occured while creating short URL", Url: ""}
		respJson, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respJson))
	}

	byteKey, byteUrl := []byte(newCode), []byte(urlStr)
	err = dbConn.Update(func(tx *bolt.Tx) error {
		//@todo : move this code to main function
		bucket, err := tx.CreateBucketIfNotExists(shortUrlBkt)
		if err != nil {
			return err
		}

		err = bucket.Put(byteKey, byteUrl)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Println(err)
		resp := &Response{Status: http.StatusInternalServerError, Msg: "Some error occured while creating short URL:", Url: ""}
		respJson, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respJson))
		return
	}

	shortUrl := baseUrl + newCode
	resp := &Response{Status: http.StatusOK, Msg: "Short URL created successfully", Url: shortUrl}
	respJson, _ := json.Marshal(resp)
	fmt.Fprint(w, string(respJson))
}

func Redirect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := ps.ByName("code")
	originalUrl, err := getCodeURL(code)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	http.Redirect(w, r, originalUrl, http.StatusFound)
}

func GetOriginalURL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := ps.ByName("code")
	originalUrl, err := getCodeURL(code)

	if err != nil {
		resp := &Response{Status: http.StatusInternalServerError, Msg: "Some error occured while reading URL", Url: ""}
		respJson, _ := json.Marshal(resp)
		fmt.Fprint(w, string(respJson))
		return
	}

	var resp *Response
	if len(originalUrl) != 0 {
		resp = &Response{Status: http.StatusOK, Msg: "Found", Url: originalUrl}
	} else {
		resp = &Response{Status: http.StatusNotFound, Msg: "URL not found", Url: ""}
	}

	respJson, err := json.Marshal(resp)

	if err != nil {
		fmt.Fprint(w, "Error occurred while creating json response")
		return
	}

	fmt.Fprint(w, string(respJson))
}

func getCodeURL(code string) (string, error) {
	key := []byte(code)
	var originalUrl string

	err := dbConn.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(shortUrlBkt)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", shortUrlBkt)
		}

		value := bucket.Get(key)
		originalUrl = string(value)
		return nil
	})

	if err != nil {
		return "", err
	}
	return originalUrl, nil
}

func GetNextCode() (string, error) {
	var newCode string
	err := dbConn.Update(func(tx *bolt.Tx) error {
		// by using locking on db file BoldDB makes sure it will be thread safe operation
		// and no two goroutines can can get same a short code at a time
		bucket, err := tx.CreateBucketIfNotExists(shortUrlBkt)
		if err != nil {
			return err
		}

		existingCodeByteKey := []byte("existingCodeKey")
		existingCode := bucket.Get(existingCodeByteKey)
		newCode, err = GenerateNextCode(string(existingCode))
		if err != nil {
			return err
		}

		err = bucket.Put(existingCodeByteKey, []byte(newCode))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	return newCode, nil
}

/*
	Following method is used to generate alphanumeric incremental code, which will be helpful
	for generating short urls
	this function will return new code like, input > output
	a > b, ax > ay, az > aA, aZ > a1, a9 > ba, 99 > aaa
	it will create shortest alphanumeric code possible for using in url
*/
func GenerateNextCode(code string) (string, error) {
	if code == "" {
		return string(aChar), nil
	}
	codeBytes := []byte(code)
	codeByteLen := len(codeBytes)

	codeCharIndex := -1
	for i := (codeByteLen - 1); i >= 0; i-- {
		codeCharIndex = bytes.IndexByte(seedChars, codeBytes[i])
		if codeCharIndex == -1 || codeCharIndex >= seedCharsLen {
			return "", errors.New("Invalid exisitng code")
		} else if codeCharIndex == (seedCharsLen - 1) {
			codeBytes[i] = aChar
		} else {
			codeBytes[i] = seedChars[(codeCharIndex + 1)]
			return string(codeBytes), nil
		}
	}
	for _, byteVal := range codeBytes {
		if byteVal != aChar {
			return string(codeBytes), nil
		}
	}
	// prepending "a" for generating new incremental code
	return "a" + string(codeBytes), nil
}
