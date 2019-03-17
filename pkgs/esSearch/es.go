package esSearch

import (
	"context"
	"fmt"
	log "github.com/AlexStocks/log4go"
	"github.com/olivere/elastic"
)

const (
	esCluserUrl = "http://119.23.211.243:9200"
)

type User struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Hobby string `json:"hobby"`
}

func NewEsClient() *elastic.Client {
	esCli, err := elastic.NewClient(elastic.SetURL(esCluserUrl), elastic.SetSniff(false))
	if err != nil {
		log.Error("===:", err)
		fmt.Println(err)
		return nil
	}
	return esCli
}

func IsEsCluserOk() bool {
	esCli := NewEsClient()

	ctx := context.Background()

	info, num, err := esCli.Ping(esCluserUrl).Do(ctx)
	if err != nil {
		log.Error(err)
		return false
	}

	fmt.Printf("es info:%+v\n", info)
	fmt.Printf("es num:%+v\n", num)
	return true
}

func UseEs() error {

	esCli := NewEsClient()

	createRes, err := esCli.CreateIndex("usertest").Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("createRes:", createRes)
	return nil
}

func AddDocmToIndex() error {
	esCli := NewEsClient()
	someone := User{
		Name:  "xiaohong123",
		Id:    "100002",
		Hobby: "cook",
	}
	res, err := esCli.Index().
		Index("usertest").
		Type("doc").
		Id("4").
		BodyJson(someone).
		Refresh("wait_for").
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("add doc:", res)
	return nil
}

func SearchEs() error {
	esCli := NewEsClient()

	termQuery := elastic.NewPrefixQuery("name", "xiao")
	searchResult, err := esCli.Search().
		Index("usertest").          // search in index "tweets"
		Query(termQuery).           // specify the query
		Sort("name.keyword", true). // sort by "user" field, ascending
		From(0).Size(10).           // take documents 0-9
		Pretty(true).               // pretty print request and response JSON
		Do(context.Background())    // execute
	if err != nil {
		return err
	}

	fmt.Printf("Query took: %+v\n", searchResult.Hits.TotalHits)
	for _, v := range searchResult.Hits.Hits {
		fmt.Printf("Query took: %+v\n", *v)
	}

	return nil
}

func readLog() {
	esCli := NewEsClient()
	a := elastic.MatchAllQuery{}
	ret := esCli.Search("index").SearchType("log").Query(a)
	fmt.Printf("%+v", ret)
}
