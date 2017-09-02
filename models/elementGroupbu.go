package models

import (
	"fmt"
)

func doNothing() {
	fmt.Println("do nothing")

}

// package models

// import (
// 	"log"
// 	"strconv"
// 	"strings"
// )

// // type GroupStructure struct {
// // 	ID      int    `json:"id"`
// // 	Type    string `json:"type"`
// // 	Amount  int    `json:"amount"`
// // 	GroupID int    `json:"groupId"`
// // }
// // type GroupStructure struct {
// // 	ID      int `json:"id"`
// // 	GroupID int `json:"goupid"`
// // 	Title   int `json:"title"`
// // 	Blurb   int `json:"blurb"`
// // 	Image   int `json:"image"`
// // 	Link    int `json:"link"`
// // }

// const ELEMENT_TYPES = []string{"title", "blurb", "image", "link"}

// // ElementGroup model to be used around
// type ElementGroup struct {
// 	ID        int             `json:"id"`
// 	PageID    int             `json:"pageId"`
// 	Name      string          `json:"name"`
// 	Structure *GroupStructure `json:"structure"`
// }

// // GetGroupsByPageID gets all the element groups by page id
// func GetGroupsByPageID(pageID int) ([]*ElementGroup, error) {
// 	grps := make([]*ElementGroup, 0)
// 	// grptrs := make([]*GroupStructure, 0)

// 	rows, err := db.Query(`SELECT *
// 	FROM elementgroups
// 	INNER JOIN groupstructures on elementgroups.id = groupstructures.groupid
// 	WHERE pageid = $1`, pageID)
// 	if err != nil {
// 		return grps, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		eg := new(ElementGroup)
// 		str := new(GroupStructure)

// 		err := rows.Scan(&eg.ID, &eg.PageID, &eg.Name, &str.ID, &str.Title, &str.Blurb, &str.Image, &str.Link)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		grps = append(grps, eg)

// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 		return grps, err
// 	}

// 	return grps, nil
// }

// // CreateNewGroup creates a new group for a page
// // returns groupid
// func CreateNewGroup(g *ElementGroup) ([]*Element, error) {
// 	newGroup := new(ElementGroup)
// 	pgs := make([]*Element, 0)

// 	err := db.QueryRow(`INSERT INTO elementgroups(pageid, name)
// 						VALUES($1, $2)
// 						RETURNING *`, g.PageID, g.Name).Scan(&newGroup.ID, &newGroup.PageID, &newGroup.Name)
// 	if err != nil {
// 		log.Fatal(err)
// 		return pgs, err
// 	}

// 	//insert structure in
// 	// ns, err := CreateNewStructures(g.Structure, newGroup.ID)
// 	ns, err := CreateNewStructure(g.Structure, newGroup.ID)
// 	if err != nil {
// 		log.Fatal(err)
// 		return pgs, err
// 	}

// 	// pgs, err = CreateNewPageElementsFromStructure(ns, newGroup.PageID)
// 	pes, err = CreateNewPageElementsFromStructure(ns, newGroup.PageID)

// 	newGroup.Structure = ns

// 	return pgs, nil
// }

// func CreateNewPageElementsFromStructure(ns *GroupStructure, pageID int) ([]*Element, error) {
// 	els := make([]*Element, 0)

// 	txn, err := db.Begin()
// 	if err != nil {
// 		log.Fatal(err)
// 		return els, err
// 	}

// 	stmt, err := db.Prepare(`INSERT INTO elements (pageid, groupid, type, sortorder, groupsortorder, name)
// 							VALUES ($1, $2, $3, $4, $5, $6)`)

// 	//make elements
// 	el
// 	for i := 0; i < ns.Title; i++ {
// 		name := strings.Join([]string{"title", strconv.Itoa(i)}, "")
// 		el := Element{
// 			PageID:         pageID,
// 			GroupID:        ns.GroupID,
// 			Type:           "title",
// 			SortOrder:      0,
// 			GroupSortOrder: i,
// 			Name:           name,
// 		}
// 		_, err := stmt.Exec(pageID, ns.GroupID, "title", 0, i, name)
// 		if err != nil {
// 			log.Fatal(err)
// 			return els, err
// 		}
// 		els = append(els, el)
// 	}

// 	err = stmt.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 		return els, err
// 	}

// 	err = txn.Commit()
// 	if err != nil {
// 		log.Fatal(err)
// 		return els, err
// 	}

// 	return els, nil
// }

// // CreateNewStructure returns a map of element types and their amounts
// func CreateNewStructure(str *GroupStructure, groupID int) (*GroupStructure, error) {
// 	newStr := new(GroupStructure)

// 	err := db.QueryRow(`INSERT INTO groupstructures(groupid, title, blurb, image, link)
// 						VALUES($1, $2, $3, $4, $5)
// 						RETURNING *`, str.GroupID, str.Title, str.Blurb, str.Image, str.Link).Scan(&newStr.ID, &newStr.GroupID, &newStr.Title, &newStr.Blurb, &newStr.Image, &newStr.Link)
// 	if err != nil {
// 		log.Fatal(err)
// 		return newStr, err
// 	}

// 	return newStr, nil
// }

// // CreateNewStructures returns an array of structures..didnt like this
// // func CreateNewStructures(gs []*GroupStructure, groupID int) ([]*GroupStructure, error) {
// // 	ns := make([]*GroupStructure, 0)

// // 	txn, err := db.Begin()
// // 	if err != nil {
// // 		log.Fatal(err)
// // 		return ns, err
// // 	}

// // 	stmt, err := db.Prepare(`INSERT INTO groupstructures (type, amount, groupid)
// // 							 VALUES ($1, $2, $3)`)

// // 	for _, g := range gs {
// // 		g.GroupID = groupID
// // 		_, err := stmt.Exec(g.Type, g.Amount, groupID)
// // 		if err != nil {
// // 			log.Fatal(err)
// // 		}
// // 	}

// // 	err = stmt.Close()
// // 	if err != nil {
// // 		log.Fatal(err)
// // 		return ns, err
// // 	}

// // 	err = txn.Commit()
// // 	if err != nil {
// // 		log.Fatal(err)
// // 		return ns, err
// // 	}

// // 	return gs, nil
// // }

// // func CreateNewPageElementsFromStructure(gs []*GroupStructure, pageID int) ([]*Element, error) {
// // 	els := make([]*Element, 0)

// // 	txn, err := db.Begin()
// // 	if err != nil {
// // 		log.Fatal(err)
// // 		return els, err
// // 	}

// // 	// stmt, err := db.Prepare(`INSERT INTO elements (type, amount, groupid)
// // 	// 						 VALUES ($1, $2, $3)`)
// // 	stmt, err := db.Prepare(`INSERT INTO elements (pageid, groupid, type, sortorder, groupsortorder, name)
// // 							VALUES ($1, $2, $3, $4, $5, $6)`)

// // 	for index, g := range gs {
// // 		for i := 0; i < g.Amount; i++ {
// // 			groupSortOrder := i + index
// // 			name := strings.Join([]string{g.Type, strconv.Itoa(i)}, "")
// // 			_, err := stmt.Exec(pageID, g.GroupID, g.Type, 0, groupSortOrder, name)
// // 			if err != nil {
// // 				log.Fatal(err)
// // 			}

// // 			pe := Element{
// // 				PageID:         pageID,
// // 				GroupID:        g.GroupID,
// // 				GroupSortOrder: groupSortOrder,
// // 				Name:           name,
// // 				Type:           g.Type,
// // 			}

// // 			els = append(els, &pe)
// // 		}
// // 	}

// // 	err = stmt.Close()
// // 	if err != nil {
// // 		log.Fatal(err)
// // 		return els, err
// // 	}

// // 	err = txn.Commit()
// // 	if err != nil {
// // 		log.Fatal(err)
// // 		return els, err
// // 	}

// // 	return els, nil

// // }
