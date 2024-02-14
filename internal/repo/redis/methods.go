package redis

import (
	"bot/internal/entities"
	log2 "bot/internal/log"
	"context"
	"encoding/json"
	"log"
	"strconv"
)

func (r Redis) NewCartProduct(userID int64, idx int, product entities.Product) error {
	strID := strconv.Itoa(int(userID))

	bytes, err := json.Marshal(product)
	if err != nil {
		r.logger.Error("NewCartProduct: marshal error", log2.Fields{
			"error": err,
		})

		return err
	}

	r.Client.HSet(context.Background(), strID, idx, bytes)

	return nil
}

func (r Redis) CartLen(userID int64) (*int64, error) {
	strID := strconv.Itoa(int(userID))

	length, err := r.Client.HLen(context.Background(), strID).Result()
	if err != nil {
		r.logger.Error("CartLen: hlen failed", log2.Fields{
			"error": err,
		})

		return nil, err
	}

	return &length, nil
}

func (r Redis) GetCartProduct(userID int64, idx int) (*entities.Product, error) {
	strID := strconv.Itoa(int(userID))
	strIdx := strconv.Itoa(idx)

	res, err := r.Client.HGet(context.Background(), strID, strIdx).Result()
	if err != nil {
		r.logger.Error("GetCartProduct: hget failed", log2.Fields{
			"error": err,
		})

		return nil, err
	}

	var prod entities.Product

	err = json.Unmarshal([]byte(res), &prod)
	if err != nil {
		r.logger.Error("GetCartProduct: unmarshal failed", log2.Fields{
			"error": err,
		})

		return nil, err
	}

	return &prod, nil
}

func (r Redis) GetCart(userID int64) (map[int]entities.Product, error) {
	strID := strconv.Itoa(int(userID))
	res, err := r.Client.HGetAll(context.Background(), strID).Result()
	if err != nil {
		log.Println(err)
	}

	cart := make(map[int]entities.Product)

	for strKey, strValue := range res {
		var value entities.Product

		err = json.Unmarshal([]byte(strValue), &value)
		if err != nil {
			r.logger.Error("GetCart: unmarshal failed", log2.Fields{
				"error": err,
			})

			return nil, err
		}

		key, err := strconv.Atoi(strKey)
		if err != nil {
			r.logger.Error("GetCart: converting string to int failed", log2.Fields{
				"error": err,
			})

			return nil, err
		}

		cart[key] = value
	}

	return cart, nil
}

func (r Redis) ClearCart(userID int64) {
	strID := strconv.Itoa(int(userID))

	r.Client.Del(context.Background(), strID)
}

func (r Redis) DeleteProductFromCart(userID int64, idx int) {
	strID := strconv.Itoa(int(userID))
	strIdx := strconv.Itoa(idx)

	r.Client.HDel(context.Background(), strID, strIdx)
}
