package pkg

import "fmt"

type Content struct {
	Ext *ContentExt
}
type ContentExt struct {
	SimClusterID *string
}
type Answer struct {
	SimClusterID *string
}

func (a *Answer) GetSimClusterID() *string {
	return a.SimClusterID
}

func PointerTest() {
	i := 10
	fmt.Println(&i)

	contentAnswers := make(map[string]*Content)
	house1 := "aaa"
	house2 := "bbb"
	//house3 := "ccc"
	contentAnswers["1"] = &Content{
		Ext: &ContentExt{
			SimClusterID: &house1,
		},
	}
	contentAnswers["2"] = &Content{
		Ext: &ContentExt{
			SimClusterID: &house2,
		},
	}
	for _, contentAnswer := range contentAnswers {
		fmt.Printf("Answer：%+v\n", *contentAnswer.Ext.SimClusterID)
	}

}
