package http_runner

import (
	"reflect"
	"testing"
	"time"

	NetworkRunner "github.com/Tanreon/go-network-runner"
)

func TestDirectHttpGetJson(t *testing.T) {
	t.Run("TestDirectHttpGetJson", func(t *testing.T) {
		dialOptions := NetworkRunner.DialOptions{
			DialTimeout:  120,
			RelayTimeout: 60,
		}
		directDialer, err := NetworkRunner.NewDirectDialer(dialOptions)
		if err != nil {
			t.Fatal(err)
		}

		directHttpRunner, err := NewDirectHttpRunner(directDialer)
		if err != nil {
			t.Fatal(err)
		}

		jsonRequestData := NewJsonRequestData("https://httpbin.org/get")
		jsonRequestData.SetHeaders(map[string]string{
			"x-test": "true",
		})
		jsonRequestData.SetRetryOption(3)
		jsonRequestData.SetTimeoutOption(time.Second * 60)
		jsonRequestData.SetFollowRedirectOption(true)

		response, err := directHttpRunner.GetJson(jsonRequestData)
		if err != nil {
			t.Fatal(err)
		}
		if got := response.StatusCode(); got != 200 {
			t.Fatalf("response.StatusCode() = %v, want 200", got)
		}
	})
}
func TestDirectHttpPostJson(t *testing.T) {
	t.Run("TestDirectHttpPostJson", func(t *testing.T) {
		dialOptions := NetworkRunner.DialOptions{
			DialTimeout:  120,
			RelayTimeout: 60,
		}
		directDialer, err := NetworkRunner.NewDirectDialer(dialOptions)
		if err != nil {
			t.Fatal(err)
		}

		directHttpRunner, err := NewDirectHttpRunner(directDialer)
		if err != nil {
			t.Fatal(err)
		}

		jsonRequestData := NewJsonRequestData("https://httpbin.org/post")
		jsonRequestData.SetHeaders(map[string]string{
			"x-test": "true",
		})
		jsonRequestData.SetValue([]byte("test"))
		jsonRequestData.SetRetryOption(3)
		jsonRequestData.SetTimeoutOption(time.Second * 60)
		jsonRequestData.SetFollowRedirectOption(true)

		response, err := directHttpRunner.PostJson(jsonRequestData)
		if err != nil {
			t.Fatal(err)
		}
		if got := response.StatusCode(); got != 200 {
			t.Fatalf("response.StatusCode() = %v, want 200", got)
		}
	})
}

func TestJsonRequestData(t *testing.T) {
	t.Run("TestJsonRequestData-Url", func(t *testing.T) {
		want := "https://httpbin.org/get"

		jsonRequestData := NewJsonRequestData(want)

		if got := jsonRequestData.Url(); got != want {
			t.Errorf("jsonRequestData.Url() = %v, want %v", got, want)
		}
	})
	t.Run("TestJsonRequestData-Value", func(t *testing.T) {
		want := []byte("value")

		jsonRequestData := NewJsonRequestData("https://httpbin.org/get")
		jsonRequestData.SetValue(want)

		if got := jsonRequestData.IsValueSet(); got != true {
			t.Errorf("jsonRequestData.IsValueSet() = %v, want %v", got, true)
		}
		if got := jsonRequestData.Value(); !reflect.DeepEqual(got, want) {
			t.Errorf("jsonRequestData.Value() = %v, want %v", got, want)
		}
	})
	t.Run("TestJsonRequestData-Headers", func(t *testing.T) {
		want := map[string]string{
			"x-test-header": "test",
		}

		jsonRequestData := NewJsonRequestData("https://httpbin.org/get")
		jsonRequestData.SetHeaders(want)

		if got := jsonRequestData.IsHeadersSet(); got != true {
			t.Errorf("jsonRequestData.IsHeadersSet() = %v, want %v", got, true)
		}
		if got := jsonRequestData.Headers(); !reflect.DeepEqual(got, want) {
			t.Errorf("jsonRequestData.Headers() = %v, want %v", got, want)
		}
	})
	t.Run("TestJsonRequestData-RetryOption", func(t *testing.T) {
		want := 3

		jsonRequestData := NewJsonRequestData("https://httpbin.org/get")
		jsonRequestData.SetRetryOption(want)

		if got := jsonRequestData.IsRetryOptionSet(); got != true {
			t.Errorf("jsonRequestData.IsRetryOptionSet() = %v, want %v", got, true)
		}
		if got := jsonRequestData.RetryOption(); got != want {
			t.Errorf("jsonRequestData.RetryOption() = %v, want %v", got, want)
		}
	})
	t.Run("TestJsonRequestData-FollowRedirectOption", func(t *testing.T) {
		want := true

		jsonRequestData := NewJsonRequestData("https://httpbin.org/get")
		jsonRequestData.SetFollowRedirectOption(want)

		if got := jsonRequestData.IsFollowRedirectOptionSet(); got != true {
			t.Errorf("jsonRequestData.IsFollowRedirectOptionSet() = %v, want %v", got, true)
		}
		if got := jsonRequestData.FollowRedirectOption(); got != want {
			t.Errorf("jsonRequestData.FollowRedirectOption() = %v, want %v", got, want)
		}
	})
	t.Run("TestJsonRequestData-TimeoutOption", func(t *testing.T) {
		want := time.Second * 5

		jsonRequestData := NewJsonRequestData("https://httpbin.org/get")
		jsonRequestData.SetTimeoutOption(want)

		if got := jsonRequestData.IsTimeoutOptionSet(); got != true {
			t.Errorf("jsonRequestData.IsTimeoutOptionSet() = %v, want %v", got, true)
		}
		if got := jsonRequestData.TimeoutOption(); got != want {
			t.Errorf("jsonRequestData.TimeoutOption() = %v, want %v", got, want)
		}
	})
}
