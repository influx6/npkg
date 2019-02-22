package nreflect_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	nreflect "github.com/gokit/reflection"
)

type bull string

type speaker interface {
	Speak() string
}

// mosnter provides a basic struct test case type.
type monster struct {
	Name  string
	Items []bull
}

// Speak returns the sound the monster makes.
func (m *monster) Speak() string {
	return "Raaaaaaarggg!"
}

func get(t *testing.T, sm speaker) {
	name, embedded, err := nreflect.ExternalTypeNames(sm)
	if err != nil {
		t.Fatalf("Should be able to retrieve field names arguments lists\n")
	}
	t.Logf("Name: %s", name)
	t.Logf("Fields: %+q", embedded)
	t.Logf("Should be able to retrieve function arguments lists")
}

type Addrs struct {
	Addr string
}

type addrFunc func(Addrs) error

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func TestIsSettableType(t *testing.T) {
	m2 := errors.New("invalid error")
	if !nreflect.IsSettableType(errorType, reflect.TypeOf(m2)) {
		t.Fatalf("Should have error matching type")
	}
}

func TestIsSettable(t *testing.T) {
	m2 := errors.New("invalid error")
	if !nreflect.IsSettable(errorType, reflect.ValueOf(m2)) {
		t.Fatalf("Should have error matching type")
	}

	m1 := mo("invalid error")
	if !nreflect.IsSettable(errorType, reflect.ValueOf(m1)) {
		t.Fatalf("Should have error matching type")
	}
}

type mo string

func (m mo) Error() string {
	return string(m)
}

func TestValidateFunc_Bad(t *testing.T) {
	var testFunc = func(v string) string {
		return "Hello " + v
	}

	err := nreflect.ValidateFunc(testFunc, []nreflect.TypeValidation{
		func(types []reflect.Type) error {
			if len(types) == 1 {
				return nil
			}
			return errors.New("bad")
		},
	}, []nreflect.TypeValidation{
		func(types []reflect.Type) error {
			if len(types) == 0 {
				return nil
			}
			return errors.New("bad")
		},
	})

	if err == nil {
		t.Fatalf("Should have function invalid to conditions")
	}
}

func TestValidateFunc(t *testing.T) {
	var testFunc = func(v string) string {
		return "Hello " + v
	}

	err := nreflect.ValidateFunc(testFunc, []nreflect.TypeValidation{
		func(types []reflect.Type) error {
			if len(types) == 1 {
				return nil
			}
			return errors.New("bad")
		},
	}, []nreflect.TypeValidation{
		func(types []reflect.Type) error {
			if len(types) == 1 {
				return nil
			}
			return errors.New("bad")
		},
	})

	if err != nil {
		t.Fatalf("Should have function valid to conditions")
	}
}

func TestFunctionApply_OneArgument(t *testing.T) {
	var testFunc = func(v string) string {
		return "Hello " + v
	}

	res, err := nreflect.CallFunc(testFunc, "Alex")
	if err != nil {
		t.Fatalf("Should have executed function")
	}
	t.Logf("Should have executed function")

	if !reflect.DeepEqual(res, []interface{}{"Hello Alex"}) {
		t.Logf("Received: %q", res)
		t.Fatalf("Expected value unmatched")
	}
}

func TestFunctionApply_ThreeArgumentWithError(t *testing.T) {
	bad := errors.New("bad")
	var testFunc = func(v string, i int, d bool) ([]interface{}, error) {
		return []interface{}{v, i, d}, bad
	}

	res, err := nreflect.CallFunc(testFunc, "Alex", 1, false)
	if err != nil {
		t.Fatalf("Should have executed function")
	}
	t.Logf("Should have executed function")

	if !reflect.DeepEqual(res, []interface{}{[]interface{}{"Alex", 1, false}, bad}) {
		t.Logf("Expected: %q", []interface{}{[]interface{}{"Alex", 1, false}, bad})
		t.Logf("Received: %q", res)
		t.Fatalf("Expected value unmatched")
	}
}

func TestFunctionApply_ThreeArgumentWithVariadic(t *testing.T) {
	var testFunc = func(v string, i int, d ...bool) []interface{} {
		return []interface{}{v, i, d}
	}

	res, err := nreflect.CallFunc(testFunc, "Alex", 1, []bool{false})
	if err != nil {
		t.Fatalf("Should have executed function")
	}
	t.Logf("Should have executed function")

	if !reflect.DeepEqual(res, []interface{}{[]interface{}{"Alex", 1, []bool{false}}}) {
		t.Logf("Expected: %q", []interface{}{[]interface{}{"Alex", 1, []bool{false}}})
		t.Logf("Received: %q", res)
		t.Fatalf("Expected value unmatched")
	}
}

func TestFunctionApply_ThreeArgument(t *testing.T) {
	var testFunc = func(v string, i int, d bool) string {
		return "Hello " + v
	}

	res, err := nreflect.CallFunc(testFunc, "Alex", 1, false)
	if err != nil {
		t.Fatalf("Should have executed function")
	}
	t.Logf("Should have executed function")

	if !reflect.DeepEqual(res, []interface{}{"Hello Alex"}) {
		t.Logf("Received: %q", res)
		t.Fatalf("Expected value unmatched")
	}
}

func TestMatchFunction(t *testing.T) {
	var addr1 = func(_ Addrs) error { return nil }
	var addr2 = func(_ Addrs) error { return nil }

	if !nreflect.MatchFunction(addr1, addr2) {
		t.Fatalf("Should have matched argument types successfully")
	}
	t.Logf("Should have matched argument types successfully")

	if !nreflect.MatchFunction(&addr1, &addr2) {
		t.Fatalf("Should have matched argument types successfully")
	}
	t.Logf("Should have matched argument types successfully")

	if nreflect.MatchFunction(&addr1, addr2) {
		t.Fatalf("Should have failed matched argument types successfully")
	}
	t.Logf("Should have failed matched argument types successfully")
}

func TestMatchElement(t *testing.T) {
	if !nreflect.MatchElement(Addrs{}, Addrs{}, false) {
		t.Fatalf("Should have matched argument types successfully")
	}
	t.Logf("Should have matched argument types successfully")

	if !nreflect.MatchElement(new(Addrs), new(Addrs), false) {
		t.Fatalf("Should have matched argument types successfully")
	}
	t.Logf("Should have matched argument types successfully")

	if nreflect.MatchElement(new(Addrs), Addrs{}, false) {
		t.Fatalf("Should have failed matched argument types successfully")
	}
	t.Logf("Should have failed matched argument types successfully")
}

func TestStructMapperWithSlice(t *testing.T) {
	mapper := nreflect.NewStructMapper()

	profile := struct {
		List []Addrs
	}{
		List: []Addrs{{Addr: "Tokura 20"}},
	}

	mapped, err := mapper.MapFrom("json", profile)
	if err != nil {
		t.Fatalf("Should have successfully converted struct")
	}
	t.Logf("Should have successfully converted struct")

	t.Logf("Map of Struct: %+q", mapped)

	profile2 := struct {
		List []Addrs
	}{}

	if err := mapper.MapTo("json", &profile2, mapped); err != nil {
		t.Fatalf("Should have successfully mapped data back to struct")
	}
	t.Logf("Should have successfully mapped data back to struct")

	if len(profile.List) != len(profile2.List) {
		t.Fatalf("Mapped struct should have same length: %d - %d ", len(profile.List), len(profile2.List))
	}
	t.Logf("Mapped struct should have same length: %d - %d ", len(profile.List), len(profile2.List))

	for ind, item := range profile.List {
		nxItem := profile2.List[ind]
		if item.Addr != nxItem.Addr {
			t.Fatalf("Item at %d should have equal value %+q -> %+q", ind, item.Addr, nxItem.Addr)
		}
	}

	t.Logf("All items should be exactly the same")
}

func TestStructMapperWthFieldStruct(t *testing.T) {
	layout := "Mon Jan 2 2006 15:04:05 -0700 MST"
	timeType := reflect.TypeOf((*time.Time)(nil))

	mapper := nreflect.NewStructMapper()
	mapper.AddAdapter(timeType, nreflect.TimeMapper(layout))
	mapper.AddInverseAdapter(timeType, nreflect.TimeInverseMapper(layout))

	profile := struct {
		Addr Addrs
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}{
		Addr: Addrs{Addr: "Tokura 20"},
		Name: "Johnson",
		Date: time.Now(),
	}

	mapped, err := mapper.MapFrom("json", profile)
	if err != nil {
		t.Fatalf("Should have successfully converted struct")
	}
	t.Logf("Should have successfully converted struct")

	t.Logf("Map of Struct: %+q", mapped)

	profile2 := struct {
		Addr Addrs
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}{}

	if err := mapper.MapTo("json", &profile2, mapped); err != nil {
		t.Fatalf("Should have successfully mapped data back to struct")
	}
	t.Logf("Should have successfully mapped data back to struct")

	if profile2.Addr.Addr != profile.Addr.Addr {
		t.Fatalf("Mapped struct should have same %q value", "Addr.Addr")
	}
	t.Logf("Mapped struct should have same %q value", "Addr.Addr")
}

func TestGetFieldByTagAndValue(t *testing.T) {
	profile := struct {
		Addrs
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}{
		Addrs: Addrs{Addr: "Tokura 20"},
		Name:  "Johnson",
		Date:  time.Now(),
	}

	_, err := nreflect.GetFieldByTagAndValue(profile, "json", "name")
	if err != nil {
		t.Fatalf("Should have successfully converted struct")
	}
}

func TestStructMapperWthEmbeddedStruct(t *testing.T) {
	layout := "Mon Jan 2 2006 15:04:05 -0700 MST"
	timeType := reflect.TypeOf((*time.Time)(nil))

	mapper := nreflect.NewStructMapper()
	mapper.AddAdapter(timeType, nreflect.TimeMapper(layout))
	mapper.AddInverseAdapter(timeType, nreflect.TimeInverseMapper(layout))

	profile := struct {
		Addrs
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}{
		Addrs: Addrs{Addr: "Tokura 20"},
		Name:  "Johnson",
		Date:  time.Now(),
	}

	mapped, err := mapper.MapFrom("json", profile)
	if err != nil {
		t.Fatalf("Should have successfully converted struct")
	}
	t.Logf("Should have successfully converted struct")

	t.Logf("Map of Struct: %+q", mapped)

	profile2 := struct {
		Addrs
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}{}

	if err := mapper.MapTo("json", &profile2, mapped); err != nil {
		t.Fatalf("Should have successfully mapped data back to struct")
	}
	t.Logf("Should have successfully mapped data back to struct")

	if profile2.Addr != profile.Addr {
		t.Fatalf("Mapped struct should have same %q value", "Addr.Addr")
	}
	t.Logf("Mapped struct should have same %q value", "Addr.Addr")
}

func TestStructMapper(t *testing.T) {
	layout := "Mon Jan 2 2006 15:04:05 -0700 MST"
	timeType := reflect.TypeOf((*time.Time)(nil))

	mapper := nreflect.NewStructMapper()
	mapper.AddAdapter(timeType, nreflect.TimeMapper(layout))
	mapper.AddInverseAdapter(timeType, nreflect.TimeInverseMapper(layout))

	profile := struct {
		Addr        string
		CountryName string
		Name        string    `json:"name"`
		Date        time.Time `json:"date"`
	}{
		Addr:        "Tokura 20",
		Name:        "Johnson",
		CountryName: "Nigeria",
		Date:        time.Now(),
	}

	mapped, err := mapper.MapFrom("json", profile)
	if err != nil {
		t.Fatalf("Should have successfully converted struct")
	}
	t.Logf("Should have successfully converted struct")

	t.Logf("Map of Struct: %+q", mapped)

	if _, ok := mapped["name"]; !ok {
		t.Fatalf("Map should have %q field", "name")
	}
	t.Logf("Map should have %q field", "name")

	if _, ok := mapped["date"]; !ok {
		t.Fatalf("Map should have %q field", "date")
	}
	t.Logf("Map should have %q field", "date")

	if _, ok := mapped["addr"]; !ok {
		t.Fatalf("Map should have %q field", "addr")
	}
	t.Logf("Map should have %q field", "addr")

	if _, ok := mapped["data"].(string); ok {
		t.Fatalf("Map should have field %q be a string", "date")
	}
	t.Logf("Map should have field %q be a string", "date")

	profile2 := struct {
		Addr        string
		CountryName string
		Name        string    `json:"name"`
		Date        time.Time `json:"date"`
	}{}

	if err := mapper.MapTo("json", &profile2, mapped); err != nil {
		t.Fatalf("Should have successfully mapped data back to struct")
	}
	t.Logf("Should have successfully mapped data back to struct")

	t.Logf("Mapped Struct: %+q", profile2)

	if profile2.Name != profile.Name {
		t.Fatalf("Mapped struct should have same %q value", "Name")
	}
	t.Logf("Mapped struct should have same %q value", "Name")

	if profile2.Date.Format(layout) != profile.Date.Format(layout) {
		t.Fatalf("Mapped struct should have same %q value", "Date")
	}
	t.Logf("Mapped struct should have same %q value", "Date")

	if profile2.CountryName != profile.CountryName {
		t.Fatalf("Mapped struct should have same %q value", "CountryName")
	}
	t.Logf("Mapped struct should have same %q value", "CountryName")

	if profile2.Addr != profile.Addr {
		t.Fatalf("Mapped struct should have same %q value", "Addr")
	}
	t.Logf("Mapped struct should have same %q value", "Addr")
}

// TestGetArgumentsType validates nreflect API GetArgumentsType functions
// results.
func TestGetArgumentsType(t *testing.T) {
	f := func(m monster) string {
		return fmt.Sprintf("Monster[%s] is ready!", m.Name)
	}

	args, err := nreflect.GetFuncArgumentsType(f)
	if err != nil {
		t.Fatalf("Should be able to retrieve function arguments lists")
	}
	t.Logf("Should be able to retrieve function arguments lists")

	name, embedded, err := nreflect.ExternalTypeNames(monster{Name: "Bob"})
	if err != nil {
		t.Fatalf("Should be able to retrieve field names arguments lists")
	}
	t.Logf("Name: %s", name)
	t.Logf("Fields: %+q", embedded)
	t.Logf("Should be able to retrieve function arguments lists")

	get(t, &monster{Name: "Bob"})

	newVals := nreflect.MakeArgumentsValues(args)
	if nlen, alen := len(newVals), len(args); nlen != alen {
		t.Fatalf("Should have matching new values lists for arguments")
	}
	t.Logf("Should have matching new values lists for arguments")

	mstring := reflect.TypeOf((*monster)(nil)).Elem()

	if mstring.Kind() != newVals[0].Kind() {
		t.Fatalf("Should be able to match argument kind")
	}
	t.Logf("Should be able to match argument kind")

}

func TestMatchFUncArgumentTypeWithValues(t *testing.T) {
	f := func(m monster) string {
		return fmt.Sprintf("Monster[%s] is ready!", m.Name)
	}

	var vals []reflect.Value
	vals = append(vals, reflect.ValueOf(monster{Name: "FireHouse"}))

	if index := nreflect.MatchFuncArgumentTypeWithValues(f, vals); index != -1 {
		t.Fatalf("Should have matching new values lists for arguments: %d", index)
	}
	t.Logf("Should have matching new values lists for arguments")
}
