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

		tgs, err := transformGroups(els, grps)
		if err != nil {
			return s, err
		}

		p.Elements = els
		p.Groups = grps
		p.ElementsGroups = tgs
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

func transformGroups(els []*Element, grps []*ElementGroup) (map[int][][]*Element, error) {
	var slidesPerGroup int
	slideGroup := make([][]*Element, 0)
	elsGroups := make(map[int][][]*Element)

	groupEls := filterElements(els, func(el *Element) bool {
		return el.GroupID != 0
	})

	slides := make([]*Element, 0)
	for i, el := range groupEls {
		currentGroup := findByKey(grps, el.GroupID)
		var slideNum int

		for _, groupEl := range currentGroup.Structure {
			slidesPerGroup += groupEl.Amount
		}

		slideNum = i / slidesPerGroup
		nextSlide := ((i + 1) / slidesPerGroup) > slideNum
		slides = append(slides, el)

		if nextSlide {
			slideGroup = append(slideGroup, slides)
			slides = make([]*Element, 0)

		}

		slidesPerGroup = 0

		elsGroups[el.GroupID] = slideGroup

		if (i+1 < len(groupEls)) && groupEls[i+1].GroupID != el.GroupID {
			slideGroup = make([][]*Element, 0)
		}

	}
	return elsGroups, nil

}

func findByKey(grps []*ElementGroup, groupID int) *ElementGroup {
	group := new(ElementGroup)
	for _, g := range grps {
		if g.ID == groupID {
			group = g
			break
		}
	}

	return group
}

func filterElements(els []*Element, f func(*Element) bool) []*Element {
	elsf := make([]*Element, 0)
	for _, el := range els {
		if f(el) {
			elsf = append(elsf, el)
		}
	}

	return elsf
}
