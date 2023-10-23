package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Define data models
type Person struct {
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}

type PersonService struct {
	db *sql.DB
}

func (ps *PersonService) Init() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Open a connection to the MySQL database
	var err error
	ps.db, err = sql.Open("mysql", "root:1423@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
}

func (ps *PersonService) enrichPersonData(name string) (int, string, string) {
	age, gender, nationality := 0, "Unknown", "Unknown"

	// Fetch data from an API
	apiEndpoints := []struct {
		name   string
		apiURL string
	}{
		{"Age", fmt.Sprintf("https://api.agify.io/?name=%s", name)},
		{"Gender", fmt.Sprintf("https://api.genderize.io/?name=%s", name)},
		{"Nationality", fmt.Sprintf("https://api.nationalize.io/?name=%s", name)},
	}

	for _, endpoint := range apiEndpoints {
		resp, err := http.Get(endpoint.apiURL)
		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}
		defer resp.Body.Close()

		var responseData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			continue
		}

		switch endpoint.name {
		case "Age":
			if v, ok := responseData["age"].(float64); ok {
				age = int(v)
			}
		case "Gender":
			if v, ok := responseData["gender"].(string); ok {
				gender = v
			}
		case "Nationality":
			countries, ok := responseData["country"].([]interface{})
			if ok && len(countries) > 0 {
				country := countries[0].(map[string]interface{})
				if v, ok := country["country_id"].(string); ok {
					nationality = v
				}
			}
		}
	}

	return age, gender, nationality
}

func main() {
	ps := &PersonService{}
	ps.Init()

	// Example usage
	name := "Dmitriy"
	surname := "Ushakov"
	patronymic := "Vasilevich"

	// Enrich the data using external APIs
	age, gender, nationality := ps.enrichPersonData(name)

	// Insert the person data into the People table
	insertQuery := "INSERT INTO People (Name, Surname, Patronymic, Age, Gender, Nationality) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := ps.db.Exec(insertQuery, name, surname, patronymic, age, gender, nationality)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Added person: Name: %s, Surname: %s, Patronymic: %s, Age: %d, Gender: %s, Nationality: %s\n", name, surname, patronymic, age, gender, nationality)
}
