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
	// ElementsGroups [][]*Element    `json:"elementsGroups"`
	ElementsGroups map[int][][]*Element `json:"elementsGroups"`
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

//CreateNewPage inserts new page into db.
func CreateNewPage(p *Page) (*Page, error) {
	var np = new(Page)
	fmt.Printf("%v", p)
	err := db.QueryRow(`INSERT INTO pages (title, path, parentid, name, siteid, showinnav, sortorder, template)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`, &p.Title, &p.Path, &p.ParentID, &p.Name, &p.SiteID, &p.ShowInNav, &p.SortOrder, &p.Template,
	).Scan(&np.ID, &np.Title, &np.Path, &np.ParentID, &np.Name, &np.SiteID, &np.ShowInNav, &np.SortOrder, &np.Template)

	if err != nil {
		log.Fatal(err)
		return np, err
	}

	return np, nil
}

//SavePage saves page and all elements and element groups in page. Returns page id on success.
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

//DeletePage delete page and all page elements and page groups related
func DeletePage(id int) (int, error) {
	var pageID int
	// var groupID int
	err := db.QueryRow(`DELETE from pages
						WHERE id = $1
						RETURNING id`, id).Scan(&pageID)
	if err != nil {
		log.Fatal(err)
		return pageID, err
	}

	// db.QueryRow(`DELETE from elements WHERE pageid = $1`, id)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return pageID, err
	// }

	// err = db.QueryRow(`DELETE from elementgroups WHERE pageid = $1`, id).Scan(&groupID)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return pageID, err
	// }

	return pageID, nil
}
