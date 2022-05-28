package github.com/raphaelsty/gokapi

import (
	"encoding/json"
	"log"
)

func (retriever Retriever) TF(token string) map[string]float32 {

	var value []byte

	key := firstN(token, 3)

	value, errRead := retriever.diskTF.Read(key)

	match := make(map[string]map[string]float32)
	err := json.Unmarshal(value, &match)
	if (err != nil) && (errRead == nil) {
		log.Println("Warning: parse ", token, " TF from disk.")
	}

	return match[token]
}

func (retriever Retriever) IDF(token string) float32 {

	var value []byte

	key := firstN(token, 3)

	value, errRead := retriever.diskIDF.Read(key)

	match := make(map[string]float32)
	err := json.Unmarshal(value, &match)
	if (err != nil) && (errRead == nil) {
		log.Println("Warning: parse ", token, " IDF from disk.")
	}

	return match[token]
}

func (retriever Retriever) Mean() (mean float32) {

	var value []byte

	value, errRead := retriever.diskMeta.Read("mean")

	err := json.Unmarshal(value, &mean)
	if (err != nil) && (errRead == nil) {
		log.Println("Unable to read the mean length of documents.")
	}

	return mean
}

func (retriever Retriever) Size() (n float32) {

	var value []byte

	value, errRead := retriever.diskMeta.Read("n")

	err := json.Unmarshal(value, &n)
	if (err != nil) && (errRead == nil) {
		log.Println("Unable to read the number of documents in the corpora.")
	}

	return n
}
