package decode

import (
	"errors"
	"reflect"
	"testing"
)

func TestDecodeStructMap(t *testing.T) {
	type nested struct {
		Field  string      `copy:"field"`
		Custom interface{} `copy:"custom"`
	}

	testIn := struct {
		Name   string `copy:"name"`
		Age    int    `copy:"age"`
		Nested nested `copy:"nested"`
		NoCopy string
	}{
		Name: "John",
		Age:  30,
		Nested: nested{
			Field:  "value",
			Custom: &[]int{1, 2, 3},
		},
		NoCopy: "value",
	}

	t.Run("test struct to map strong type", func(t *testing.T) {
		testOut := make(map[string]interface{})

		if err := Decode(testIn, &testOut, "copy", 0); err != nil {
			t.Errorf("Decode() error = %v", err)
			return
		}

		if testIn.Age != testOut["age"] || testIn.Name != testOut["name"] || testIn.Nested.Field != testOut["nested"].(nested).Field ||
			testIn.Nested.Custom == testOut["nested"].(nested).Custom || len(testOut) != 3 ||
			!reflect.DeepEqual(*(testIn.Nested.Custom.(*[]int)), *(testOut["nested"].(nested).Custom.(*[]int))) {
			t.Errorf("Decode() = %v, want %v", testOut, testIn)
		}
	})

	t.Run("test struct to map not strong type", func(t *testing.T) {
		testOut := make(map[string]interface{})

		if err := Decode(testIn, &testOut, "copy", DecoderUnwrapStructToMap); err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if testIn.Age != testOut["age"] || testIn.Name != testOut["name"] || len(testOut) != 3 ||
			testIn.Nested.Field != testOut["nested"].(map[string]interface{})["field"] ||
			testIn.Nested.Custom == testOut["nested"].(map[string]interface{})["custom"] {
			t.Errorf("Decode() = %v, want %v", testOut, testIn)
		}
	})

	t.Run("test struct to map not strong type by name", func(t *testing.T) {
		testOut := make(map[string]interface{})

		if err := Decode(testIn, &testOut, "", DecoderUnwrapStructToMap); err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if testIn.Age != testOut["Age"] || testIn.Name != testOut["Name"] || testIn.NoCopy != testOut["NoCopy"] ||
			len(testOut) != 4 || testIn.Nested.Field != testOut["Nested"].(map[string]interface{})["Field"] ||
			testIn.Nested.Custom == testOut["Nested"].(map[string]interface{})["Custom"] {
			t.Errorf("Decode() = %v, want %v", testOut, testIn)
		}
	})
}

func TestDecodeStructStruct(t *testing.T) {
	testIn := struct {
		Name   string `copy:"name"`
		Age    int    `copy:"age"`
		Year   int    `copy:"year"`
		Nested struct {
			Field  string `copy:"field"`
			Field2 int    `copy:"field2"`
		} `copy:"nested"`
	}{
		Name: "John",
		Age:  30,
		Year: 2024,
		Nested: struct {
			Field  string `copy:"field"`
			Field2 int    `copy:"field2"`
		}{
			Field:  "value",
			Field2: 42,
		},
	}
	testOut := struct {
		Name   string `copy:"name"`
		Age    int    `copy:"age"`
		Nested struct {
			Field string `copy:"field"`
		} `copy:"nested"`
	}{}
	wantOut := struct {
		Name   string `copy:"name"`
		Age    int    `copy:"age"`
		Nested struct {
			Field string `copy:"field"`
		} `copy:"nested"`
	}{
		Name: "John",
		Age:  30,
		Nested: struct {
			Field string `copy:"field"`
		}{Field: "value"},
	}

	t.Run("test struct to struct by tag", func(t *testing.T) {
		if err := Decode(testIn, &testOut, "copy", 0); err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if !reflect.DeepEqual(testOut, wantOut) {
			t.Errorf("Decode() = %v, want %v", testOut, wantOut)
		}
	})

	t.Run("test struct to struct by tag not found dst", func(t *testing.T) {
		if err := Decode(testIn, &testOut, "copy", DecoderStrongFoundDst); !errors.Is(err, ErrorDstNotFound) {
			t.Errorf("Decode() error = %v", err)
		}
	})
}

func TestDecodeStructStructByName(t *testing.T) {
	testIn := struct {
		Name   string
		Age    int
		Year   int
		Nested struct {
			Field  string
			Field2 int
		}
	}{
		Name: "John",
		Age:  30,
		Year: 2024,
		Nested: struct {
			Field  string
			Field2 int
		}{
			Field:  "value",
			Field2: 42,
		},
	}
	wantOut := struct {
		Name   string
		Age    int
		Nested struct {
			Field string
		}
	}{
		Name: "John",
		Age:  30,
		Nested: struct {
			Field string
		}{Field: "value"},
	}

	t.Run("test struct to struct by field name", func(t *testing.T) {
		testOut := struct {
			Name   string
			Age    int
			Nested struct {
				Field string
			}
		}{}

		if err := Decode(testIn, &testOut, "", 0); err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if !reflect.DeepEqual(testOut, wantOut) {
			t.Errorf("Decode() = %v, want %v", testOut, wantOut)
		}
	})

	t.Run("test struct to struct by field name not found dst", func(t *testing.T) {
		testOut := struct {
			Name   string
			Age    int
			Nested struct {
				Field string
			}
		}{}

		if err := Decode(testIn, &testOut, "", DecoderStrongFoundDst); !errors.Is(err, ErrorDstNotFound) {
			t.Errorf("Decode() error = %v", err)
		}
	})
}

func TestDecodeMapStruct(t *testing.T) {
	type nested struct {
		Field  string      `copy:"field"`
		Custom interface{} `copy:"custom"`
	}

	type testStruct struct {
		Name   string `copy:"name"`
		Age    int    `copy:"age"`
		Nested nested `copy:"nested"`
		NoCopy string
	}

	testIn := map[string]interface{}{
		"name": "John",
		"age":  30,
		"nested": nested{
			Field:  "value",
			Custom: &[]int{1, 2, 3},
		},
		"noCopy": "value",
	}

	t.Run("test map to struct strong type", func(t *testing.T) {
		testOut := testStruct{}

		if err := Decode(testIn, &testOut, "copy", 0); err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if testOut.Age != testIn["age"] || testOut.Name != testIn["name"] || !reflect.DeepEqual(testOut.Nested, testIn["nested"]) ||
			testOut.Nested.Custom == testIn["nested"].(nested).Custom {
			t.Errorf("Decode() = %v, want %v", testOut, testIn)
		}
	})

	t.Run("test full map to struct strong type", func(t *testing.T) {
		testOut := testStruct{}
		testIn := map[string]interface{}{
			"name": "John",
			"age":  float64(30),
			"nested": map[string]interface{}{
				"field":  "value",
				"custom": &[]int{1, 2, 3},
			},
			"noCopy": "value",
		}

		if err := Decode(testIn, &testOut, "copy", 0); err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if testOut.Age != int(testIn["age"].(float64)) || testOut.Name != testIn["name"] ||
			testOut.Nested.Custom == testIn["nested"].(map[string]interface{})["custom"] ||
			testOut.Nested.Field != testIn["nested"].(map[string]interface{})["field"] {
			t.Errorf("Decode() = %v, want %v", testOut, testIn)
		}
	})

	t.Run("test full map to non struct strong type", func(t *testing.T) {
		testOut := testStruct{}
		testIn := map[string]interface{}{
			"name": 100,
			"age":  "30",
			"nested": map[string]interface{}{
				"field":  "value",
				"custom": &[]int{1, 2, 3},
			},
			"noCopy": "value",
		}

		if err := Decode(testIn, &testOut, "copy", 0); err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if testOut.Age != 30 || testOut.Name != "100" ||
			testOut.Nested.Custom == testIn["nested"].(map[string]interface{})["custom"] ||
			testOut.Nested.Field != testIn["nested"].(map[string]interface{})["field"] {
			t.Errorf("Decode() = %v, want %v", testOut, testIn)
		}
	})

	t.Run("test map to struct deep copy", func(t *testing.T) {
		testOut := testStruct{}

		if err := Decode(testIn, &testOut, "copy", 0); err != nil {
			t.Errorf("Decode() error = %v", err)
			return
		}

		if testOut.Age != testIn["age"] || testOut.Name != testIn["name"] || testOut.Nested.Field != testIn["nested"].(nested).Field ||
			testOut.Nested.Custom == testIn["nested"].(nested).Custom ||
			!reflect.DeepEqual(*(testOut.Nested.Custom.(*[]int)), *(testIn["nested"].(nested).Custom.(*[]int))) {
			t.Errorf("Decode() = %v, want %v", testOut, testIn)
		}
	})
}

func TestDecodeMapMap(t *testing.T) {
	type nested struct {
		Field  string      `copy:"field"`
		Custom interface{} `copy:"custom"`
	}

	testIn := map[string]interface{}{
		"name": "John",
		"age":  30,
		"nested": nested{
			Field:  "value",
			Custom: &[]int{1, 2, 3},
		},
		"noCopy": "value",
	}

	t.Run("test map to struct strong type", func(t *testing.T) {
		testOut := map[string]interface{}{}

		if err := Decode(testIn, &testOut, "copy", DecoderStrongType); err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if testOut["age"] != testIn["age"] || testOut["name"] != testIn["name"] || !reflect.DeepEqual(testOut["nested"], testIn["nested"]) ||
			testOut["nested"].(nested).Custom != testIn["nested"].(nested).Custom {
			t.Errorf("Decode() = %v, want %v", testOut, testIn)
		}
	})
}

func TestDecodeMapMapTypes(t *testing.T) {
	testIn := map[string]interface{}{
		"1": "1",
		"2": 2.1,
		"3": true,
	}

	t.Run("test map to map convert to int", func(t *testing.T) {
		testOut := map[string]int64{}
		want := map[string]int64{
			"1": 1,
			"2": 2,
			"3": 1,
		}

		if err := Decode(testIn, &testOut, "copy", 0); err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if reflect.DeepEqual(testOut, want) {
			t.Errorf("Decode() = %v, want %v", testOut, want)
		}
	})

	t.Run("test map to map convert to float", func(t *testing.T) {
		testOut := map[string]float64{}
		want := map[string]float64{
			"1": 1.0,
			"2": 2.1,
			"3": 1.0,
		}

		if err := Decode(testIn, &testOut, "copy", 0); err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if reflect.DeepEqual(testOut, want) {
			t.Errorf("Decode() = %v, want %v", testOut, want)
		}
	})
}
