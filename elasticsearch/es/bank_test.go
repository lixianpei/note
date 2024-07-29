package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"testing"
)

const BankIndexName = "bank"

func TestBankSearchAggregations(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		return
	}
	ctx := context.TODO()

	//同时按多个字段多种聚合查询
	//{"aggregations":{"genders":{"aggregations":{"balances1":{"avg":{"field":"balance","format":".00"}}},"terms":{"field":"gender.keyword","order":[{"balances1":"asc"}]}},"states":{"aggregations":{"balances":{"avg":{"field":"balance","format":".00"}}},"terms":{"field":"state.keyword","order":[{"balances":"asc"}]}}},"size":0}
	//按州维度-聚合1
	agg1 := elastic.NewTermsAggregation().Field("state.keyword")             //统计各州人数
	avgBalance := elastic.NewAvgAggregation().Field("balance").Format(".00") //各州的平均帐户余额是多少？
	agg1.SubAggregation("balances", avgBalance)
	//agg1.OrderByCount(true)  //按当前聚合的key进行排序-升序，即state.count
	//agg1.OrderByCount(false) //按当前聚合的key进行排序-降序，即state.count
	//按性别维度-聚合2
	agg2 := elastic.NewTermsAggregation().Field("gender.keyword")             //按性别统计人数
	avgBalance2 := elastic.NewAvgAggregation().Field("balance").Format(".00") //性别维度的平均帐户余额是多少？
	agg2.SubAggregation("balances1", avgBalance2)
	agg2.OrderByAggregation("balances1", true) //按子查询balance字段进行升序
	agg2.Missing("10")                         //参数定义应如何处理缺少值的文档，指定缺失的文档安装10进行统计

	//普通字段聚合-统计最大年龄
	agg3 := elastic.NewMaxAggregation().Field("age")

	//桶聚合bucket 可以理解为一个桶，它会遍历文档中的内容，凡是符合某一要求的就放入一个桶中，分桶相当于 SQL 中的 group by。从另外一个角度，可以将指标聚合看成单桶聚合，即把所有文档放到一个桶中，而桶聚合是多桶型聚合，它根据相应的条件进行分组。
	//统计各个性别中的余额最大值
	agg4 := elastic.NewTermsAggregation().Field("gender.keyword").
		SubAggregation("max_balance", elastic.NewMaxAggregation().Field("balance").Format(".00"))

	res, err := client.Search(BankIndexName).
		Aggregation("count_states", agg1).
		Aggregation("count_genders", agg2).
		Aggregation("max_age", agg3).
		Aggregation("max_gender_balance", agg4).
		Size(0).
		Do(ctx)

	fmt.Println("==================按各州聚合==========================")
	s1, found := res.Aggregations.Terms("count_states")
	if found {
		for _, bucket := range s1.Buckets {
			b := bucket.Aggregations["balances"]
			fmt.Printf("Force: %s, Count: %d ,balances: %s \n", bucket.Key, bucket.DocCount, b)
		}
	}

	fmt.Println("=======================按性别聚合=======================")
	s2, found := res.Aggregations.Terms("count_genders")
	if found {
		for _, bucket := range s2.Buckets {
			b := bucket.Aggregations["balances1"]
			fmt.Printf("Force: %s, Count: %d ,balances: %s \n", bucket.Key, bucket.DocCount, b)
		}
	}

	fmt.Println("=======================按最大年龄聚合=======================")
	s3, found := res.Aggregations.Max("max_age")
	if found {
		fmt.Printf("最大年龄：%+v, 其他信息：%+v \n", *s3.Value, s3)
	}

	fmt.Println("=======================统计各个性别中的余额最大值=======================")
	s4, found := res.Aggregations.Terms("max_gender_balance")
	if found {
		for _, bucket := range s4.Buckets {
			b := bucket.Aggregations["max_balance"]
			fmt.Printf("Force: %s, Count: %d ,max_balance: %s \n", bucket.Key, bucket.DocCount, b)
		}
	}
}

func printSearchResult(res *elastic.SearchResult, err error) {
	fmt.Println(err, res)
	if err != nil {
		return
	}
	fmt.Printf("Query took %d milliseconds\n", res.TookInMillis)
	fmt.Printf("Found %d hits\n", res.TotalHits())
	if res.Hits.TotalHits.Value > 0 {
		for _, hit := range res.Hits.Hits {
			fmt.Printf("Document ID: %s\n", hit.Id)
			fmt.Printf("Source: %s\n", hit.Source)
		}
	}
}
