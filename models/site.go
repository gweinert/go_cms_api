package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gweinert/cms_scratch/services"
)

//Site needs comment
type Site struct {
	ID       int     `json:"id"`
	Domain   string  `json:"domain"`
	UserID   int     `json:"userId"`
	DateTime string  `json:"dateTime"`
	Pages    []*Page `json:"pages"`
}

func BuildStaticJsonAndUpload(sessionID string, bucketName string) (string, error) {

	user, err := GetUserFromSessionID(sessionID)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	site, err := GetSiteByUserID(user.ID)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	site, err = transformSite(site) // transform for better as static site

	//create json file and upload to google cloud
	b, err := json.Marshal(site)
	if err != nil {
		fmt.Println("json err:", err)
		return "", err
	}

	buf := bytes.NewReader(b)
	fileName := "site.json"

	fileURL, err := services.GoogleCloudUpload(buf, bucketName, fileName)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return fileURL, nil
}

// GetSiteByUserID returns one site based on user id
// gets all pages and and page elements
func GetSiteByUserID(userID int) (*Site, error) {
	s := new(Site)
	pgs := make([]*Page, 0)

	rows, err := db.Query(` SELECT sites.id, domain, userid, datetime, pages.id, title, path, parentid, name, siteid, showinnav, sortorder, template 
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

		grps, err = AddElementsToGroups(els, grps)

		// tgs, err := transformGroups(els, grps)
		// if err != nil {
		// 	return s, err
		// }

		p.Elements = els
		p.Groups = grps
		// p.ElementsGroups = tgs
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

func GetSiteFromPageId(pageID int) (*Site, error) {
	s := new(Site)

	err := db.QueryRow(`SELECT sites.id, domain, userid
						FROM sites INNER JOIN pages on pages.siteid = sites.id
						WHERE pages.id = $1`, pageID).Scan(&s.ID, &s.Domain, &s.UserID)
	if err != nil {
		fmt.Println("fail here")
		log.Fatal(err)
		return nil, err
	}

	return s, nil
}

// func addElementsToGroups(els []*Element, grps []*ElementGroup) ([]*ElementGroup, error) {
// 	for _, g := range grps {
// 		gels := make([]*Element, 0)
// 		for _, el := range els {

// 			if el.GroupID == g.ID {

// 				gels = append(gels, el)
// 			}

// 		}

// 		g.Elements = gels
// 	}

// 	return grps, nil
// }

// func transformGroups(els []*Element, grps []*ElementGroup) (map[int][][]*Element, error) {
// 	var slidesPerGroup int
// 	slideGroup := make([][]*Element, 0)
// 	elsGroups := make(map[int][][]*Element)

// 	groupEls := filterElements(els, func(el *Element) bool {
// 		return el.GroupID != 0
// 	})

// 	slides := make([]*Element, 0)
// 	for i, el := range groupEls {
// 		currentGroup := findByKey(grps, el.GroupID)
// 		var slideNum int

// 		for _, groupEl := range currentGroup.Structure {
// 			slidesPerGroup += groupEl.Amount
// 		}

// 		slideNum = i / slidesPerGroup
// 		nextSlide := ((i + 1) / slidesPerGroup) > slideNum
// 		slides = append(slides, el)

// 		if nextSlide {
// 			slideGroup = append(slideGroup, slides)
// 			slides = make([]*Element, 0)

// 		}

// 		slidesPerGroup = 0

// 		elsGroups[el.GroupID] = slideGroup

// 		if (i+1 < len(groupEls)) && groupEls[i+1].GroupID != el.GroupID {
// 			slideGroup = make([][]*Element, 0)
// 		}

// 	}
// 	return elsGroups, nil

// }

func transformSite(site *Site) (*Site, error) {

	for _, p := range site.Pages {
		var elementMap = make(map[string]string)
		for _, e := range p.Elements {
			switch e.Type {
			case "image":
				elementMap[e.Name] = e.ImageURL
				break
			case "link":
				elementMap[e.Name] = e.LinkPath
				elementMap[strings.Join([]string{e.Name, "Text"}, "")] = e.LinkText
				break
			default:
				elementMap[e.Name] = e.Body
			}
		}
		p.ElementMap = elementMap
	}

	return site, nil
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
