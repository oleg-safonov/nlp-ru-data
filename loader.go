package nlprudata

import (
	"archive/zip"
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/oleg-safonov/nlp"
)

//go:embed dict.zip
var dictBytes []byte

//go:embed suffix.zip
var suffixBytes []byte

func Load() (nlp.LemmatizerData, error) {
	base := nlp.LemmatizerData{}
	dictBase, err := loadDict()
	if err != nil {
		return base, err
	}
	suffixBase, err := loadSuffix()
	if err != nil {
		return base, err
	}
	base.Dictionary = *dictBase
	base.SuffixPredictor = *suffixBase

	return base, nil
}

func loadDict() (*nlp.DictionaryBase, error) {
	dataReader := bytes.NewReader(dictBytes)
	r, err := zip.NewReader(dataReader, int64(len(dictBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	var dictionary *nlp.DictionaryBase
	for _, f := range r.File {
		if f.Name == "dictionary.bin" {
			dictionary, err = loadFile[nlp.DictionaryBase](f)
			if err != nil {
				return nil, err
			}
		}
	}

	return dictionary, nil
}

func loadSuffix() (*nlp.SuffixPredictorBase, error) {
	dataReader := bytes.NewReader(suffixBytes)
	r, err := zip.NewReader(dataReader, int64(len(suffixBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	var suffixBase *nlp.SuffixPredictorBase
	for _, f := range r.File {
		if f.Name == "suffix.bin" {
			suffixBase, err = loadFile[nlp.SuffixPredictorBase](f)
			if err != nil {
				return nil, err
			}
		}
	}

	return suffixBase, nil
}

func loadFile[T any](f *zip.File) (*T, error) {
	file, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", f.Name, err)
	}

	var base T
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&base)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", f.Name, err)
	}
	return &base, nil
}
