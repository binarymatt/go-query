package query

import (
	"testing"
	"time"

	"github.com/shoenig/test/must"
)

type testModel struct {
	Id          int
	private     string
	SampleCase  string
	Tagged      float64 `db:"tagged_field"`
	Ignored     string  `db:"-"`
	CreatedAt   time.Time
	DeletedAt   time.Time `db:"deleted,omitempty"`
	AnotherTime time.Time
	Simple      int
}

func TestReflectDetails(t *testing.T) {
	now := time.Now()
	var zeroTime time.Time
	model := testModel{
		Id:         1,
		private:    "private string",
		SampleCase: "sample",
		Tagged:     1.2,
		Ignored:    "should be absent",
		CreatedAt:  now,
	}
	expectedFields := []field{
		{
			name:      "id",
			value:     1,
			omitEmpty: false,
			tagged:    false,
			excluded:  false,
			isZero:    false,
		},
		{
			name:      "sample_case",
			value:     "sample",
			omitEmpty: false,
			tagged:    false,
			excluded:  false,
			isZero:    false,
		},
		{
			name:      "tagged_field",
			value:     1.2,
			omitEmpty: false,
			tagged:    true,
			excluded:  false,
			isZero:    false,
		},
		{
			name:      "ignored",
			value:     "should be absent",
			omitEmpty: false,
			tagged:    true,
			excluded:  true,
			isZero:    false,
		},
		{
			name:      "created_at",
			value:     now,
			omitEmpty: false,
			tagged:    false,
			excluded:  false,
			isZero:    false,
		},
		{
			name:      "deleted",
			value:     zeroTime,
			omitEmpty: true,
			tagged:    true,
			excluded:  false,
			isZero:    true,
		},
		{
			name:      "another_time",
			value:     zeroTime,
			omitEmpty: false,
			tagged:    false,
			excluded:  false,
			isZero:    true,
		},
		{
			name:      "simple",
			value:     0,
			omitEmpty: false,
			tagged:    false,
			excluded:  false,
			isZero:    true,
		},
	}
	name, fields := reflectDetails(model)
	must.Eq(t, "test_models", name)
	for i, field := range fields {
		must.Eq(t, expectedFields[i], field)
	}
	must.Eq(t, expectedFields, fields)
}
