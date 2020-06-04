package demodata

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	"webix/usermanager/data"
)

var demoDataFolder = "demodata"

func ImportDemoData(db *gorm.DB) {
	var c int
	db.Model(&data.User{}).Count(&c)

	if c > 0 {
		return
	}

	stat, err := os.Stat(demoDataFolder)
	if err != nil || !stat.IsDir() {
		return
	}

	importDemoStruct(db, data.Role{})
	importDemoStruct(db, data.RoleRule{})

	importDemoStruct(db, data.User{})
	importDemoStruct(db, data.UserRole{})
	importDemoStruct(db, data.UserRule{})

	importDemoStruct(db, data.Credential{})
	importDemoStruct(db, data.Log{})
}

func importDemoStruct(db *gorm.DB, t interface{}) {
	dt := reflect.TypeOf(t)
	name := strings.ToLower(dt.Name())
	cont, err := ioutil.ReadFile(filepath.Join(demoDataFolder, name+".json"))
	if err != nil {
		return
	}

	log.Println("[demo-data]", name)

	slicePtr := reflect.New(reflect.SliceOf(dt))
	slice := slicePtr.Interface()

	err = json.Unmarshal(cont, &slice)
	if err != nil {
		log.Fatal(err)
	}

	sliceObj := slicePtr.Elem()
	for i := 0; i < sliceObj.Len(); i++ {
		el := sliceObj.Index(i)
		err = db.Save(el.Addr().Interface()).Error
		if err != nil {
			log.Fatal(err)
		}
	}
}
