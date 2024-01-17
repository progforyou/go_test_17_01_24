package data

import (
	"encoding/json"
	"gorm.io/gorm"
)

type Gender string

const (
	GENDER_FEMALE Gender = "female"
	GENDER_MALE   Gender = "male"
)

type Nationalize struct {
	CountryId   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type Person struct {
	gorm.Model
	Name               string        `json:"name"`
	Surname            string        `json:"surname"`
	Patronymic         string        `json:"patronymic"`
	Age                uint64        `json:"age"`
	Gender             Gender        `json:"gender"`
	NationalizeMarshal string        `json:"-"`
	Nationalize        []Nationalize `json:"nationalize" gorm:"-"`
}

func (p *Person) BeforeSave(tx *gorm.DB) error {
	bts, err := json.Marshal(p.Nationalize)
	if err != nil {
		return err
	}
	p.NationalizeMarshal = string(bts)
	return nil
}

func (p *Person) AfterFind(tx *gorm.DB) error {
	err := json.Unmarshal([]byte(p.NationalizeMarshal), &p.Nationalize)
	if err != nil {
		return err
	}
	return nil
}
