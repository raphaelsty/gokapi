package gokapi

import (
	"log"
	"os"
)

// Delete index.
func (retriever Retriever) Reset() {

	errIDF := os.RemoveAll(retriever.path + "/idf")
	if errIDF != nil {
		log.Fatal("Unable to delete index: ", retriever.path+"/idf")
	}

	errTF := os.RemoveAll(retriever.path + "/tf")
	if errTF != nil {
		log.Fatal("Unable to delete index: ", retriever.path+"/tf")
	}

	errMeta := os.RemoveAll(retriever.path + "/meta")
	if errMeta != nil {
		log.Fatal("Unable to delete metadata: ", retriever.path+"/meta")
	}
}
