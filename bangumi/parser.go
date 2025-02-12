package bangumi

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bangumilite/bangumilite-component/utils"
	"regexp"
	"strconv"
	"strings"
)

var BackgroundImageUrlRegex = regexp.MustCompile(`url\(["']?(//[^\)"']+)["']?\)`)

// ParseSubjectIDs extracts unique subject IDs from the given document.
func ParseSubjectIDs(doc *goquery.Document) []int {
	var ids []int

	doc.Find("ul#browserItemList li a[href^='/subject/']").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		id, err := GetID(href)
		if err != nil {
			return
		}

		ids = append(ids, *id)
	})

	uniqueIDs := utils.RemoveDuplicates(ids)

	return uniqueIDs
}

// ParseImageURLFromSrc converts a relative image URL starting with '//' to an absolute URL
func ParseImageURLFromSrc(src string) *string {
	if len(src) == 0 {
		return nil
	}

	var image string
	if strings.HasPrefix(src, "//") {
		image = "https:" + src
	}

	return &image
}

// ParseImageURLFromStyle uses regex to extract the image src from background url style.
func ParseImageURLFromStyle(src string) *string {
	if len(src) == 0 {
		return nil
	}

	res := BackgroundImageUrlRegex.FindStringSubmatch(src)
	if len(res) <= 1 {
		return nil
	}

	image := res[1]
	if strings.HasPrefix(image, "//") {
		image = "https:" + image
	}

	// fallback to original image url
	return &image
}

// ParseSubjectType extracts the subject type id from span.ico_subject_type class and returns subject type in int
func ParseSubjectType(s string) (*int, error) {
	if len(s) == 0 {
		return nil, errors.New("subject type class is empty")
	}

	parts := strings.Split(s, " ")

	for _, part := range parts {
		if strings.HasPrefix(part, "subject_type_") {
			subjectType, err := strconv.Atoi(strings.TrimPrefix(part, "subject_type_"))

			if err != nil {
				return nil, err
			}

			return &subjectType, nil
		}
	}

	return nil, errors.New("no match subject type pattern has found")
}

// GetID extracts an ID from a URL-like string (href) in the form of "/<type>/127791".
func GetID(href string) (*int, error) {
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
