package entity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Nationality struct {
	CountryID   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type Person struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Surname     string        `json:"surname"`
	Patronymic  string        `json:"patronymic"`
	Age         int           `json:"age"`
	Gender      string        `json:"gender"`
	Nationality []Nationality `json:"country"`
}

func (p *Person) EnrichPersonInfo() error {

	age, err := getAge(p.Name)
	if err != nil {
		return err
	}
	p.Age = age

	gender, err := getGender(p.Name)
	if err != nil {
		return err
	}
	p.Gender = gender

	nationality, err := getNationality(p.Name)
	if err != nil {
		return err
	}
	p.Nationality = nationality

	return nil
}

// Функция для получения возраста из внешнего API
func getAge(firstName string) (int, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", firstName)
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var ageData map[string]interface{}
	err = json.Unmarshal(body, &ageData)
	if err != nil {
		return 0, err
	}

	age := int(ageData["age"].(float64))
	return age, nil
}

// Функция для получения пола из внешнего API
func getGender(firstName string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", firstName)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var genderData map[string]interface{}
	err = json.Unmarshal(body, &genderData)
	if err != nil {
		return "", err
	}

	gender := genderData["gender"].(string)

	return gender, nil
}

// Функция для получения национальности из внешнего API
func getNationality(firstName string) (res []Nationality, err error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", firstName)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var nationalityData Person
	err = json.Unmarshal(body, &nationalityData)
	if err != nil {
		return nil, err
	}

	return nationalityData.Nationality, nil
}
