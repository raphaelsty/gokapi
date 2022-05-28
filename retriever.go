package github.com/raphaelsty/gokapi

import (
	"encoding/json"
	"math"
	"path/filepath"
	"sort"
	"strings"

	"github.com/peterbourgon/diskv"
)

type Document struct {
	id    string
	score float32
}

type Retriever struct {
	path     string
	diskTF   diskv.Diskv
	diskIDF  diskv.Diskv
	diskMeta diskv.Diskv
	k1       float32
	b        float32
}

func BM25(path string) Retriever {

	var n float32
	var mean float32

	diskTF := diskv.New(diskv.Options{
		BasePath:     filepath.FromSlash(path + "/tf"),
		CacheSizeMax: 1024 * 1024,
	})

	diskIDF := diskv.New(diskv.Options{
		BasePath:     filepath.FromSlash(path + "/idf"),
		CacheSizeMax: 1024 * 1024,
	})

	diskMeta := diskv.New(diskv.Options{
		BasePath:     filepath.FromSlash(path + "/meta"),
		CacheSizeMax: 1024 * 1024,
	})

	value, _ := diskMeta.Read("n")
	err := json.Unmarshal(value, &n)
	if err != nil {
		n = 0
	}

	value, _ = diskMeta.Read("mean")
	err = json.Unmarshal(value, &mean)
	if err != nil {
		mean = 0
	}

	retriever := Retriever{path, *diskTF, *diskIDF, *diskMeta, 1.6, 0.75}
	return retriever
}

func (retriever Retriever) Query(q string, k int) []Document {

	q = clean(q)

	scores := make(map[string]float32)

	tokens := strings.Split(q, " ")

	n := retriever.Size()
	mean := retriever.Mean()

	for _, token := range tokens {

		docsTF := retriever.TF(token)

		idf := retriever.IDF(token)
		idf = float32(math.Log(float64((n - idf + 0.5) / (idf + 0.5))))

		for id, tf := range docsTF {

			tf := (tf * (retriever.k1 + 1)) / (tf + retriever.k1*(1-retriever.b+retriever.b*(n/mean)))

			if _, exist := scores[id]; !exist {
				scores[id] = tf * idf
			} else {
				scores[id] += tf * idf
			}
		}
	}

	match := make([]string, 0, len(scores))
	for id := range scores {
		match = append(match, id)
	}

	sort.Slice(match, func(i, j int) bool {
		return scores[match[i]] > scores[match[j]]
	})

	var results []Document

	for top, id := range match {

		if top >= k {
			break
		}
		results = append(results, Document{id: id, score: scores[id]})
	}

	return results
}
