package persist

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

type Tuido struct {
	Lists []TuidoList `yaml:"tuido"`
}

type TuidoList struct {
	Title string `yaml:"title"`
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Created     CustomTime `yaml:"created"`
	Updated     CustomTime `yaml:"updated"`
	Index       int        `yaml:"index"`
	Title       string     `yaml:"title"`
	Description string     `yaml:"description"`
}

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02T15:04:05.000-0700"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func LoadTuido(tuidoFile string) *Tuido {
	var t Tuido
	_, err := os.Stat(tuidoFile)
	if errors.Is(err, os.ErrNotExist) {
		return &t
	}

	yamlData, err := os.ReadFile(tuidoFile)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
		return &t
	}

	err = yaml.Unmarshal([]byte(yamlData), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &t
}

func (t *Tuido) Persist(tuidoFile string) {
	yamlData, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = os.WriteFile(tuidoFile, yamlData, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
}
