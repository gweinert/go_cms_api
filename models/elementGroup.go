package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// GroupStructure saves structure of group elements so front end knows how to add more slides
type GroupStructure struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Amount  int    `json:"amount"`
	GroupID int    `json:"groupId"`
}

// ElementGroup model to be used around
type ElementGroup struct {
	ID        int               `json:"id"`
	PageID    int               `json:"pageId"`
	Name      string            `json:"name"`
	Structure []*GroupStructure `json:"structure"`
}

// GetGroupsByPageID gets all the element groups by page id
func GetGroupsByPageID(pageID int) ([]*ElementGroup, error) {
	grps := make([]*ElementGroup, 0)

	rows, err := db.Query(`SELECT * 
						FROM elementgroups 
						WHERE pageid = $1`, pageID)
	// INNER JOIN groupstructures on elementgroups.id = groupstructures.groupid
	if err != nil {
		return grps, err
	}
	defer rows.Close()

	for rows.Next() {
		eg := new(ElementGroup)
		// str := new(GroupStructure)

		// err := rows.Scan(&eg.ID, &eg.PageID, &eg.Name, &str.ID, &str.Type, &str.Amount, &str.GroupID)
		err := rows.Scan(&eg.ID, &eg.PageID, &eg.Name)

		if err != nil {
			log.Fatal(err)
		}
		strs, err := getGroupStructuresByElementGroupID(eg.ID)
		if err != nil {
			log.Fatal(err)
			fmt.Println("getgroupstructs fail")
			return grps, err
		}
		eg.Structure = strs
		grps = append(grps, eg)

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return grps, err
	}

	return grps, nil
}

func getGroupStructuresByElementGroupID(groupID int) ([]*GroupStructure, error) {
	strs := make([]*GroupStructure, 0)
	rows, err := db.Query(`SELECT * 
							FROM groupstructures 
							WHERE groupid = $1`, groupID)
	if err != nil {
		return strs, err
	}
	defer rows.Close()

	for rows.Next() {
		str := new(GroupStructure)

		err := rows.Scan(&str.ID, &str.Type, &str.Amount, &str.GroupID)
		if err != nil {
			log.Fatal(err)
		}

		strs = append(strs, str)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return strs, err
	}

	return strs, nil

}

// CreateNewGroup creates a new group for a page. Returns pagelements and group
func CreateNewGroup(g *ElementGroup) ([]*Element, *ElementGroup, error) {
	newGroup := new(ElementGroup)
	pgs := make([]*Element, 0)

	err := db.QueryRow(`INSERT INTO elementgroups(pageid, name)
						VALUES($1, $2) 
						RETURNING *`, g.PageID, g.Name).Scan(&newGroup.ID, &newGroup.PageID, &newGroup.Name)
	if err != nil {
		log.Fatal(err)
		return pgs, newGroup, err
	}

	//insert structure in
	ns, err := CreateNewStructures(g.Structure, newGroup.ID)
	if err != nil {
		log.Fatal(err)
		return pgs, newGroup, err
	}

	pgs, err = CreateNewPageElementsFromStructure(ns, newGroup.PageID)

	newGroup.Structure = ns

	return pgs, newGroup, nil
}

// CreateNewStructures inserts all structures from group into db
func CreateNewStructures(gs []*GroupStructure, groupID int) ([]*GroupStructure, error) {
	ns := make([]*GroupStructure, 0)

	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return ns, err
	}

	stmt, err := db.Prepare(`INSERT INTO groupstructures (type, amount, groupid)
							 VALUES ($1, $2, $3)`)

	for _, g := range gs {
		g.GroupID = groupID
		_, err := stmt.Exec(g.Type, g.Amount, groupID)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
		return ns, err
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
		return ns, err
	}

	return gs, nil
}

// CreateNewPageElementsFromStructure builds page elements from structures of group and inserts into db
func CreateNewPageElementsFromStructure(gs []*GroupStructure, pageID int) ([]*Element, error) {
	els := make([]*Element, 0)

	for index, g := range gs {
		for i := 0; i < g.Amount; i++ {
			p := new(Element)
			groupSortOrder := i + index
			name := strings.Join([]string{g.Type, strconv.Itoa(i)}, "")

			err := db.QueryRow(`INSERT INTO elements (pageid, groupid, type, sortorder, groupsortorder, name)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, pageid, groupid, groupsortorder, name, type
			`, pageID, g.GroupID, g.Type, 0, groupSortOrder, name).Scan(&p.ID, &p.PageID, &p.GroupID, &p.GroupSortOrder, &p.Name, &p.Type)
			if err != nil {
				log.Fatal(err)
			}

			els = append(els, p)
		}
	}

	return els, nil

}
