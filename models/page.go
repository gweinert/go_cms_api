package models

import "log"
import "fmt"

//Page used for somethign?
type Page struct {
	ID        int             `json:"id"`
	Title     string          `json:"title"`
	Path      string          `json:"path"`
	ParentID  int             `json:"parentId"`
	Name      string          `json:"name"`
	ShowInNav bool            `json:"showInNav"`
	SiteID    int             `json:"siteId"`
	SortOrder int             `json:"sortOrder"`
	Template  string          `json:"template"`
	Elements  []*Element      `json:"elements"`
	Groups    []*ElementGroup `json:"groups"`
}

//GetPagesBySiteID obviously
func GetPagesBySiteID(siteID int) ([]*Page, error) {
	pgs := make([]*Page, 0)
	rows, err := db.Query("SELECT * FROM pages WHERE siteid = $1", siteID)
	if err != nil {
		return pgs, err
	}
	defer rows.Close()

	for rows.Next() {
		p := new(Page)

		err := rows.Scan(&p.ID, &p.Title, &p.Path, &p.ParentID, &p.Name, &p.SiteID, &p.ShowInNav, &p.SortOrder, &p.Template)
		if err != nil {
			log.Fatal(err)
		}
		pgs = append(pgs, p)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return pgs, err
	}

	return pgs, nil
}

//GetPageByID gets a page by id
func GetPageByID(pageID int) (*Page, error) {
	p := new(Page)

	rows, err := db.Query("SELECT * FROM pages WHERE id = $1", pageID)
	if err != nil {
		return p, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&p.ID, &p.Title, &p.Path, &p.ParentID, &p.Name, &p.SiteID, &p.ShowInNav, &p.SortOrder, &p.Template)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return p, err
	}

	return p, nil
}

func SavePage(up *Page) (int, error) {
	stmt, err := db.Prepare(`UPDATE pages 
							SET name = $1, path = $2, sortorder = $3 
							WHERE id = $4 
							RETURNING id;`)
	if err != nil {
		fmt.Println("fail at prepare")
		log.Fatal(err)
	}
	res, err := stmt.Exec(up.Name, up.Path, up.SortOrder, up.ID)
	if err != nil {
		fmt.Println("fail at Exec")
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		fmt.Println("fail at rows affe")
		log.Fatal(err)
	}

	if len(up.Elements) > 0 {
		_, err := CreateOrUpdateElementIfExists(up.Elements)
		if err != nil {
			fmt.Println("fail at save el")
			log.Fatal(err)
		}
	}
	log.Printf("ID = %d, affected = %d\n", up.ID, rowCnt)
	return up.ID, nil
}
