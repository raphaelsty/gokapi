package gokapi

import (
	"regexp"
	"strings"
	"sync"
)

type docMeta struct {
	tf  map[string]map[string]map[string]float32
	len float32
}

func stripRegex(in string) string {
	re, _ := regexp.Compile(`[^\w]`)
	return re.ReplaceAllString(in, " ")
}

func clean(input string) string {
	return stripRegex(strings.ToLower(input))
}

func firstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func process(id string, document string) docMeta {

	document = clean(document)

	content := make(map[string]map[string]map[string]float32)

	tokens := strings.Split(document, " ")

	len := float32(len(tokens))

	for _, token := range tokens {

		key := firstN(token, 3)

		f := float32(strings.Count(document, token))

		if _, exist := content[key]; !exist {
			content[key] = make(map[string]map[string]float32)
		}

		content[key][token] = make(map[string]float32)
		content[key][token][id] = f

	}

	return docMeta{content, len}
}

func (retriever Retriever) Add(documents map[string]string) {

	var wg sync.WaitGroup

	tf := make(map[string]map[string]map[string]float32)
	idf := make(map[string]map[string]float32)
	queue := make(chan docMeta, len(documents))

	for id, document := range documents {

		wg.Add(1)
		go func(id string, document string) {
			defer wg.Done()
			queue <- process(id, document)
		}(id, document)
	}

	wg.Wait()

	close(queue)

	mean := retriever.Mean()
	n := retriever.Size()

	for sample := range queue {

		n++
		mean += (sample.len - mean) / n

		for key, tokenIDTF := range sample.tf {

			if _, ok := tf[key]; !ok {
				tf[key] = make(map[string]map[string]float32)
				idf[key] = make(map[string]float32)
			}

			for token, idtf := range tokenIDTF {

				// Unknown token
				if _, ok := tf[key][token]; !ok {
					tf[key][token] = make(map[string]float32)
					idf[key][token] = retriever.IDF(token) + 1
				} else {
					// Known token.
					idf[key][token]++
				}

				for id, value := range idtf {
					tf[key][token][id] = value
				}
			}
		}
	}

	retriever.writeMeta(n, mean)

	var wg2 sync.WaitGroup

	for key, value := range tf {

		wg2.Add(1)

		go func(retriever *Retriever, key string, value map[string]map[string]float32) {
			defer wg2.Done()
			retriever.writeTF(key, value)
		}(&retriever, key, value)
	}

	for key, value := range idf {

		wg2.Add(1)

		go func(retriever *Retriever, key string, value map[string]float32) {
			defer wg2.Done()
			retriever.writeIDF(key, value)
		}(&retriever, key, value)
	}

	wg2.Wait()

}
