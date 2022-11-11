package repositories

import (
	"api/src/config"
	"api/src/database"
	"api/src/models"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

func GetAllUpdatableCurrencies() ([]models.Currency, error) {
	redisClient := database.Connect()
	defer redisClient.ClientKill(config.DBPort)

	var currencies []models.Currency

	dbResultIterator := redisClient.Scan(0, "*", 0).Iterator()

	for i := 0; dbResultIterator.Next(); i++ {

		currencyFromDatabase, err := GetCurrencyByName(dbResultIterator.Val())
		if err != nil {
			return nil, err
		}

		if currencyFromDatabase.IsAutoUpdatable {
			currencies = append(currencies, currencyFromDatabase)
		}
	}

	return currencies, nil
}

func GetCurrencyByName(currencyName string) (models.Currency, error) {

	redisClient := database.Connect()
	defer redisClient.ClientKill(config.DBPort)

	var currency models.Currency

	dbResultJSON, err := redisClient.Get(currencyName).Result()

	if err != nil {
		fmt.Println("error getting currency from database:", err)
		return models.Currency{}, err
	}

	err = json.Unmarshal([]byte(dbResultJSON), &currency)
	if err != nil {
		fmt.Println("error unmarshalling dbResultJSON:", err)
	}

	return currency, nil
}

func InsertCurrency(currency models.Currency) error {

	redisClient := database.Connect()
	defer redisClient.ClientKill(config.DBPort)

	currencyJSON, err := json.Marshal(currency)

	if err != nil {
		return err
	}

	err = redisClient.Set(currency.Name, currencyJSON, 0).Err()

	if err != nil {
		return err
	}

	return nil
}

func DeleteCurrency(currencyName string) *redis.IntCmd {
	redisClient := database.Connect()
	defer redisClient.ClientKill(config.DBPort)

	return redisClient.Del(currencyName)
}
