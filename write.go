package gokapi

import (
	"encoding/json"
	"log"
)

func (retriever Retriever) writeTF(key string, value map[string]map[string]float32) {

	jsonStr, err := json.Marshal(value)
	if err != nil {
		log.Println("Warning: unable to parse document TF ", key, " as json.")
	}

	err = retriever.diskTF.Write(key, jsonStr)
	if err != nil {
		log.Println("Warning: Unable to save document TF ", key, "on disk.")
	}

}

func (retriever Retriever) writeIDF(key string, value map[string]float32) {

	jsonStr, _ := json.Marshal(value)

	retriever.diskIDF.Write(key, jsonStr)

}

func (retriever Retriever) writeMeta(n float32, mean float32) {

	jsonStr, err := json.Marshal(n)
	if err != nil {
		log.Println("Warning: unable to convert the total number of documents as JSON.")
	}

	err = retriever.diskMeta.Write("n", jsonStr)
	if err != nil {
		log.Println("Warning: Unable to write the total number of documents on disk.")
	}

	jsonStr, err = json.Marshal(mean)
	if err != nil {
		log.Println("Warning: unable to convert the average length number of documents as JSON.")
	}

	err = retriever.diskMeta.Write("mean", jsonStr)
	if err != nil {
		log.Println("Warning: Unable to write the average length number of documents on disk.")
	}

}
