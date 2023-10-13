package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}

}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postData := url.Values{
		"a": {"a"},
		"b": {"b"},
		"c": {"c"},
	}

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postData
	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	if form.Has("a") {
		t.Error("form shows has field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	if !form.Has("a") {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("x", 3)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData := url.Values{}
	postedData.Add("some_field", "abcd")
	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("form shows min lenght of 100 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("another_field", "abcd")
	form = New(postedData)

	form.MinLength("another_field", 3)
	if !form.Valid() {
		t.Error("form shows lack min length when it does")
	}

	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")

	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("wrong_email", "abcd")
	form = New(postedData)

	form.IsEmail("wrong_email")

	if form.Valid() {
		t.Error("form shows email is valid when it is not")
	}

	postedData = url.Values{}
	postedData.Add("correct_email", "abcd@gmail.com")
	form = New(postedData)

	form.IsEmail("correct_email")

	if !form.Valid() {
		t.Error("form shows email is not valid when it is")
	}
}
