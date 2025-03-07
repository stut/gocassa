package gocassa

type multimapMkT struct {
	t               Table
	fieldsToIndexBy []string
	idField         []string
}

func (mm *multimapMkT) Table() Table                        { return mm.t }
func (mm *multimapMkT) Create() error                       { return mm.Table().Create() }
func (mm *multimapMkT) CreateIfNotExist() error             { return mm.Table().CreateIfNotExist() }
func (mm *multimapMkT) Name() string                        { return mm.Table().Name() }
func (mm *multimapMkT) Recreate() error                     { return mm.Table().Recreate() }
func (mm *multimapMkT) CreateStatement() (Statement, error) { return mm.Table().CreateStatement() }
func (mm *multimapMkT) CreateIfNotExistStatement() (Statement, error) {
	return mm.Table().CreateIfNotExistStatement()
}

func (mm *multimapMkT) Update(field, id map[string]interface{}, m map[string]interface{}) Op {
	return mm.Table().
		Where(mm.ListOfEqualRelations(field, id)...).
		Update(m)
}

func (mm *multimapMkT) Set(v interface{}) Op {
	return mm.Table().
		Set(v)
}

func (mm *multimapMkT) Delete(field, id map[string]interface{}) Op {
	return mm.Table().
		Where(mm.ListOfEqualRelations(field, id)...).
		Delete()
}

func (mm *multimapMkT) DeleteAll(field map[string]interface{}) Op {
	return mm.Table().
		Where(mm.ListOfEqualRelations(field, nil)...).
		Delete()
}

func (mm *multimapMkT) Read(field, id map[string]interface{}, pointer interface{}) Op {
	return mm.Table().
		Where(mm.ListOfEqualRelations(field, id)...).
		ReadOne(pointer)
}

func (mm *multimapMkT) MultiRead(field, id map[string]interface{}, pointerToASlice interface{}) Op {
	return mm.Table().
		Where(mm.ListOfEqualRelations(field, id)...).
		Read(pointerToASlice)
}

func (mm *multimapMkT) List(field, startId map[string]interface{}, limit int, pointerToASlice interface{}) Op {
	rels := mm.ListOfEqualRelations(field, nil)
	if startId != nil {
		for _, field := range mm.idField {
			if value := startId[field]; value != "" {
				rels = append(rels, GTE(field, value))
			}
		}
	}
	return mm.
		WithOptions(Options{
			Limit: limit,
		}).
		Table().
		Where(rels...).
		Read(pointerToASlice)
}

func (mm *multimapMkT) WithOptions(o Options) MultimapMkTable {
	return &multimapMkT{
		t:               mm.Table().WithOptions(o),
		fieldsToIndexBy: mm.fieldsToIndexBy,
		idField:         mm.idField,
	}
}

func (mm *multimapMkT) ListOfEqualRelations(fieldsToIndex, ids map[string]interface{}) []Relation {
	relations := make([]Relation, 0)

	for _, field := range mm.fieldsToIndexBy {
		if value := fieldsToIndex[field]; value != nil && value != "" {
			relation := Eq(field, value)
			relations = append(relations, relation)
		}
	}

	for _, field := range mm.idField {
		if value := ids[field]; value != nil && value != "" {
			relation := Eq(field, value)
			relations = append(relations, relation)
		}
	}

	return relations
}
