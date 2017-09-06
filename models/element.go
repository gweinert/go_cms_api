package models

import (
	"fmt"
	"log"
	"strings"
)

type Element struct {
	ID             int    `json:"id"`
	PageID         int    `json:"pageId"`
	GroupID        int    `json:"groupId"`
	SortOrder      int    `json:"sortOrder"`
	GroupSortOrder int    `json:"groupSortOrder"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	Body           string `json:"body"`
	ImageURL       string `json:"imageURL"`
	LinkPath       string `json:"linkPath"`
	LinkText       string `json:"linkText"`
}

//GetElementsByGroupID returns all elements from group
func GetElementsByGroupID(groupID int) ([]*Element, error) {
	els, err := getElementsBy("groupid", groupID)
	if err != nil {
		return nil, err
	}

	return els, nil
}

//GetElementsByPageID returns all elements on page
func GetElementsByPageID(pageID int) ([]*Element, error) {
	els, err := getElementsBy("pageid", pageID)
	if err != nil {
		return nil, err
	}

	return els, nil
}

func getElementsBy(col string, colVar int) ([]*Element, error) {
	els := make([]*Element, 0)

	queryStr := strings.Join([]string{"SELECT * FROM elements WHERE", col, "= $1 ORDER BY sortorder"}, " ")
	rows, err := db.Query(queryStr, colVar)
	if err != nil {
		return els, err
	}
	defer rows.Close()

	for rows.Next() {
		el := new(Element)

		err := rows.Scan(&el.ID, &el.PageID, &el.GroupID, &el.SortOrder, &el.GroupSortOrder, &el.Name, &el.Type, &el.Body, &el.ImageURL, &el.LinkPath, &el.LinkText)
		if err != nil {
			log.Fatal(err)
		}
		els = append(els, el)

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return els, err
	}

	return els, nil
}

//CreateOrUpdateElementIfExists prepares two statements depending on if update succeeded, inserts new page element
func CreateOrUpdateElementIfExists(els []*Element) ([]*Element, error) {
	fmt.Println("CreateOrUpdateElementIfExists")

	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	updStmt, err := db.Prepare(`UPDATE elements
							SET sortorder = $1, groupsortorder = $2, name = $3, body = $4, imageurl = $5, linkpath = $6, linktext = $7
							WHERE id = $8;
							`)
	insStmt, err := db.Prepare(`INSERT INTO elements (pageid, groupid, type, sortorder, groupsortorder, name, body, imageurl, linkpath, linktext)
							 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
							 `)

	for _, el := range els {
		fmt.Printf("%v\n", el)
		res, err := updStmt.Exec(el.SortOrder, el.GroupSortOrder, el.Name, el.Body, el.ImageURL, el.LinkPath, el.LinkText, el.ID)

		if err != nil {
			fmt.Println("fail Exec")

			log.Fatal(err)
		}

		rowCnt, err := res.RowsAffected()
		if err != nil {
			fmt.Println("fail at rows affect")
			log.Fatal(err)
		} else if rowCnt == 0 {
			fmt.Println("inserting elements...")

			// @DEV needs to return ids yo
			// res, err = insStmt.Exec(el.PageID, el.GroupID, el.Type, el.SortOrder, el.GroupSortOrder, el.Name, el.Body, el.ImageURL, el.LinkPath, el.LinkText)
			if err != nil {
				fmt.Println("fail inserExec")

				log.Fatal(err)
			}
		}

	}

	err = insStmt.Close()
	if err != nil {
		fmt.Println("close")

		log.Fatal(err)
	}

	err = updStmt.Close()
	if err != nil {
		fmt.Println("close")

		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		fmt.Println("cpommit")
		log.Fatal(err)
		return 1, err
	}

	return 1, nil
}

//DeleteElements given an array of element ids, deletes page elements and returns array of ids
func DeleteElements(ids []int) ([]int, error) {

	stmt, _ := db.Prepare(`DELETE from elements WHERE id = $1`)
	for _, id := range ids {
		_, err := stmt.Query(id)
		if err != nil {
			log.Fatal(err)
			fmt.Println("query fail")
			return nil, err
		}
	}

	return ids, nil
}
