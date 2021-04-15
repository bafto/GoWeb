package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Food struct {
	Name string
}

func CheckFileExist(file string) (*bool, error) {
	var b bool
	if _, err := os.Stat(file); err == nil {
		b = true
		return &b, nil
	} else if os.IsNotExist(err) {
		b = false
		return &b, nil
	} else {
		return nil, err
	}
}

//add data to foods.json
func AddFoodToList(data Food) error {
	if b, err := CheckFileExist("static/foods.json"); err != nil {
		log.Println("File static/foods.json may or may not exist: " + err.Error())
	} else if !*b {
		os.Create("static/foods.json")
		ioutil.WriteFile("static/foods.json", []byte("[]"), 0644)
	}
	file, err := ioutil.ReadFile("static/foods.json")
	if err != nil {
		return err
	}
	var Foods []Food
	err = json.Unmarshal(file, &Foods)
	if err != nil {
		return err
	}
	Foods = append(Foods, data)
	newFile, err := json.MarshalIndent(Foods, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("static/foods.json", newFile, 0644)
	if err != nil {
		return err
	}
	return nil
}

//remove data from foods.json
func DeleteFoodFromList(data Food) error {
	file, err := ioutil.ReadFile("static/foods.json")
	if err != nil {
		return err
	}
	var Foods []Food
	err = json.Unmarshal(file, &Foods)
	if err != nil {
		return err
	}
	var updatedFoods []Food = make([]Food, 0, cap(Foods))
	for _, v := range Foods {
		if v != data {
			updatedFoods = append(updatedFoods, v)
		}
	}
	newFile, err := json.MarshalIndent(updatedFoods, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("static/foods.json", newFile, 0644)
	if err != nil {
		return err
	}
	return nil
}
