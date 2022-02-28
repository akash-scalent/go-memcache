package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var cache = newCache("localhost:11211")

func main() {
	photo := &Photo{}
	photoItem, err := cache.Get("1")
	if err != nil {
		fmt.Println("From server")
		photo, err = getPhoto(1)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = json.Unmarshal(photoItem.Value, photo)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(photo)
}

func getPhoto(id int) (*Photo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/photos/%d", id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	photo := &Photo{}
	err = json.Unmarshal(body, photo)
	if err != nil {
		return nil, err
	}
	err = cache.Set(&memcache.Item{Key: fmt.Sprint(id), Value: body, Expiration: 10})
	if err != nil {
		return nil, err
	}
	return photo, nil
}
