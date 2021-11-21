package database

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Database struct {
	raw map[string]map[uuid.UUID]RawValue
	dic map[string]map[string]map[uuid.UUID]struct{}
}

var dbInstance *Database

func DatabaseInstance() *Database {
	if dbInstance == nil {
		dbInstance = newDb()
	}
	return dbInstance
}

func (d *Database) GetRawValues(guild string, tags []string) []RawValue {
	result := []RawValue{}
	if _, ok := d.dic[guild]; !ok {
		return result
	}
	gDic := d.dic[guild]
	toCheckSet := []map[uuid.UUID]struct{}{}
	for _, tag := range tags {
		if _, ok := gDic[tag]; ok {
			toCheckSet = append(toCheckSet, gDic[tag])
		}
	}
	if len(toCheckSet) == 0 {
		return result
	}
	subset := toCheckSet[0]
	for i := 1; i < len(toCheckSet); i++ {
		subset = intersection(subset, toCheckSet[i])
	}
	if len(subset) == 0 {
		return result
	}
	for id := range subset {
		if content, ok := d.raw[guild][id]; ok {
			result = append(result, content)
		}
	}
	return result
}

func (d *Database) GetRawValue(guild string, tags []string) *RawValue {
	values := d.GetRawValues(guild, tags)
	if len(values) == 0 {
		return nil
	}
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	return &values[rand.Intn(len(values))]
}

func (d *Database) GetRawValuesByAuthor(guild string, author string) []RawValue {
	result := []RawValue{}
	if _, ok := d.raw[guild]; !ok {
		return nil
	}
	for _, v := range d.raw[guild] {
		if v.Author == author {
			result = append(result, v)
		}
	}
	return result
}

func (d *Database) Init() {
	rawData := SpreadSheetInstance().RawValues
	d.raw = make(map[string]map[uuid.UUID]RawValue)
	for k, v := range rawData {
		d.raw[k] = make(map[uuid.UUID]RawValue)
		for _, value := range v {
			d.raw[k][value.Id] = value
		}
	}
	d.initDic()
}

func (d *Database) AddValue(content *RawValue, guild string) {
	if _, ok := d.dic[guild]; !ok {
		d.dic[guild] = make(map[string]map[uuid.UUID]struct{})
	}

	for _, tag := range content.Tags {
		gDic := d.dic[guild]
		if _, ok := gDic[tag]; !ok {
			gDic[tag] = make(map[uuid.UUID]struct{})
		}
		tSls := gDic[tag]
		if _, ok := tSls[content.Id]; !ok {
			tSls[content.Id] = struct{}{}
		}
	}

	if _, ok := d.raw[guild]; !ok {
		d.raw[guild] = make(map[uuid.UUID]RawValue)
	}
	if _, ok := d.raw[guild][content.Id]; !ok {
		d.raw[guild][content.Id] = *content
	}

}

func intersection(a map[uuid.UUID]struct{}, b map[uuid.UUID]struct{}) map[uuid.UUID]struct{} {
	result := make(map[uuid.UUID]struct{})
	if len(a) == 0 || len(a) == 0 {
		return result
	}
	for id := range b {
		if _, ok := a[id]; ok {
			result[id] = struct{}{}
		}
	}
	return result
}

func (d *Database) initDic() {
	d.dic = make(map[string]map[string]map[uuid.UUID]struct{})
	for k, v := range d.raw {
		for _, raw := range v {
			d.AddValue(&raw, k)
		}
	}
}

func newDb() *Database {
	db := Database{}
	db.raw = make(map[string]map[uuid.UUID]RawValue)
	db.dic = make(map[string]map[string]map[uuid.UUID]struct{})
	return &db
}
