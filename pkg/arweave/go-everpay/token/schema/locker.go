package schema

import "fmt"

type HeightLocker struct {
	Address     string
	BeginHeight int64
	EndHeight   int64
}

type Locker struct {
	ChainType string
	ChainId   string
	List      []HeightLocker
}

type Lockers map[string]Locker // key: oracle chainType

func (a Lockers) GetLocker(chainType string) (Locker, error) {
	chainLocker, ok := a[chainType]
	if !ok {
		return Locker{}, fmt.Errorf("GetLocker not found the chainType HeightLocker; chainType:%s", chainType)
	}
	return chainLocker, nil
}

func (c Locker) GetLockerAddrFromHeight(height int64) (addr string) {
	for _, lk := range c.List {
		if lk.BeginHeight <= height && lk.EndHeight >= height {
			addr = lk.Address
			return
		}
	}
	return
}

func (c Locker) CurLockerAddr() string {
	lockerArr := c.List
	return lockerArr[len(lockerArr)-1].Address
}
