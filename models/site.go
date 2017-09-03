package models

import "log"

//Site needs comment
type Site struct {
	ID       int     `json:"id"`
	Domain   string  `json:"domain"`
	UserID   int     `json:"userId"`
	DateTime string  `json:"dateTime"`
	Pages    []*Page `json:"pages"`
}

// GetSiteByUserID returns one site based on user id
// gets all pages and and page elements
func GetSiteByUserID(userID int) (*Site, error) {
	s := new(Site)
	pgs := make([]*Page, 0)

	rows, err := db.Query(` SELECT * 
							FROM sites INNER JOIN pages on sites.id = pages.siteid
						   	WHERE userid = $1`, userID)

	if err != nil {
		return s, err
	}
	defer rows.Close()

	for rows.Next() {

		p := new(Page)

		err := rows.Scan(
			&s.ID, &s.Domain, &s.UserID, &s.DateTime,
			&p.ID, &p.Title, &p.Path, &p.ParentID, &p.Name, &p.SiteID, &p.ShowInNav, &p.SortOrder, &p.Template)
		if err != nil {
			log.Fatal(err)
		}

		els, err := GetElementsByPageID(p.ID)
		if err != nil {
			log.Fatal(err)
		}

		grps, err := GetGroupsByPageID(p.ID)
		if err != nil {
			log.Fatal(err)
		}

		p.Elements = els
		p.Groups = grps
		pgs = append(pgs, p)

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return s, err
	}
	s.Pages = pgs

	return s, nil
}
