package main

import (
	"day01/internal/dbreader"
	"day01/internal/model"
	"flag"
	"fmt"
	"log"
)

type changes struct {
	addedCake         []string
	removedCake       []string
	cakeChangedTime   []string
	addedIngredient   []string
	removedIngredient []string
	changedIngredient []string
}

func printChanges(ch changes) {
	for _, val := range ch.addedCake {
		fmt.Print(val)
	}
	for _, val := range ch.removedCake {
		fmt.Print(val)
	}
	for _, val := range ch.cakeChangedTime {
		fmt.Print(val)
	}
	for _, val := range ch.addedIngredient {
		fmt.Print(val)
	}
	for _, val := range ch.removedIngredient {
		fmt.Print(val)
	}
	for _, val := range ch.changedIngredient {
		fmt.Print(val)
	}
}

func findChanges(old, new model.Storage) *changes {
	var ch changes

	cakesOld := make(map[string]model.Cake)
	for _, cake := range old.Cake {
		cakesOld[cake.Name] = cake
	}
	cakesNew := make(map[string]model.Cake)
	for _, cake := range new.Cake {
		cakesNew[cake.Name] = cake
	}

	for nameOld := range cakesOld {
		_, ok := cakesNew[nameOld]
		if !ok {
			ch.removedCake = append(ch.removedCake, fmt.Sprintf("REMOVED cake \"%s\"\n", nameOld))
		}
	}

	for nameNew, valNew := range cakesNew {
		valOld, ok := cakesOld[nameNew]
		if !ok {
			ch.addedCake = append(ch.addedCake, fmt.Sprintf("ADDED cake \"%s\"\n", nameNew))
		} else {
			if valOld.Time != valNew.Time {
				ch.cakeChangedTime = append(ch.cakeChangedTime, fmt.Sprintf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", nameNew, valNew.Time, valOld.Time))
			}

			ingredientOld := make(map[string]model.Ingredient)
			for _, ingr := range valOld.Ingredients {
				ingredientOld[ingr.Ingredient_name] = ingr
			}

			ingredientNew := make(map[string]model.Ingredient)
			for _, ingr := range valNew.Ingredients {
				ingredientNew[ingr.Ingredient_name] = ingr
			}

			for ingr := range ingredientNew {
				_, ok := ingredientOld[ingr]
				if !ok {
					ch.addedIngredient = append(ch.addedIngredient, fmt.Sprintf("ADDED ingredient \"%s\" for cake \"%s\"\n", ingr, valNew.Name))
				}
			}

			for ingr, ingrOld := range ingredientOld {
				ingrNew, ok := ingredientNew[ingr]
				if !ok {
					ch.removedIngredient = append(ch.removedIngredient, fmt.Sprintf("REMOVED ingredient \"%s\" for cake \"%s\"\n", ingr, valNew.Name))
				} else {
					if ingrOld.Ingredient_unit != "" && ingrNew.Ingredient_unit == "" {
						ch.changedIngredient = append(ch.changedIngredient, fmt.Sprintf("REMOVED uint \"%s\" for ingredient \"%s\" of cake \"%s\"\n", ingrOld.Ingredient_unit, ingr, valNew.Name))
					} else if ingrOld.Ingredient_unit != ingrNew.Ingredient_unit {
						ch.changedIngredient = append(ch.changedIngredient, fmt.Sprintf("CHANGED uint for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", ingr, valNew.Name, ingrNew.Ingredient_unit, ingrOld.Ingredient_unit))
					} else if ingrOld.Ingredient_count != ingrNew.Ingredient_count {
						ch.changedIngredient = append(ch.changedIngredient, fmt.Sprintf("CHANGED uint count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", ingr, valNew.Name, ingrNew.Ingredient_count, ingrOld.Ingredient_count))
					}
				}
			}
		}
	}
	return &ch
}

func getData() (*model.Storage, *model.Storage, error) {
	oldFile := flag.String("old", "", "old db file name")
	newFile := flag.String("new", "", "new db file name")
	flag.Parse()

	oldReader, err := dbreader.NewReader(*oldFile)
	if err != nil {
		return nil, nil, err
	}
	oldDB, err := oldReader.Read()
	if err != nil {
		return nil, nil, err
	}

	newReader, err := dbreader.NewReader(*newFile)
	if err != nil {
		return nil, nil, err
	}
	newDB, err := newReader.Read()
	if err != nil {
		return nil, nil, err
	}

	return oldDB, newDB, nil
}

func main() {
	oldDB, newDB, err := getData()
	if err != nil {
		log.Fatal(err)
	}

	ch := findChanges(*oldDB, *newDB)
	printChanges(*ch)
}
