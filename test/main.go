package main

import (
	"fmt"
	"log"
	"reflect"
)

type collection interface {
	fetchAll() error
	fetchWithTitle(string) error
	itemWithTitle(string) (interface{}, bool)
}

type collectionItem interface {
	getTitle() string
	collection
}

func fetchCollectionItemWithTitle(c collection, t string, i interface{}) error {
	item, ok := c.itemWithTitle(t)
	if !ok {
		return fmt.Errorf("item not found with title %s", t)
	}

	itemValue := reflect.ValueOf(item)
	iValue := reflect.ValueOf(i).Elem()
	if itemValue.Type() != iValue.Type() {
		return fmt.Errorf("found item type (%s) differs from passed item type (%s)", itemValue.Type(), iValue.Type())
	}

	iValue.Set(itemValue)

	return nil
}

func refreshCollectionItem(i collectionItem) error {
	i.fetchWithTitle(i.getTitle())

	return fetchCollectionItemWithTitle(i, i.getTitle(), i)
}

type roleType struct {
	title       string
	description string
	rolesType
}

func (r roleType) getTitle() string {
	return r.title
}

type rolesType []roleType

func (r *rolesType) fetchAll() error {
	*r = rolesType{
		{
			title:       "admin",
			description: "Admin role",
		},
		{
			title:       "user",
			description: "User role",
		},
	}

	return nil
}

func (r *rolesType) fetchWithTitle(t string) error {
	r.fetchAll()
	foundRole := roleType{}

	if err := fetchCollectionItemWithTitle(r, t, &foundRole); err != nil {
		return err
	}

	*r = rolesType{foundRole}

	return nil
}

func (r rolesType) itemWithTitle(title string) (interface{}, bool) {
	for _, role := range r {
		if role.title == title {
			return role, true
		}
	}

	return nil, false
}

func main() {
	roles := rolesType{}
	roles.fetchAll()

	roleName := "admin"
	foundRole := roleType{}
	if err := fetchCollectionItemWithTitle(&roles, roleName, &foundRole); err != nil {
		log.Fatalf("unable to find role %s: %s", roleName, err)
	}

	fmt.Printf("foundRole: %#v\n", foundRole)

	specificRole := roleType{
		title: "user",
		// rolesType: rolesType{
		// 	roleType{
		// 		title:       "user",
		// 		description: "User role",
		// 	},
		// },
	}
	if err := refreshCollectionItem(&specificRole); err != nil {
		log.Fatalf("unable to refresh role: %s: %s", specificRole.title, err)
	}

	fmt.Printf("specificRole: %#v\n", specificRole)
}
