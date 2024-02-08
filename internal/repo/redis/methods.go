package redis

import (
	"bot/internal/entities"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

func (r Redis) NewCartProduct(userID int64, idx int, product entities.Product) {
	strID := strconv.Itoa(int(userID))

	bytes, err := json.Marshal(product)
	if err != nil {
		log.Println(err)
	}

	r.Client.HSet(context.Background(), strID, idx, bytes)
}

func (r Redis) CartLen(userID int64) int64 {
	strID := strconv.Itoa(int(userID))

	length, err := r.Client.HLen(context.Background(), strID).Result()
	if err != nil {
		log.Println(err)
	}

	return length
}

func (r Redis) GetCartProduct(userID int64, idx int) entities.Product {
	strID := strconv.Itoa(int(userID))
	strIdx := strconv.Itoa(idx)

	res, err := r.Client.HGet(context.Background(), strID, strIdx).Result()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(res)

	var prod entities.Product

	err = json.Unmarshal([]byte(res), &prod)
	if err != nil {
		log.Println(err)
	}

	return prod
}

func (r Redis) GetCart(userID int64) []entities.Product {
	strID := strconv.Itoa(int(userID))
	res, err := r.Client.HGetAll(context.Background(), strID).Result()
	if err != nil {
		log.Println(err)
	}

	var cart []entities.Product

	for _, strValue := range res {
		var value entities.Product

		err = json.Unmarshal([]byte(strValue), &value)
		if err != nil {
			log.Println(err)
		}

		cart = append(cart, value)
	}

	return cart
}

func (r Redis) ClearCart(userID int64) {
	strID := strconv.Itoa(int(userID))

	r.Client.Del(context.Background(), strID)
}
