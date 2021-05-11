package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

//helper to check if a []string contains a certain value
func ContainsString(s []string, val string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}

type Food struct {
	ID          string //represented as a date string generated by time.Now().Format()
	Name        string
	Label       []string
	Ingredients []string
}

func GenerateFoodID() string {
	return time.Now().Format("Mon-Jan-2-15:04:05-2006")
}

func (f Food) Equals(comp Food) bool { //compare foods by ID
	return f.ID == comp.ID
}

func (f Food) HasLabel(label []string) bool { //check if a food has a Label set to true
	ret := true
	for _, v := range label {
		if !ContainsString(f.Label, v) {
			ret = false
		}
	}
	return ret
}

//removes all Labels that are not in the Label list, and returns true if at least one label was removed
func (f *Food) ValidateLabel() (bool, error) {
	Label, err := GetWholeLabelList()
	if err != nil {
		return false, err
	}
	ret := false
	for i := 0; i < len(f.Label); i++ {
		if !ContainsString(Label, f.Label[i]) {
			var newLabel []string = make([]string, 0, len(f.Label))
			for _, v := range f.Label {
				if v != f.Label[i] {
					newLabel = append(newLabel, v)
				}
			}
			f.Label = newLabel
			i--
			ret = true
		}
	}
	return ret, nil
}

type JsonHolder struct {
	LabelList      []string
	IngredientList []string
	FoodList       []Food
}

func writeJsonHolder(holder JsonHolder) error {
	newFile, err := json.MarshalIndent(holder, "", "\t") //encode the holder to json
	if err != nil {
		return err
	}
	return ioutil.WriteFile("resources/foods.json", newFile, 0644) //and write it to foods.json
}

func writeFoodList(foods []Food) error {
	holder, err := GetWholeFile() //get the whole file
	if err != nil {
		return err
	}
	holder.FoodList = foods        //change its food list
	return writeJsonHolder(holder) //and write it to the file
}

func GetWholeFile() (JsonHolder, error) {
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

func GetEveryFoodWithLabel(label []string) ([]Food, error) {
	Foods, err := GetWholeFoodList()
	if err != nil {
		return nil, err
	}
	var retFoods []Food
	for _, v := range Foods {
		if v.HasLabel(label) {
			retFoods = append(retFoods, v)
		}
	}
	return retFoods, nil
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
	Foods = append(Foods, data)
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
	return writeJsonHolder(holder)
}

func DeleteLabelFromList(data string) error {
	holder, err := GetWholeFile()
	if err != nil {
		return err
	}
	Label := holder.LabelList
	holder.LabelList = make([]string, 0, len(Label))
	for _, v := range Label { //write every label to the new list that is not equal to the label that we want to delete
		if v != data {
			holder.LabelList = append(holder.LabelList, v)
		}
	}
	for i := range holder.FoodList { //delete the label from every food in the list
		if holder.FoodList[i].HasLabel([]string{data}) {
			var newLabel []string
			for _, v := range holder.FoodList[i].Label {
				if v != data {
					newLabel = append(newLabel, v)
				}
			}
			holder.FoodList[i].Label = newLabel
		}
	}
	return writeJsonHolder(holder)
}
