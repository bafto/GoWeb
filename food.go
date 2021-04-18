package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Food struct {
	ID    string //represented as a date string generated by time.Now().Format()
	Name  string
	Label map[string]bool
}

func NewFood(food Food) *Food { //genearte an ID for a new Food
	f := food
	f.ID = time.Now().Format("Mon-Jan-2-15:04:05-2006")
	return &f
}

func (f Food) Equals(comp Food) bool { //compare foods by ID
	return f.ID == comp.ID
}

func (f Food) HasLabel(label string) bool { //check if a food has a Label set to true
	return f.Label[label]
}

//check if a file exists (mainly used for food.json)
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

//add a food to the list in foods.json
func AddFoodToList(data Food) error {
	Foods, err := GetWholeFoodList()
	if err != nil {
		return err
	}
	Foods = append(Foods, data)                         //append the new Food to the list
	newFile, err := json.MarshalIndent(Foods, "", "\t") //and encode the food list to json
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("resources/foods.json", newFile, 0644) //finally write the new json to the file
	if err != nil {
		return err
	}
	return nil
}

//change labels of a food in the list
func ChangeFoodInList(data Food) error {
	Foods, err := GetWholeFoodList()
	if err != nil {
		return err
	}
	for i, v := range Foods { //loop through the list and overwrite the food to change
		if v.Equals(data) {
			Foods[i] = data
			break
		}
	}
	newFile, err := json.MarshalIndent(Foods, "", "\t") //decode the list into json
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("resources/foods.json", newFile, 0644) //finally write the new list to foods.json
	if err != nil {
		return err
	}
	return nil
}

//remove a food from foods.json
func DeleteFoodFromList(data Food) error {
	Foods, err := GetWholeFoodList()
	if err != nil {
		return err
	}
	var updatedFoods []Food = make([]Food, 0, cap(Foods)) //write every food to the new list that is not equal to the food that we want to delete
	for _, v := range Foods {
		if !v.Equals(data) {
			updatedFoods = append(updatedFoods, v)
		}
	}
	newFile, err := json.MarshalIndent(updatedFoods, "", "\t") //decode list into json
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("resources/foods.json", newFile, 0644) //finally write the list to foods.json
	if err != nil {
		return err
	}
	return nil
}

//returns a []Food containing the food list in foods.json
func GetWholeFoodList() ([]Food, error) {
	//if foods.json does not exist create it with an empty food list
	if b, err := CheckFileExist("resources/foods.json"); err != nil {
		log.Println("File resources/foods.json may or may not exist: " + err.Error())
	} else if !*b {
		os.Create("resources/foods.json")
		ioutil.WriteFile("resources/foods.json", []byte("[]"), 0644)
	}
	file, err := ioutil.ReadFile("resources/foods.json") //read foods.json
	if err != nil {
		return nil, err
	}
	var Foods []Food
	err = json.Unmarshal(file, &Foods) //obtain the food list
	if err != nil {
		return nil, err
	}
	return Foods, nil
}

//add a label to every food in the list
func AddLabelToList(label string) error {
	Foods, err := GetWholeFoodList()
	if err != nil {
		return err
	}
	for i := range Foods { //loop through the list and add the label
		Foods[i].Label[label] = false
	}
	newFile, err := json.MarshalIndent(Foods, "", "\t") //decode the new list into json
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("resources/foods.json", newFile, 0644) //write the new list into foods.json
	if err != nil {
		return err
	}
	return nil
}
