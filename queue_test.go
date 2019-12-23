package tinycache

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func Test_priorityQueue_pushItem(t *testing.T) {
	queue := newPriorityQueue()
	for i := 0; i <= 5; i++ {
		t.Run(fmt.Sprintf("item%d", i), func(t *testing.T) {
			itemData := fmt.Sprintf("itemData%d", i)
			queue.Push(&cacheItem{data: itemData})
			if queueLen := queue.Len(); queueLen != i+1 {
				t.Errorf("after queue.Push has wrong leng")
			}
		})
	}
}

func TestPriorityQueue_popItemAndOrder(t *testing.T) {
	queue := newPriorityQueue()
	for i := 0; i < 5; i++ {
		queue.Push(&cacheItem{data: fmt.Sprintf("data%d", i)})
	}

	for i := 4; i >= 0; i-- {
		item := queue.Pop()
		if queue.Len() != i {
			t.Errorf("wrong queue len")
		}

		cacheItem := item.(*cacheItem)
		wantData := fmt.Sprintf("data%d", i)
		if cacheItem.data != wantData {
			t.Errorf("want data:%v, ret:%v", wantData, cacheItem.data)
		}
	}
}

func Test_priorityQueue_updateAndRemoveItem(t *testing.T) {
	queue := newPriorityQueue()
	var chosenItem *cacheItem
	for i := 0; i < 5; i++ {
		item := newCacheItem(fmt.Sprintf("key%d", i), "data", time.Duration(i)*time.Second)
		if i == 1 {
			chosenItem = item
		}
		queue.pushItem(item)
	}

	//for _, item := range *queue {
	//	fmt.Printf("%v\n", item)
	//}

	chosenItem.update("data2", time.Duration(6)*time.Hour)
	queue.updateItem(chosenItem)
	sort.Sort(queue)

	t.Run("update item", func(t *testing.T) {
		if chosenItem.index != 4 {
			t.Errorf("updatedItem has index %d, expect:%d", chosenItem.index, 4)
		}
	})
	queue.removeItem(chosenItem)

	t.Run("remove item", func(t *testing.T) {
		if queue.Len() != 4 {
			t.Errorf("expected queue len: 4, ret:%d", queue.Len())
		}
	})

}
