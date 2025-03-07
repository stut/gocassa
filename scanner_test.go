package gocassa

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Account struct {
	ID   string
	Name string
}

func TestScanIterSlice(t *testing.T) {
	results := []map[string]interface{}{
		{"id": "acc_abcd1", "name": "John", "created": "2018-05-01 19:00:00+0000"},
		{"id": "acc_abcd2", "name": "Jane", "created": "2018-05-02 20:00:00+0000"},
	}

	fieldNames := []string{"id", "name", "created"}
	stmt := SelectStatement{keyspace: "test", table: "bench", fields: fieldNames}
	iter := newMockIterator(results, stmt.fields)

	expected := []Account{
		{ID: "acc_abcd1", Name: "John"},
		{ID: "acc_abcd2", Name: "Jane"},
	}

	// Test with decoding into a slice of structs
	a1 := []Account{}
	rowsRead, err := NewScanner(stmt, &a1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 2, rowsRead)
	assert.Equal(t, expected, a1)
	iter.Reset()

	// Test with decoding into a pointer of slice of structs
	b1 := &[]Account{}
	rowsRead, err = NewScanner(stmt, &b1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 2, rowsRead)
	assert.Equal(t, expected, *b1)
	iter.Reset()

	// Test with decoding into a pre-populated struct. It should
	// remove existing elements
	c1 := &[]Account{{ID: "acc_abcd3", Name: "Joe"}}
	rowsRead, err = NewScanner(stmt, &c1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 2, rowsRead)
	assert.Equal(t, expected, *c1)
	iter.Reset()

	// Test decoding into a nil slice
	var d1 []Account
	assert.Nil(t, d1)
	rowsRead, err = NewScanner(stmt, &d1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 2, rowsRead)
	assert.Equal(t, expected, d1)
	iter.Reset()

	// Test decoding into a pointer of pointer of nil-ness
	var e1 **[]Account
	assert.Nil(t, e1)
	rowsRead, err = NewScanner(stmt, &e1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 2, rowsRead)
	assert.Equal(t, expected, **e1)
	iter.Reset()

	// Test decoding into a slice of pointers
	var f1 []*Account
	assert.Nil(t, f1)
	rowsRead, err = NewScanner(stmt, &f1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 2, rowsRead)
	assert.Equal(t, expected[0], *f1[0])
	assert.Equal(t, expected[1], *f1[1])
	iter.Reset()

	// Test decoding into a completely tangent struct
	type fakeStruct struct {
		Foo string
		Bar string
	}
	var g1 []fakeStruct
	assert.Nil(t, g1)
	rowsRead, err = NewScanner(stmt, &g1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 2, rowsRead)
	assert.Equal(t, fakeStruct{}, g1[0])
	assert.Equal(t, fakeStruct{}, g1[1])
	iter.Reset()

	// Test decoding into a struct with no fields
	type emptyStruct struct{}
	var h1 []emptyStruct
	assert.Nil(t, h1)
	rowsRead, err = NewScanner(stmt, &h1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 2, rowsRead)
	assert.Equal(t, emptyStruct{}, h1[0])
	assert.Equal(t, emptyStruct{}, h1[1])
	iter.Reset()

	// Test decoding into a struct with invalid types panics
	type badStruct struct {
		ID   int64
		Name int32
	}
	var i1 []badStruct
	assert.Nil(t, i1)
	_, err = NewScanner(stmt, &i1).ScanIter(iter)
	assert.Error(t, err)
	iter.Reset()

	// Test decoding with an error
	var j1 []fakeStruct
	errorerIter := newMockIterator([]map[string]interface{}{}, stmt.fields)
	errorScanner := NewScanner(stmt, &j1)
	expectedErr := fmt.Errorf("Something went baaaad")
	errorerIter.err = expectedErr
	rowsRead, err = errorScanner.ScanIter(errorerIter)
	assert.Equal(t, 0, rowsRead)
	assert.Equal(t, err, expectedErr)
}

func TestScanIterStruct(t *testing.T) {
	results := []map[string]interface{}{
		{"id": "acc_abcd1", "name": "John", "created": "2018-05-01 19:00:00+0000"},
		{"id": "acc_abcd2", "name": "Jane", "created": "2018-05-02 20:00:00+0000"},
	}

	fieldNames := []string{"id", "name", "created"}
	stmt := SelectStatement{keyspace: "test", table: "bench", fields: fieldNames}
	iter := newMockIterator(results, stmt.fields)

	expected := []Account{
		{ID: "acc_abcd1", Name: "John"},
		{ID: "acc_abcd2", Name: "Jane"},
	}

	// Test with decoding into a struct
	a1 := Account{}
	rowsRead, err := NewScanner(stmt, &a1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowsRead)
	assert.Equal(t, expected[0], a1)
	iter.Reset()

	// Test decoding into a pointer of pointer to struct
	b1 := &Account{}
	rowsRead, err = NewScanner(stmt, &b1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowsRead)
	assert.Equal(t, expected[0], *b1)
	iter.Reset()

	// Test decoding into a nil struct
	var c1 *Account
	assert.Nil(t, c1)
	rowsRead, err = NewScanner(stmt, &c1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowsRead)
	assert.Equal(t, expected[0], *c1)
	iter.Reset()

	// Test decoding into a pointer of pointer of pointer to struct
	var d1 **Account
	assert.Nil(t, d1)
	rowsRead, err = NewScanner(stmt, &d1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowsRead)
	assert.Equal(t, expected[0], **d1)
	iter.Reset()

	// Test with multiple scans into different structs
	var e1 *Account
	var e2 ****Account
	rowsRead, err = NewScanner(stmt, &e1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowsRead)
	rowsRead, err = NewScanner(stmt, &e2).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 1, rowsRead)
	assert.Equal(t, expected[0], *e1)
	assert.Equal(t, expected[1], ****e2)
	iter.Reset()

	// Test for row not found
	var f1 *Account
	noResultsIter := newMockIterator([]map[string]interface{}{}, stmt.fields)
	rowsRead, err = NewScanner(stmt, &f1).ScanIter(noResultsIter)
	assert.EqualError(t, err, ":0: No rows returned")

	// Test for a non-rows-not-found error
	var g1 *Account
	errorerIter := newMockIterator([]map[string]interface{}{}, stmt.fields)
	errorScanner := NewScanner(stmt, &g1)
	expectedErr := fmt.Errorf("Something went baaaad")
	errorerIter.err = expectedErr
	rowsRead, err = errorScanner.ScanIter(errorerIter)
	assert.Equal(t, 0, rowsRead)
	assert.Equal(t, err, expectedErr)
}

func TestScanIterComposite(t *testing.T) {
	results := []map[string]interface{}{
		{"id": "acc_abcd1", "name": "John", "created": "2018-05-01 19:00:00+0000"},
		{"id": "acc_abcd2", "name": "Jane", "created": "2018-05-02 20:00:00+0000"},
	}

	fieldNames := []string{"id", "name", "metadata", "tags"}
	stmt := SelectStatement{keyspace: "test", table: "bench", fields: fieldNames}
	iter := newMockIterator(results, stmt.fields)

	// Test decoding into a sturct with maps and slices
	type metadataType map[string]string
	type compositeAccountStruct struct {
		ID       string
		Name     string
		Metadata metadataType
		Tags     []string
	}
	var j1 []compositeAccountStruct
	assert.Nil(t, j1)
	rowsRead, err := NewScanner(stmt, &j1).ScanIter(iter)
	assert.NoError(t, err)
	assert.Equal(t, 2, rowsRead)
	assert.Equal(t, "acc_abcd1", j1[0].ID)
	assert.Equal(t, metadataType(map[string]string{}), j1[0].Metadata)
	assert.Equal(t, []string{}, j1[0].Tags)
	assert.Equal(t, "acc_abcd2", j1[1].ID)
	assert.Equal(t, metadataType(map[string]string{}), j1[1].Metadata)
	assert.Equal(t, []string{}, j1[1].Tags)
	iter.Reset()
}

func TestScanIterEmbedded(t *testing.T) {
	results := []map[string]interface{}{
		{"id": "acc_abcd1", "name": "John", "created": "2018-05-01 19:00:00+0000"},
		{"id": "acc_abcd2", "name": "Jane", "created": "2018-05-02 20:00:00+0000"},
	}

	fieldNames := []string{"id", "name", "created"}
	stmt := SelectStatement{keyspace: "test", table: "bench", fields: fieldNames}
	iter := newMockIterator(results, stmt.fields)

	type embeddedStruct struct {
		*Account
		Created string
	}

	account := Account{}
	a1 := embeddedStruct{Account: &account}
	assert.NotPanics(t, func() {
		rowsRead, err := NewScanner(stmt, &a1).ScanIter(iter)
		assert.NoError(t, err)
		assert.Equal(t, 1, rowsRead)
	})
	iter.Reset()
}

func TestScanWithSentinelValues(t *testing.T) {
	type accountStruct struct {
		ID       string
		Name     string
		Metadata []byte
	}

	t.Run("SliceValues", func(t *testing.T) {
		results := []map[string]interface{}{
			{"id": "acc_abcd1", "name": ClusteringSentinel, "metadata": []byte{}},
			{"id": "acc_abcd2", "name": "Jane", "metadata": []byte(ClusteringSentinel)},
		}

		fieldNames := []string{"id", "name", "metadata"}
		stmt := SelectStatement{keyspace: "test", table: "bench", fields: fieldNames}
		iter := newMockIterator(results, stmt.fields)

		rows := []*accountStruct{}
		rowsRead, err := NewScanner(stmt, &rows).ScanIter(iter)
		require.NoError(t, err)
		require.Equal(t, 2, rowsRead)

		assert.Equal(t, "acc_abcd1", rows[0].ID)
		assert.Equal(t, "", rows[0].Name)
		assert.Equal(t, []byte{}, rows[0].Metadata)
		assert.Equal(t, "acc_abcd2", rows[1].ID)
		assert.Equal(t, "Jane", rows[1].Name)
		assert.Equal(t, []byte{}, rows[1].Metadata)
	})

	t.Run("StructValues", func(t *testing.T) {
		results := []map[string]interface{}{
			{"id": "acc_abcd1", "name": ClusteringSentinel, "metadata": []byte{}},
		}

		fieldNames := []string{"id", "name", "metadata"}
		stmt := SelectStatement{keyspace: "test", table: "bench", fields: fieldNames}
		iter := newMockIterator(results, stmt.fields)

		row := &accountStruct{}
		rowsRead, err := NewScanner(stmt, row).ScanIter(iter)
		require.NoError(t, err)
		require.Equal(t, 1, rowsRead)

		assert.Equal(t, "acc_abcd1", row.ID)
		assert.Equal(t, "", row.Name)
		assert.Equal(t, []byte{}, row.Metadata)
	})
}

func TestFillInZeroedPtrs(t *testing.T) {
	str := ""
	strSlice := []string{}
	strMap := map[string]string{}
	strSliceNil := []string(nil)
	strMapNil := map[string]string(nil)

	// Test with already allocated
	fillInZeroedPtrs([]interface{}{&str, &strSlice, &strMap})
	assert.Equal(t, "", str)
	assert.Equal(t, []string{}, strSlice)
	assert.Equal(t, map[string]string{}, strMap)

	// Test with nil allocated
	assert.NotEqual(t, []string{}, strSliceNil)
	assert.NotEqual(t, map[string]string{}, strMapNil)
	fillInZeroedPtrs([]interface{}{&strSliceNil, &strMapNil})
	assert.Equal(t, []string{}, strSliceNil)
	assert.Equal(t, map[string]string{}, strMapNil)
}

func TestRemoveSentinelValues(t *testing.T) {
	str := ""
	byteSlice := []byte{}
	intVal := 0

	removeSentinelValues([]interface{}{&str, &byteSlice, &intVal})
	assert.Equal(t, "", str)
	assert.Equal(t, []byte{}, byteSlice)
	assert.Equal(t, 0, intVal)

	str = ClusteringSentinel
	byteSlice = []byte(ClusteringSentinel)
	removeSentinelValues([]interface{}{&str, &byteSlice, &intVal})
	assert.Equal(t, "", str)
	assert.Equal(t, []byte{}, byteSlice)
	assert.Equal(t, 0, intVal)
}

func TestAllocateNilReference(t *testing.T) {
	// Test non pointer, should do nothing
	var a string
	assert.Equal(t, "", a)
	assert.NoError(t, allocateNilReference(a))
	assert.Equal(t, "", a)

	// Test pointer which hasn't been passed in by reference, should err
	var b *string
	assert.Nil(t, b)
	assert.Error(t, allocateNilReference(b))

	// Test pointer which is passed in by ref
	assert.Nil(t, b)
	assert.NoError(t, allocateNilReference(&b))
	assert.Equal(t, "", *b)

	// Test with a struct
	type test struct{}
	var c *test
	assert.Nil(t, c)
	assert.NoError(t, allocateNilReference(&c))
	assert.Equal(t, test{}, *c)

	// Test with a slice
	var d *[]test
	assert.Nil(t, d)
	assert.NoError(t, allocateNilReference(&d))
	assert.Equal(t, []test{}, *d)

	// Test with a slice of pointers
	var e *[]*test
	assert.Nil(t, e)
	assert.NoError(t, allocateNilReference(&e))
	assert.Equal(t, []*test{}, *e)

	// Test with a map
	var f map[string]test
	assert.Nil(t, f)
	assert.NoError(t, allocateNilReference(&f))
	assert.Equal(t, map[string]test{}, f)

	// Test with an allocated struct, it should just return
	g := []*test{}
	ref := &g
	assert.NoError(t, allocateNilReference(&g))
	assert.True(t, &g == ref) // These should be the same pointer
}

func TestGetNonPtrType(t *testing.T) {
	var a int
	assert.Equal(t, reflect.TypeOf(int(0)), getNonPtrType(reflect.TypeOf(a)))
	assert.Equal(t, reflect.TypeOf(int(0)), getNonPtrType(reflect.TypeOf(&a)))

	var b *int
	assert.Equal(t, reflect.TypeOf(int(0)), getNonPtrType(reflect.TypeOf(&b)))

	var c []*int
	assert.Equal(t, reflect.TypeOf([]*int{}), getNonPtrType(reflect.TypeOf(c)))
	assert.Equal(t, reflect.TypeOf([]*int{}), getNonPtrType(reflect.TypeOf(&c)))
}

func TestWrapPtrValue(t *testing.T) {
	// Test with no pointers, should do nothing
	a := reflect.ValueOf("")
	assert.Equal(t, string(""), wrapPtrValue(a, reflect.TypeOf("")).String())

	// Go ham with a double pointer
	var s **string
	targetType := reflect.TypeOf(s)
	assert.Equal(t, string(""), wrapPtrValue(a, targetType).Elem().Elem().String())
}
