package gocassa

import (
	"reflect"
	"testing"
)

type Customer2 struct {
	Id   string
	Name string
	Tag  string
}

func TestMultimapTableInsertRead(t *testing.T) {
	tbl := ns.MultimapTable("customer91", "Tag", "Id", Customer2{})
	createIf(tbl.(TableChanger), t)
	joe := Customer2{
		Id:   "33",
		Name: "Joe",
		Tag:  "A",
	}
	err := tbl.Set(joe).Run()
	if err != nil {
		t.Fatal(err)
	}
	res := &Customer2{}
	err = tbl.Read("A", "33", res).Run()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*res, joe) {
		t.Fatal(*res, joe)
	}
	err = tbl.Read("B", "33", res).Run()
	if err == nil {
		t.Fatal(*res)
	}
	err = tbl.Update("A", "33", map[string]interface{}{
		"Name": "John",
	}).Run()
	if err != nil {
		t.Fatal(err)
	}
	err = tbl.Read("A", "33", res).Run()
	if err != nil {
		t.Fatal(err)
	}
	if res.Name != "John" {
		t.Fatal(*res)
	}
	list := &[]Customer{}
	err = tbl.List("A", nil, 20, list).Run()
	if err != nil {
		t.Fatal(err)
	}
	if len(*list) != 1 {
		t.Fatal(*list)
	}
}

func TestMultimapTableDelete(t *testing.T) {
	tbl := ns.MultimapTable("customer92", "Tag", "Id", Customer2{})
	createIf(tbl.(TableChanger), t)
	joe := Customer2{
		Id:   "33",
		Name: "Joe",
		Tag:  "A",
	}
	err := tbl.Set(joe).Run()
	if err != nil {
		t.Fatal(err)
	}
	res := &Customer2{}
	err = tbl.Read("A", "33", res).Run()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*res, joe) {
		t.Fatal(*res, joe)
	}
	err = tbl.Delete("A", "33").Run()
	if err != nil {
		t.Fatal(err)
	}
	err = tbl.Read("A", "33", res).Run()
	if err == nil {
		t.Fatal(res)
	}
}

func TestMultimapTableMultiListOrder(t *testing.T) {
	tbl := ns.MultimapTable("customer93", "Tag", "Id", Customer2{})
	createIf(tbl.(TableChanger), t)
	joe := Customer2{
		Id:   "33",
		Name: "Joe",
		Tag:  "A",
	}
	err := tbl.Set(joe).Run()
	if err != nil {
		t.Fatal(err)
	}
	jane := Customer2{
		Id:   "34",
		Name: "Jane",
		Tag:  "A",
	}
	err = tbl.Set(jane).Run()
	if err != nil {
		t.Fatal(err)
	}
	customers := []Customer2{}
	op := tbl.List("A", nil, 20, &customers).WithOptions(Options{
		ClusteringOrder: []ClusteringOrderColumn{
			{DESC, "Id"},
		},
	})

	err = op.Run()
	if err != nil {
		t.Fatal(err)
	}
	if len(customers) != 2 {
		t.Fatalf("Expected to multiread 2 records, got %d", len(customers))
	}
	if !reflect.DeepEqual((customers)[0], jane) {
		t.Fatalf("Expected to find jane, got %v", (customers)[0])
	}
	if !reflect.DeepEqual((customers)[1], joe) {
		t.Fatalf("Expected to find joe, got %v", (customers)[1])
	}
}
