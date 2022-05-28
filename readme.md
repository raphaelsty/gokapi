# Gokapi

Gokapi implements the [Okapi BM25](https://www.wikiwand.com/en/Okapi_BM25) retriever in Go.

## Install

```sh
go get github.com/raphaelsty/gokapi@0.0.2
```

Gokapi is suitable when you want to:

- Implement the search engine on a single machine.
- Store the data on disk and thus be able to search among millions of documents without memory constraints.
- Update the retriever with new documents.
- Avoid the constraints of Elasticsearch and the JVM.

Here are some comparisons between Gokapi and the retrievers of the tool [Cherche](https://github.com/raphaelsty/cherche)

|       Retriever       | batch | disk storage | stand-alone |
|:---------------------:|:-----:|:------------:|:----------:|
|      Gokapi BM25      |   ✅   |       ✅      |      ✅     |
| Cherche Elasticsearch |   ✅   |       ✅      |      ❌     |
|     Cherche TF-IDF    |   ❌   |       ❌      |      ✅     |
|      Cherche BM25     |   ❌   |       ❌      |      ✅     |
|     Cherche Lunar     |   ❌   |       ❌      |      ✅     |

Gokapi stores the data necessary for the search on the disk (frequency of terms and metadata on the corpus). The reading (query) and writing (document indexing) times are higher than a retriever which stores data in memory, but Gokapi allows to process more documents without overloading the memory. Writing and reading on the disk are done with the [Diskv](https://github.com/peterbourgon/diskv) library. Gokapi does not store the content of the documents on the disk.

## Quick start

```go
package main

import (
	"fmt"

	"github.com/raphaelsty/gokapi"
)

func main() {

	data := make(map[string]string)

	data["document_0"] = "Paris is the capital of France"
	data["document_1"] = "Montreal is the capital of Canada"
	data["document_2"] = "Madrid is the capital of Spain"
	data["document_3"] = "Rome is the capital of Italy"

	retriever := gokapi.BM25("index")

	// Add the documents to the retriever.
	retriever.Add(data)

	// Top five answers.
	answers := retriever.Query("Paris France Canada", 5)

	for _, answer := range answers {
		fmt.Println(answer)
	}

	// Delete the index.
	retriever.Reset()

}
```

```
{document_0 2.002704}
{document_1 1.001352}
```

## Work in progress

Gokapi is under construction and may change soon.

Here are some "short-term" goals:

1. To be called by the Cherche library via Python to provide a lighter alternative to Elasticsearch.

2. Provide a command-line client to search for documents locally on your machine with the terminal at lightning speed.

3. Gokapi code needs to be enhanced.

4. Provide benchmarks.

