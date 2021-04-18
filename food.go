package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

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

type Food struct {
	ID    string //represented as a date string generated by time.Now().Format()
	Name  string
	Label map[string]bool
}

func GenerateFoodID() string {
	return time.Now().Format("Mon-Jan-2-15:04:05-2006")
}

func (f Food) Equals(comp Food) bool { //compare foods by ID
	return f.ID == comp.ID
}

func (f Food) HasLabel(label string) bool { //check if a food has a Label set to true
	return f.Label[label]
}

type JsonHolder struct {
	LabelList []string
	FoodList  []Food
}

func writeJsonHolder(holder JsonHolder) error {
	newFile, err := json.MarshalIndent(holder, "", "\t") //encode the holder to json
	if err != nil {
		return err
	}
	return ioutil.WriteFile("resources/foods.json", newFile, 0644) //and write it to foods.json
}

func writeFoodList(foods []Food) error {
	holder, err := GetWholeFile()
	if err != nil {
		return err
	}
	holder.FoodList = foods
	return writeJsonHolder(holder)
}

//at the moment unused but maybe usefull later
func writeLabelList(label []string) error {
	holder, err := GetWholeFile()
	if err != nil {
		return err
	}
	holder.LabelList = label
	return writeJsonHolder(holder)
}

func GetWholeFile() (JsonHolder, error) {
	//if foods.json does not exist create it with an empty food and label list
	if b, err := CheckFileExist("resources/foods.json"); err != nil {
		log.Println("File resources/foods.json may or may not exist: " + err.Error())
	} else if !*b {
		os.Create("resources/foods.json")
		ioutil.WriteFile("resources/foods.json", []byte(`{"LabelList":[],"FoodLIst":[]}`), 0644)
	}
	file, err := ioutil.ReadFile("resources/foods.json") //read foods.json
	if err != nil {
		return JsonHolder{}, err
	}
	var holder JsonHolder
	err = json.Unmarshal(file, &holder) //obtain the file as a JsonHolder
	if err != nil {
		return JsonHolder{}, err
	}
	return holder, nil
}

//returns a []Food containing the food list in foods.json
func GetWholeFoodList() ([]Food, error) {
	holder, err := GetWholeFile()
	return holder.FoodList, err
}

//returns a []string containing the label list in foods.json
func GetWholeLabelList() ([]string, error) {
	holder, err := GetWholeFile()
	return holder.LabelList, err
}

//add a food to the list in foods.json
func AddFoodToList(data Food) error {
	Foods, err := GetWholeFoodList()
	if err != nil {
		return err
	}
	Foods = append(Foods, data) //append the new Food to the list
	return writeFoodList(Foods)
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
	return writeFoodList(Foods)
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
	return writeFoodList(updatedFoods)
}

//add a label to every food in the list
func AddLabelToList(label string) error {
	holder, err := GetWholeFile()
	if err != nil {
		return err
	}
	holder.LabelList = append(holder.LabelList, label) //append the label to the LabelList
	for i := range holder.FoodList {                   //and add the label to every food
		holder.FoodList[i].Label[label] = false
	}
	return writeJsonHolder(holder)
}
