package services

import (
	"api/src/models"
	"fmt"
	"log"
	"strings"
)

type CurrencySearchService interface {
	GetAllUpdatableCurrencies() ([]models.Currency, error)
	GetCurrenciesBasedOnUSDFromAPI(string, []string) ([]models.ConversionRateFromAPI, error)
}

type SyncRepository interface {
	UpdateCurrency(models.Currency) error
	InsertCurrency(models.Currency) error
	DeleteCurrency(string) error
}

type SyncService struct {
	repository    SyncRepository
	searchService CurrencySearchService
}

func NewSyncService(repository SyncRepository, currencyService CurrencySearchService) *SyncService {
	return &SyncService{repository, currencyService}
}

func (syncService SyncService) UpdateAllUpdatableCurrencies() {
	fmt.Println("##### NEW JOB RUN #####")

	currencies, err := syncService.searchService.GetAllUpdatableCurrencies()
	if err != nil {
		log.Fatal(err)
	}

	var currencyNames []string
	for _, currency := range currencies {
		currencyNames = append(currencyNames, currency.Name)
	}

	conversionRatesByCurrency, err := syncService.searchService.GetCurrenciesBasedOnUSDFromAPI("USD",
		currencyNames)
	if err != nil {
		log.Fatal(err)
	}

	for _, newConversionRate := range conversionRatesByCurrency {

		newCurrency := models.Currency{
			Name:            newConversionRate.Name,
			ConversionRate:  newConversionRate.ConversionRate,
			IsAutoUpdatable: true,
		}

		if err := syncService.repository.UpdateCurrency(newCurrency); err != nil {
			log.Fatal(err)
		}
	}
}

func (syncService SyncService) InsertCurrency(currency models.Currency) error {

	currency.Name = strings.ToUpper(currency.Name)

	if err := syncService.repository.InsertCurrency(currency); err != nil {
		return err
	}

	return nil
}

func (syncService SyncService) DeleteCurrency(currencyName string) error {

	if err := syncService.repository.DeleteCurrency(currencyName); err != nil {
		return err
	}

	return nil
}

func (syncService SyncService) UpdateCurrency(currency models.Currency) error {

	if err := syncService.repository.UpdateCurrency(currency); err != nil {
		return err
	}

	return nil
}