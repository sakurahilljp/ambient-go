package ambient

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSend(t *testing.T) {

	c := NewClient(4460, "9c1d1656cfff47a7")
	{
		dp := NewDataPoint()
		dp["d1"] = rand.Float64()*2 + 20
		err := c.Send(dp)
		fmt.Println("err = ", err)
	}

	{
		t := time.Now()
		dp1 := NewDataPoint(t)
		dp1["d1"] = rand.Float64()*2 + 20

		dp2 := NewDataPoint(t.Add(time.Second * 2))
		dp2["d1"] = rand.Float64()*2 + 20

		dp3 := NewDataPoint(t.Add(time.Second * 4))
		dp3["d1"] = rand.Float64()*2 + 20

		// Wait to avoid limitation
		time.Sleep(time.Second * 6)

		err := c.Send(dp1, dp2, dp3)
		fmt.Println("err = ", err)
	}
}

func TestRead(t *testing.T) {
	c := NewClient(4460, "9c1d1656cfff47a7", ReadKey("6a0158c51ef97701"))
	{
		values, err := c.Read(Count(5), Skip(5))

		fmt.Println("Read with Count and Skip")
		fmt.Println("err = ", err)
		fmt.Println("values = ", values)
		fmt.Println("len(values) = ", len(values))
	}

	{
		values, err := c.Read(Date(time.Now().Add(-time.Hour * 24)))

		fmt.Println("Read with date")
		fmt.Println("err = ", err)
		fmt.Println("values = ", values)
		fmt.Println("len(values) = ", len(values))
	}

	{
		values, err := c.Read(Range(time.Now().Add(-time.Second*30), time.Now().Add(time.Second*30)))

		fmt.Println("Read with Range")
		fmt.Println("err = ", err)
		fmt.Println("values = ", values)
		fmt.Println("len(values) = ", len(values))
	}
}

func TestGetProp(t *testing.T) {
	c := NewClient(4460, "9c1d1656cfff47a7", ReadKey("6a0158c51ef97701"))
	prop, err := c.GetProp()
	fmt.Println("err = ", err)
	fmt.Println("prop = ", prop)

}
