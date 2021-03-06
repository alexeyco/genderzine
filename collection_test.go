package genderize_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/alexeyco/genderize"
)

// nolint:gochecknoglobals,golint,stylecheck
var testCollectionGenders = []*genderize.Gender{
	{
		Name:        "Alice",
		Gender:      "female",
		Probability: 0.9,
		Count:       12345,
	},
	{
		Name:        "John",
		Gender:      "male",
		Probability: 0.9,
		Count:       87890,
	},
}

type testCollectionRoundTripper struct {
	genders []*genderize.Gender
}

func (c *testCollectionRoundTripper) RoundTrip(_ *http.Request) (res *http.Response, err error) {
	b, err := json.Marshal(c.genders)
	if err != nil {
		return
	}

	h := http.Header{}
	h.Set(genderize.HdrXRateLimitLimit, "123")
	h.Set(genderize.HdrXRateLimitRemaining, "456")
	h.Set(genderize.HdrXRateReset, "789")

	res = &http.Response{
		StatusCode: http.StatusOK,
		Header:     h,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}

	return
}

func testCollectionClient(genders ...*genderize.Gender) *http.Client {
	return &http.Client{
		Transport: &testCollectionRoundTripper{
			genders: genders,
		},
	}
}

func TestCollection_Limit(t *testing.T) {
	httpClient := testCollectionClient(testCollectionGenders...)

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	if c.Limit() != 123 {
		t.Errorf(`Should be %d, %d given`, 123, c.Limit())
	}
}

func TestCollection_LimitRemaining(t *testing.T) {
	httpClient := testCollectionClient(testCollectionGenders...)

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	if c.LimitRemaining() != 456 {
		t.Errorf(`Should be %d, %d given`, 456, c.LimitRemaining())
	}
}

func TestCollection_LimitReset(t *testing.T) {
	httpClient := testCollectionClient(testCollectionGenders...)

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	if c.LimitReset() != 789*time.Second {
		t.Errorf(`Should be %d, %d given`, 789*time.Second, c.LimitReset())
	}
}

func TestCollection_Length(t *testing.T) {
	httpClient := testCollectionClient(testCollectionGenders...)

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	if c.Length() != 2 {
		t.Errorf(`Should be %d, %d given`, 2, c.Length())
	}
}

func TestCollection_Find(t *testing.T) {
	httpClient := testCollectionClient(testCollectionGenders...)

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	alice, err := c.Find("Alice")
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	if !reflect.DeepEqual(alice, testCollectionGenders[0]) {
		t.Error(`Should be equal`)
	}

	mike, err := c.Find("Mike")
	if err == nil {
		t.Error(`Should not be nil`)
	}

	if !errors.Is(err, genderize.ErrNothingFound) {
		t.Error(`Should be "genderize.ErrNothingFound"`)
	}

	if mike != nil {
		t.Error(`Should be nil`)
	}
}

func TestCollection_FindX(t *testing.T) {
	httpClient := testCollectionClient(testCollectionGenders...)

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	alice := c.FindX("Alice")

	if !reflect.DeepEqual(alice, testCollectionGenders[0]) {
		t.Error(`Should be equal`)
	}
}

func TestCollection_FindX_Panic(t *testing.T) {
	httpClient := testCollectionClient()

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Error(`Should not be nil`)
		}

		if !errors.Is(err.(error), genderize.ErrNothingFound) {
			t.Error(`Should not be genderize.ErrNothingFound`)
		}
	}()

	alice := c.FindX("Alice")
	if alice != nil {
		t.Error(`Should be nil`)
	}
}

func TestCollection_First(t *testing.T) {
	httpClient := testCollectionClient(testCollectionGenders...)

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	person, err := c.First()
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	if !reflect.DeepEqual(person, testCollectionGenders[0]) && !reflect.DeepEqual(person, testCollectionGenders[1]) {
		t.Error(`Should be one of two genders`)
	}
}

func TestCollection_FirstX(t *testing.T) {
	httpClient := testCollectionClient(testCollectionGenders...)

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	person := c.FirstX()

	if !reflect.DeepEqual(person, testCollectionGenders[0]) && !reflect.DeepEqual(person, testCollectionGenders[1]) {
		t.Error(`Should be one of two genders`)
	}
}

func TestCollection_FirstX_Panic(t *testing.T) {
	httpClient := testCollectionClient()

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Error(`Should not be nil`)
		}

		if !errors.Is(err.(error), genderize.ErrNothingFound) {
			t.Error(`Should not be genderize.ErrNothingFound`)
		}
	}()

	person := c.FirstX()
	if person != nil {
		t.Error(`Should be nil`)
	}
}

func TestCollection_Each(t *testing.T) {
	httpClient := testCollectionClient(testCollectionGenders...)

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	cnt := 0
	err = c.Each(func(g *genderize.Gender) {
		if !reflect.DeepEqual(g, testCollectionGenders[0]) && !reflect.DeepEqual(g, testCollectionGenders[1]) {
			t.Error(`Should be one of two genders`)
		}

		cnt++
	})

	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	if cnt != 2 {
		t.Errorf(`Should be %d, %d given`, 2, cnt)
	}
}

func TestCollection_EachX(t *testing.T) {
	httpClient := testCollectionClient(testCollectionGenders...)

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	defer func() {
		if err := recover(); err != nil {
			t.Errorf(`Should be nil, "%s" given`, err)
		}
	}()

	cnt := 0

	c.EachX(func(g *genderize.Gender) {
		if !reflect.DeepEqual(g, testCollectionGenders[0]) && !reflect.DeepEqual(g, testCollectionGenders[1]) {
			t.Error(`Should be one of two genders`)
		}

		cnt++
	})

	if cnt != 2 {
		t.Errorf(`Should be %d, %d given`, 2, cnt)
	}
}

func TestCollection_EachX_Panic(t *testing.T) {
	httpClient := testCollectionClient()

	r := genderize.NewRequest(context.TODO()).
		Name("Alice", "John")

	c, err := genderize.NewClient(genderize.WithHTTPClient(httpClient)).
		Execute(r)
	if err != nil {
		t.Errorf(`Should be nil, "%s" given`, err)
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Error(`Should not be nil`)
		}

		if !errors.Is(err.(error), genderize.ErrNothingFound) {
			t.Error(`Should not be genderize.ErrNothingFound`)
		}
	}()

	cnt := 0

	c.EachX(func(g *genderize.Gender) {
		if !reflect.DeepEqual(g, testCollectionGenders[0]) && !reflect.DeepEqual(g, testCollectionGenders[1]) {
			t.Error(`Should be one of two genders`)
		}

		cnt++
	})

	if cnt != 2 {
		t.Errorf(`Should be %d, %d given`, 2, cnt)
	}
}
