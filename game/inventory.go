package game

import (
	"bufio"
	"encoding/csv"
)

const USE_TYPE_EQUIP = 0
const USE_TYPE_USE = 1
const USE_TYPE_QUEST = 2

type ItemOrigin struct {
	id      int
	name    string
	useType int
	desc    string
}

type ItemOriginManager struct {
	dict map[int]*ItemOrigin
}

func (i *ItemOriginManager) LoadFromCSV(filename string) {
	file := OpenFile(filename)
	defer file.Close()
	rdr := csv.NewReader(bufio.NewReader(file))
	rows, _ := rdr.ReadAll()
	for idx, row := range rows {
		if idx == 0 {
			continue
		}

		item := &ItemOrigin{id: atoi(row[0]),
			name:    row[1],
			useType: atoi(row[2]),
			desc:    row[3],
		}

		i.dict[item.id] = item
	}

	dump(i.dict)
}

type Item struct {
	originId int
}

type Inventory struct {
}
