package todo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const initialSize = 5

type Record struct {
	Number      int       `json:"number"`
	Task        string    `json:"task"`
	IsDone      bool      `json:"isDone"`
	CreatedAt   time.Time `json:"createdAt"`
	CompletedAt time.Time `json:"completedAt"`
}

type List struct {
	Name    string   `json:"name"`
	Records []Record `json:"records"`
}

func (l *List) Add(record string) error {
	next := len(l.Records) + 1
	newRecord := Record{
		Number:      next,
		Task:        record,
		IsDone:      false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	l.Records = append(l.Records, newRecord)
	return nil
}

func (l *List) Remove(number int) error {
	l.Records = append(l.Records[:number-1], l.Records[number:]...)
	return nil
}

func (l *List) MarkAsDone(id int) error {
	if id <= 0 || id >= len(l.Records) {
		return errors.New("no such element found")
	}
	record := &l.Records[id-1]
	record.IsDone = true
	record.CreatedAt = time.Now()
	return nil
}

func NewList(name string) *List {
	list := List{
		Records: make([]Record, initialSize),
		Name:    name,
	}
	return &list
}

func LoadFromFile(name string) (*List, error) {
	file, err := os.Open(fmt.Sprintf("%s.json", name))
	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			panic(e)
		}
	}(file)

	if err != nil {
		return nil, err
	}
	slice := make([]byte, 1024)
	buffer := bytes.Buffer{}
	for {
		end, e := file.Read(slice)
		if e != nil && e != io.EOF {
			return nil, e
		}
		buffer.Write(slice[:end])
		if end == 0 {
			break
		}
	}
	var list List
	output := string(buffer.Bytes())
	if er := json.Unmarshal([]byte(output), &list); er != nil {
		return nil, er
	}
	return &list, nil
}

func SaveToFile(l *List) error {
	data, err := json.Marshal(l)
	if err != nil {
		log.Fatal(err)
	}
	return os.WriteFile(fmt.Sprintf("%s.json", l.Name), data, 0644)
}
