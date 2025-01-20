package bangumi

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sstp105/bangumi-component/utils"
	"strconv"
	"strings"
)

// ParseSubjectIDs extracts unique subject IDs from the given document.
func ParseSubjectIDs(doc *goquery.Document) []int {
	var ids []int

	doc.Find("ul#browserItemList li a[href^='/subject/']").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		id, err := getID(href)
		if err != nil {
			return
		}

		ids = append(ids, *id)
	})

	uniqueIDs := utils.RemoveDuplicates(ids)

	return uniqueIDs
}

// getID extracts an ID from a URL-like string (href) in the form of "/<type>/127791".
func getID(href string) (*int, error) {
	if len(href) == 0 {
		return nil, errors.New("href is empty, unable to parse the id")
	}

	parts := strings.Split(href, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("href is invalid, split by / returns %s", parts)
	}

	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("error converting id %s to int", idStr)
	}

	return &id, nil
}
