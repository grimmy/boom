package boom

import (
	"fmt"
	"time"

	. "gopkg.in/check.v1"
)

type ColSuite struct{}

var _ = Suite(&ColSuite{})

func (s *ColSuite) TestWaitTimeout(c *C) {
	col := NewCollector()

	col.Run(func(data ...interface{}) *TaskResult {
		time.Sleep(15 * time.Millisecond)
		return NewResult(1, nil)
	})
	col.Run(func(data ...interface{}) *TaskResult {
		time.Sleep(20 * time.Millisecond)
		return NewResult(2, nil)
	})
	col.Run(func(data ...interface{}) *TaskResult {
		time.Sleep(30 * time.Millisecond)
		return NewResult(3, nil)
	})

	res, err := col.Wait(10 * time.Millisecond)
	c.Check(res, IsNil)
	c.Check(err, NotNil)
	c.Check(err, Equals, ErrTimeout)
}

func (s *ColSuite) TestWaitCollect(c *C) {
	col := NewCollector()

	col.Run(func(data ...interface{}) *TaskResult {
		time.Sleep(5 * time.Millisecond)
		return NewResult(1, nil)
	})
	col.Run(func(data ...interface{}) *TaskResult {
		time.Sleep(1 * time.Millisecond)
		return NewResult(2, nil)
	})
	col.Run(func(data ...interface{}) *TaskResult {
		return NewResult(3, nil)
	})

	res, err := col.Wait(10 * time.Millisecond)
	c.Check(err, IsNil)
	c.Check(res, NotNil)
	c.Check(res[0], DeepEquals, &TaskResult{Value: 1, Err: nil})
	c.Check(res[1], DeepEquals, &TaskResult{Value: 2, Err: nil})
	c.Check(res[2], DeepEquals, &TaskResult{Value: 3, Err: nil})
}

func (s *ColSuite) TestArgs(c *C) {
	col := NewCollector()

	col.Run(func(data ...interface{}) *TaskResult {
		return NewResult(fmt.Sprintf("r %d %d", data[0], data[1]), nil)
	}, 1, 2)
	col.Run(func(data ...interface{}) *TaskResult {
		return NewResult(fmt.Sprintf("r %d %d", data[0], data[1]), nil)
	}, 3, 4)
	col.Run(func(data ...interface{}) *TaskResult {
		return NewResult(fmt.Sprintf("r %d %d", data[0], data[1]), nil)
	}, 5, 6)

	res, err := col.Wait(10 * time.Millisecond)
	c.Check(err, IsNil)
	c.Check(res, NotNil)
	c.Check(res[0], DeepEquals, &TaskResult{Value: "r 1 2", Err: nil})
	c.Check(res[1], DeepEquals, &TaskResult{Value: "r 3 4", Err: nil})
	c.Check(res[2], DeepEquals, &TaskResult{Value: "r 5 6", Err: nil})
}