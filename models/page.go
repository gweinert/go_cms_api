package models

import (
	"fmt"
	"log"
)

//Page used for somethign?
type Page struct {
	ID        int             `json:"id"`
	Title     string          `json:"title"`
	Path      string          `json:"path"`
	ParentID  int             `json:"parentId"`
	Name      string          `json:"name"`
	SiteID    int             `json:"siteId"`
	ShowInNav bool            `json:"showInNav"`
	SortOrder int             `json:"sortOrder"`
	Template  string          `json:"template"`
	Elements  []*Element      `json:"elements"`
	Groups    []*ElementGroup `json:"groups"`
	// ElementsGroups [][]*Element    `json:"elementsGroups"`
	ElementMap map[string]string `json:"elementMap"`
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
						VALUES($1, $2, $3, $4, $5, $6, $7, $8) 
						RETURNING id, title, path, parentid, name, siteid, showinnav, sortorder, template`,
		&p.Title, &p.Path, &p.ParentID, &p.Name, &p.SiteID, &p.ShowInNav, &p.SortOrder, &p.Template,
	).Scan(&np.ID, &np.Title, &np.Path, &np.ParentID, &np.Name, &np.SiteID, &np.ShowInNav, &np.SortOrder, &np.Template)

	if err != nil {
		log.Fatal(err)
		return np, err
	}

	return np, nil
}

//SavePage saves page and all elements and element groups in page. Returns page id on success.
func SavePage(up *Page) (*Page, error) {
	// sup := new(Page)
	stmt, err := db.Prepare(`UPDATE pages 
							SET name = $1, path = $2, sortorder = $3, template = $4, showinnav = $5
							WHERE id = $6 
							RETURNING id;`)
	if err != nil {
		fmt.Println("fail at prepare")
		log.Fatal(err)
	}
	res, err := stmt.Exec(up.Name, up.Path, up.SortOrder, up.Template, up.ShowInNav, up.ID)
	if err != nil {
		fmt.Println("fail at Exec")
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		fmt.Println("fail at rows affe")
		log.Fatal(err)
	}

	// if len(up.Elements) > 0 {
	nels, err := CreateOrUpdateElementIfExists(up.Elements)
	if err != nil {
		fmt.Println("fail at save el")
		log.Fatal(err)
	}

	grps, err := AddElementsToGroups(nels, up.Groups)

	up.Elements = nels
	up.Groups = grps
	log.Printf("ID = %d, affected = %d\n", up.ID, rowCnt)
	return up, nil
}

//DeletePage delete page and all page elements and page groups related
func DeletePage(id int) (int, int, int, error) {
	var pageID int
	var sortOrder int
	var parentID int
	// get page about to delete. find sort order of that element and page id and update all sort orders after it

	err := db.QueryRow(`SELECT sortorder, parentid from pages 
						WHERE id = $1`, id).Scan(&sortOrder, &parentID)

	err = db.QueryRow(`UPDATE pages
						SET sortorder = sortorder - 1
						WHERE parentid = $1 and sortorder > $2
						RETURNING id`, parentID, sortOrder).Scan(&pageID)

	err = db.QueryRow(`DELETE from pages
						WHERE id = $1
						RETURNING id`, id).Scan(&pageID)
	if err != nil {
		log.Fatal(err)
		return 0, 0, 0, err
	}

	return pageID, sortOrder, parentID, nil
}

func UpdatePageSortOrder(id int, newIndex int) (int, error) {
	var oldIndex int
	var parentID int

	err := db.QueryRow(`SELECT sortorder, parentid from pages WHERE id = $1`, id).Scan(&oldIndex, &parentID)
	if err != nil {
		fmt.Println("fail select query page sort order")
		return 0, err
	}

	db.QueryRow(`UPDATE pages
				SET sortorder = $1
				WHERE sortorder = $2 AND parentid = $3`, oldIndex, newIndex, parentID)

	db.QueryRow(`UPDATE pages
				SET sortorder = $1
				WHERE id = $2 AND parentid = $3`, newIndex, id, parentID)

	return id, nil
}
