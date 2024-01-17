package controller

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"sync"
	"testing/service/person/data"
	"testing/service/person/web"
)

func Request(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msg("request")
		return nil, err
	}
	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("request body")
		return nil, err
	}
	return resBytes, nil
}

func RequestAge(obj *data.Person, wg *sync.WaitGroup) error {
	var result data.AgeRequest
	requestURL := fmt.Sprintf("https://api.agify.io/?name=%s", obj.Name)
	resultBytes, err := Request(requestURL)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		log.Error().Err(err).Msg("parse request body")
		return err
	}
	obj.Age = result.Age
	wg.Done()
	return nil
}
func RequestGender(obj *data.Person, wg *sync.WaitGroup) error {
	var result data.GenderRequest
	requestURL := fmt.Sprintf("https://api.genderize.io/?name=%s", obj.Name)
	resultBytes, err := Request(requestURL)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		log.Error().Err(err).Msg("parse request body")
		return err
	}
	obj.Gender = result.Gender
	wg.Done()
	return nil
}
func RequestNationalize(obj *data.Person, wg *sync.WaitGroup) error {
	var result data.NationalizeRequest
	requestURL := fmt.Sprintf("https://api.nationalize.io/?name=%s", obj.Name)
	resultBytes, err := Request(requestURL)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		log.Error().Err(err).Msg("parse request body")
		return err
	}
	obj.Nationalize = result.Country
	wg.Done()
	return nil
}
func EnrichmentData(obj *data.Person) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		err := RequestAge(obj, &wg)
		if err != nil {

		}
	}()
	go func() {
		err := RequestGender(obj, &wg)
		if err != nil {

		}
	}()
	go func() {
		err := RequestNationalize(obj, &wg)
		if err != nil {

		}
	}()
	wg.Wait()
}

func NewPersonController(db *gorm.DB, baseLog zerolog.Logger) web.CrudRouterController {
	log := baseLog.With().Str("model", "person").Logger()

	if err := db.AutoMigrate(&data.Person{}); err != nil {
		log.Fatal().Err(err).Msg("auto-migrate")
	}

	return web.CrudRouterController{
		Log: log,
		GetAll: func(p web.Pagination) interface{} {
			var res []data.Person
			tx := db.Model(&data.Person{})

			if p.Filter != nil {
				for k, v := range p.Filter {
					tx = tx.Where(fmt.Sprintf("%s = ?", k), v)
				}
			}
			tx = tx.Offset(p.Page * p.PerPage).Limit(p.PerPage)
			tx.Find(&res)
			return res
		},
		Create: func(r []byte) (interface{}, error) {
			var obj data.Person
			if err := json.Unmarshal(r, &obj); err != nil {
				log.Error().Err(err).Msg("decode json")
				return nil, err
			}
			EnrichmentData(&obj)

			tx := db.Model(&data.Person{}).Create(&obj)
			if tx.Error != nil {
				log.Error().Err(tx.Error).Msg("db error")
				return nil, tx.Error
			}
			return obj, nil
		},
		Update: func(id uint64, r []byte) (interface{}, error) {
			var obj data.Person
			if err := json.Unmarshal(r, &obj); err != nil {
				log.Error().Err(err).Msg("decode json")
				return nil, err
			}
			EnrichmentData(&obj)

			obj.ID = uint(id)
			tx := db.Save(&obj)
			if tx.Error != nil {
				log.Error().Err(tx.Error).Msg("db error")
				return nil, tx.Error
			}
			return obj, nil
		},
		Delete: func(id uint64) (error, int64) {
			tx := db.Delete(&data.Person{}, id)
			return tx.Error, tx.RowsAffected
		},
	}
}
