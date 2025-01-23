package bangumi

import (
	"errors"
	"strings"
)

var ErrCategoryMismatch = errors.New("image categories do not match")
var ErrInvalidPattern = errors.New("image does not have matched pattern")

type ImagePath string

const (
	MonoImagePattern    string = "//lain.bgm.tv/pic/crt/"
	SubjectImagePattern string = "//lain.bgm.tv/pic/cover/"

	MonoImage    string = "mono"
	SubjectImage string = "subject"

	CharacterLarge  ImagePath = "//lain.bgm.tv/pic/crt/l/"
	CharacterMedium ImagePath = "//lain.bgm.tv/pic/crt/m/"
	CharacterGrid   ImagePath = "//lain.bgm.tv/pic/crt/g/"
	CharacterSmall  ImagePath = "//lain.bgm.tv/pic/crt/s/"

	SubjectLarge  ImagePath = "//lain.bgm.tv/pic/cover/l/"
	SubjectMedium ImagePath = "//lain.bgm.tv/pic/cover/m/"
	SubjectGrid   ImagePath = "//lain.bgm.tv/pic/cover/g/"
	SubjectSmall  ImagePath = "//lain.bgm.tv/pic/cover/s/"
)

func (imgPath ImagePath) Category() string {
	if strings.Contains(string(imgPath), MonoImagePattern) {
		return MonoImage
	}

	if strings.Contains(string(imgPath), SubjectImagePattern) {
		return SubjectImage
	}

	return ""
}

func Convert(from ImagePath, to ImagePath, src string) (*string, error) {
	if from.Category() != to.Category() {
		return nil, ErrCategoryMismatch
	}

	if !strings.Contains(string(src), string(from)) {
		return nil, ErrInvalidPattern
	}

	imageURL := strings.Replace(src, string(from), string(to), 1)
	return &imageURL, nil
}
