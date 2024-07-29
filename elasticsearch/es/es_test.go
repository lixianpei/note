package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"reflect"
	"testing"
)

const IndexName = "user"

type Person struct {
	Id      string   `json:"id"` //如果ID不是字符串，则生成的文档中_id会出现错误
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Married bool     `json:"married"`
	Money   float64  `json:"money"`
	Tags    []string `json:"tags"`
}

// Create a new index.
// number_of_shards 主分片
// number_of_replicas 副分片 单节点时分片数量只能设置为0，即使设置大于0也无法创建分片副本
const IndexMapping = `
{
    "settings" : {
        "index" : {
            "number_of_shards" : 1,
            "number_of_replicas" : 0 
        }
    },
	"mappings":{
		"properties":{
			"id":{
				"type":"keyword"
			},
			"name":{
				"type":"text",
				"fields": {
					"keyword": {
						"type": "keyword"
					}
				}
			},
			"age":{
				"type":"integer"
			},
			"married":{
				"type":"boolean"
			},
			"money":{
				"type":"double"
			},
			"tags":{
				"type":"text"
			}
		}
	}
}
`

// 创建索引
func TestEsCreateIndex(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		return
	}

	isExist, err := checkIndexIsExist(IndexName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if isExist {
		fmt.Println("索引已存在")
		return
	}

	res, err := client.CreateIndex(IndexName).BodyJson(IndexMapping).Do(context.TODO())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

// 检测索引是否存在
func checkIndexIsExist(indexName string) (bool, error) {
	client, err := NewClient()
	if err != nil {
		return false, err
	}
	//检测索引是否存在
	isExist, err := client.IndexExists(indexName).Do(context.TODO())
	if err != nil {
		return false, err
	}
	return isExist, nil
}

// 插入数据
func TestEsInsertData(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		return
	}
	isExist, err := checkIndexIsExist(IndexName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !isExist {
		fmt.Println("索引不存在")
		return
	}

	items := make([]Person, 0, 10)
	items = append(items,
		Person{Id: "1", Name: "赵一", Age: 18, Married: false, Money: 100, Tags: []string{"跑步", "篮球", "羽毛球"}},
		Person{Id: "2", Name: "赵二", Age: 19, Married: true, Money: 200, Tags: []string{"跑步", "篮球"}},
		Person{Id: "3", Name: "赵三", Age: 20, Married: false, Money: 200, Tags: []string{"跑步", "羽毛球"}},
		Person{Id: "4", Name: "赵四", Age: 21, Married: true, Money: 200, Tags: []string{"羽毛球"}},
		Person{Id: "5", Name: "赵五", Age: 22, Married: true, Money: 200, Tags: []string{"羽毛球"}},
		Person{Id: "6", Name: "李六", Age: 22, Married: true, Money: 600, Tags: []string{"篮球"}},
		Person{Id: "7", Name: "李七", Age: 30, Married: true, Money: 700, Tags: []string{"跑步", "篮球", "羽毛球"}},
	)
	bulkRequests := client.Bulk()
	for _, item := range items {
		req := elastic.NewBulkIndexRequest().Index(IndexName).Id(string(item.Id)).Doc(item)
		//req := elastic.NewBulkIndexRequest().Index(IndexName).Id(string(item.Id)).Doc(item)
		bulkRequests.Add(req)
	}
	bulkRes, err := bulkRequests.Do(context.TODO())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bulkRes)
	for _, item := range bulkRes.Items {
		d := item["index"]
		if bulkRes.Errors {
			fmt.Println("数据插入错误：", d.Error.Reason)
		} else {
			fmt.Printf("操作成功的数据：%+v\n", d)
		}
	}

	//p1 := Person{Name: "张三", Age: 18, Married: false}
	//p1 := Person{Name: "李四", Age: 25, Married: true}
	//p1 := Person{Name: "赵五2", Age: 35, Married: false, Money: 20056.78}
	//res, err := client.Index().
	//	Index("user").
	//	Id(p1.Name).
	//	BodyJson(p1).
	//	Do(context.Background())
	//fmt.Println(err, res)
	//fmt.Println(res.Version)
	//fmt.Println(res.Id)
	//fmt.Println(res.Type)
	//fmt.Println(res.Status)
	//fmt.Println(res.Result)
}

func TestEsUpdate(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		return
	}
	id := "赵五2"
	//p1 := Person{Name: "赵五2", Age: 22, Money: 66.77, Married: false}
	//p1 := Person{Name: "赵五2", Married: true}//如果用结构体会把其他字段也会一起更新，因此建议用map

	//使用map则只会更新map中的数据
	p1 := map[string]interface{}{
		"married": true,
	}
	res, err := client.Update().Index("user").Id(id).Doc(p1).Do(context.Background())
	fmt.Println(err, res) //若id不存在，则更新失败：elastic: Error 404 (Not Found): [_doc][赵五3]: document missing [typ  e=document_missing_exception] <nil>
	fmt.Println(res.Version)
	fmt.Println(res.Id)
	fmt.Println(res.Type)
	fmt.Println(res.Status)
	fmt.Println(res.Result)
}

func TestEsDelete(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		return
	}
	//根据文档ID进行删除数据
	//id := "赵五1"
	//res, err := client.Delete().Index("user").Id(id).Do(context.Background())
	//fmt.Println(err, res)
	//fmt.Println(res.Id)
	//fmt.Println(res.Type)
	//fmt.Println(res.Status)
	//fmt.Println(res.Result)
	//删除成功：
	//<nil> &{user _doc 赵五1 2 deleted 0xc0001a62c0 29 4   0 false}
	//赵五1
	//_doc
	//0
	//deleted
	//--- PASS: TestEsDelete (0.01s)

	//重复删除：
	//elastic: Error 404 (Not Found) &{user _doc 赵五1 3 not_found   0xc000066c80 30 4 0 false}
	//	赵五1
	//	_doc
	//	0
	//	not_found
	//	--- PASS: TestEsDelete (0.02s)

	//根据查询条件删除文档
	////q := elastic.NewTermQuery("name", "赵一") //TODO: 查询条件存在问题，无法删除
	//q := elastic.NewTermQuery("age", 35) //根据age条件删除，eg：不能用这个合并多个条件查询
	////q2 := elastic.NewTermQuery("married", true) //删除成功
	////q := elastic.NewMatchAllQuery()
	//result, err := client.DeleteByQuery().Index("user").Query(q).Do(context.TODO())
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//fmt.Println(result)
	//fmt.Println(result.Total) //删除成功的总数
	////delete ok: &{map[] 9 <nil> false 3 0 0 3 1 0 0 {0 0}  0 -1   0 []}

	//删除整个索引
	res, err := client.DeleteIndex(IndexName).Do(context.TODO())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(res.Acknowledged)
}

func TestEsSearchById(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		return
	}

	res, err := client.Get().Index(IndexName).Id("1").Do(context.TODO())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
	fmt.Println(string(res.Source))
	p := new(Person)
	err = json.Unmarshal(res.Source, p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(p)
}

// query  :查询条件
// bool:聚合查询组合条件 相当于（）
// must:必须满足 相当于and =
// should:条件可以满足 相当于or
// must_not:条件不需要满足，相当于and not
// range:范围
// gt: 大于
// lt:小于
// gte:大等于
// lte:小等于
// filter:条件过滤
// term:不分词
func TestEsSearchByParams(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		return
	}

	ctx := context.TODO()

	//根据ID查询单条记录
	//res, err := client.Get().Index(IndexName).Id("6").Do(ctx)
	//if elastic.IsNotFound(err) == true {
	//	//res = nil
	//	//err = elastic: Error 404 (Not Found)
	//	fmt.Println("找不到文档记录", err) //elastic: Error 404 (Not Found)
	//	return
	//}
	//if err != nil {
	//	fmt.Println("查询出现错误", err.Error())
	//	return
	//}
	//fmt.Println(err, res)
	//fmt.Println(res.Found)          //true
	//fmt.Println(string(res.Source)) //文档记录JSON格式: {"id":"6","name":"李六","age":22,"married":true,  "money":600,"tags":["篮球"]}

	//根据多个ID的查询文档记录
	//q := elastic.NewIdsQuery().Ids("3", "4")
	//res, err := client.Search(IndexName).Query(q).Do(ctx)
	//printEmployee(res, err)

	//查询某个字段，NewTermQuery 用于除text类型外的字段精确查询，NewMatchQuery用于text文本类型查询（需要精确查询时name.keyword，且mapping必须有keyword类型）
	//q := elastic.NewTermQuery("id", "7") // NewTermQuery 字段精确查询，除text类型字段外（ text 字段类型使用了分析器（analyzer），该分析器会将文本分词（tokenize）成多个词条（tokens），这使得 Term 查询（NewTermQuery）通常不能直接匹配 text 字段的原始输入值）
	//q := elastic.NewTermQuery("age", "22") //age=22的文档记录 字段精确查询
	//q := elastic.NewMatchQuery("name", "李六") // 模糊查询： 由于name字段是text类型， 默认情况下，作为分析的一部分，Elasticsearch会更改文本字段的值。这使得查找文本字段值的精确匹配变得困难。将会把 李六 和 李七 展示出来，李六的得分更高，优先匹配
	//q := elastic.NewMatchQuery("name.keyword", "李六") // 字段精确查询 索引mapping必须存在keyword类型，否是无法精确匹配 : {"name":{"type":"text","fields":{"keyword":{"type":"keyword"}}}}
	//res, err := client.Search(IndexName).Query(q).Do(ctx)
	//fmt.Printf("Query took %d milliseconds\n", res.TookInMillis)
	//fmt.Printf("Found %d hits\n", res.TotalHits())
	//if res.Hits.TotalHits.Value > 0 {
	//	for _, hit := range res.Hits.Hits {
	//		fmt.Printf("Document ID: %s\n", hit.Id)
	//		fmt.Printf("Source: %s\n", hit.Source)
	//	}
	//}

	//q := elastic.NewQueryStringQuery("name:赵一") //会把所有相关的都查询出来，比如包括所有 赵 的记录
	//q := elastic.NewQueryStringQuery("name.keyword:赵一") //利用 .keyword 对字段进行精确查询，只查询出 赵一  的文档记录
	//res, err := client.Search(IndexName).Query(q).Do(ctx)
	//printEmployee(res, err)

	//条件动态查询 {"query":{"bool":{"filter":[{"range":{"age":{"from":"21","include_lower":true,"include_upper":true,"to":null}}},{"match":{"name.keyword":{"query":"赵四"}}}],"must":{"match":{"married":{"query":true}}}}}}
	//q := elastic.NewBoolQuery()                      //用此拼接多个查询条件
	//q.Must(elastic.NewMatchQuery("married", true))   //bool值匹配
	//q.Filter(elastic.NewRangeQuery("age").Gte("21")) //age>21
	////q.Must(elastic.NewQueryStringQuery("name.keyword:赵四"))
	////q.Filter(elastic.NewMatchQuery("name.keyword", "赵四")) //字段进行精确匹配时，需要使用 .keyword 子字段。
	//res, err := client.Search(IndexName).Query(q).Do(ctx)
	//printEmployee(res, err)

	//查询name中包含 “李” 字符串的
	//matchPhraseQuery := elastic.NewMatchPhraseQuery("name", "李")
	//res, err := client.Search(IndexName).Query(matchPhraseQuery).Do(ctx)
	//printEmployee(res, err)

	//分页查询
	//page := 2
	//pageSize := 3
	//offset := (page - 1) * pageSize
	//res, err := client.Search(IndexName).Size(pageSize).From(offset).Do(ctx)
	//printEmployee(res, err)

	//OR 查询
	q := elastic.NewBoolQuery()
	//必须包含 篮球 和 羽毛球 的文档记录，可以再包含其他标签
	//q.Must(
	//	elastic.NewMatchPhraseQuery("tags", "篮球"),
	//	elastic.NewMatchPhraseQuery("tags", "羽毛球"),
	//)
	//只要包含其中一个
	//q.Should(
	//	elastic.NewMatchPhraseQuery("tags", "篮球"),
	//	elastic.NewMatchPhraseQuery("tags", "羽毛球"),
	//)

	//TODO : 查询只包含 篮球 的文档记录

	//查询 22 <= age < 99 的文档记录
	//q.Must(
	//	elastic.NewRangeQuery("age").Gte("22"),
	//	elastic.NewRangeQuery("age").Lt("99"),
	//)

	// 查询 age != 22 的文档记录
	q.MustNot(
		elastic.NewTermQuery("age", "22"),
	)

	//排序
	//s1 := elastic.NewFieldSort("age").Desc()
	//s2 := elastic.NewFieldSort("money").Desc()
	//Query(1).SortBy(s1,s2)

	//查询文档
	res, err := client.Search(IndexName).Query(q).Do(ctx)
	printEmployee(res, err)
}

// 打印查询到的Employee
func printEmployee(res *elastic.SearchResult, err error) {
	if err != nil {
		print(err.Error())
		return
	}
	var typ Person
	for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
		t := item.(Person)
		fmt.Printf("%#v\n", t)
	}
}
