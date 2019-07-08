package products

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func fromFile(fileName string) (*Menu, error) {
	var d *Menu
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println(err)
		return d, err
	}
	err = json.Unmarshal(dat, &d)
	if err != nil {
		log.Println(err)
		return d, err
	}
	return d, err
}

//
func fromDB() {

}
