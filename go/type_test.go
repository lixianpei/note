package type_test

import (
	"fmt"
	"testing"
)

// Map 是引用类型，函数或方法间传递，在函数或方法中进行新增或修改会影响原始数据
func TestMapEdit(t *testing.T) {
	data := make(map[string]any, 1)
	data["name"] = "LiXianPei"
	setMap(data, "age", 19)
	setMap(data, "phone", "123456")
	setMap(data, "summary", "这里是简介信息...")

	fmt.Println(data)
	//Output:
	//map[age:19 name:LiXianPei phone:123456 summary:这里是简介信息...]
}

func setMap(c map[string]any, key string, data any) {
	c[key] = data
}

// 切片作为函数的参数传递是 按值传递，切片内的数据是源数组的地址，因此当切片容量没有发生扩容时修改数据会影响源数据
// 切片在函数内使用append增加元素不会影响原切片，若发生cap扩容则切片地址发生变化，会生成一个新的切片
func TestSliceEdit(t *testing.T) {
	data := make([]int, 0, 1) //cap=1 容量大小决定了切片是否会发生扩容
	data = append(data, 1)
	fmt.Printf("元素地址：%p\n", data)

	editSlice(data)
	fmt.Println("元素修改后：", data) //此时数据已经发生了变更：[10]
}

// 切片作为值传递，但是切片指向的数组还是原来数组的地址，因此数组发生改变后会影响原数组的数据
func editSlice(items []int) {
	//若此时 append 增加的元素让切片扩容后 items 的地址发生变化，则 items 是新的地址，下面的修改则不会影响原数据
	//若此时 append 增加的元素没有让切片扩容，则 items 的修改会影响原数据，但是 append 的新元素不会出现在原数据中，可理解为两个切片只是共享元素地址信息
	items = append(items, 3, 4)
	for i := 0; i < len(items); i++ {
		items[i] = items[i] * 10
	}
	fmt.Printf("editSlice:%v,元素地址: %p\n", items, items)
}

func TestOver(t *testing.T) {
	var x uint8 = 255 // uint8范围[0-255]
	x = x + 1         // 溢出
	fmt.Println(x)    // 输出 0

	var y int8 = 127 //int8-范围[-128,127]
	y = y + 1        //溢出
	fmt.Println(y)   //output: -128
}
