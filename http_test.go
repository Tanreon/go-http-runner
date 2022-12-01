package http_runner

import (
	"os"
	"reflect"
	"testing"
	"time"

	NetworkRunner "github.com/Tanreon/go-network-runner"
)

func TestDirectHttpGetJson(t *testing.T) {
	t.Run("TestDirectHttpGetJson", func(t *testing.T) {
		directDialOptions := NetworkRunner.NewDirectDialOptions()
		directDialOptions.SetDialTimeout(60)
		directDialOptions.SetRelayTimeout(60)
		directDialer, err := NetworkRunner.NewDirectDialer(directDialOptions)
		if err != nil {
			t.Fatal(err)
		}

		directHttpRunner, err := NewDirectHttpRunner(directDialer)
		if err != nil {
			t.Fatal(err)
		}

		jsonRequest := NewJsonRequestOptions("https://httpbin.org/get")
		jsonRequest.SetHeaders(map[string]string{
			"x-test": "true",
		})
		jsonRequest.SetRetryOption(3)
		jsonRequest.SetTimeoutOption(time.Second * 60)
		jsonRequest.SetFollowRedirectOption(true)

		response, err := directHttpRunner.GetJson(jsonRequest)
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
		directDialOptions := NetworkRunner.NewDirectDialOptions()
		directDialOptions.SetDialTimeout(60)
		directDialOptions.SetRelayTimeout(60)
		directDialer, err := NetworkRunner.NewDirectDialer(directDialOptions)
		if err != nil {
			t.Fatal(err)
		}

		directHttpRunner, err := NewDirectHttpRunner(directDialer)
		if err != nil {
			t.Fatal(err)
		}

		jsonRequest := NewJsonRequestOptions("https://httpbin.org/post")
		jsonRequest.SetHeaders(map[string]string{
			"x-test": "true",
		})
		jsonRequest.SetValue([]byte("test"))
		jsonRequest.SetRetryOption(3)
		jsonRequest.SetTimeoutOption(time.Second * 60)
		jsonRequest.SetFollowRedirectOption(true)

		response, err := directHttpRunner.PostJson(jsonRequest)
		if err != nil {
			t.Fatal(err)
		}
		if got := response.StatusCode(); got != 200 {
			t.Fatalf("response.StatusCode() = %v, want 200", got)
		}
	})
}

func TestJsonRequestOptions(t *testing.T) {
	t.Run("TestJsonRequest-Url", func(t *testing.T) {
		want := "https://httpbin.org/get"

		jsonRequest := NewJsonRequestOptions(want)

		if got := jsonRequest.Url(); got != want {
			t.Errorf("jsonRequest.Url() = %v, want %v", got, want)
		}
	})
	t.Run("TestJsonRequest-Value", func(t *testing.T) {
		want := []byte("value")

		jsonRequest := NewJsonRequestOptions("https://httpbin.org/get")
		jsonRequest.SetValue(want)

		if got := jsonRequest.IsValueSet(); got != true {
			t.Errorf("jsonRequest.IsValueSet() = %v, want %v", got, true)
		}
		if got := jsonRequest.Value(); !reflect.DeepEqual(got, want) {
			t.Errorf("jsonRequest.Value() = %v, want %v", got, want)
		}
	})
	t.Run("TestJsonRequest-Headers", func(t *testing.T) {
		want := map[string]string{
			"x-test-header": "test",
		}

		jsonRequest := NewJsonRequestOptions("https://httpbin.org/get")
		jsonRequest.SetHeaders(want)

		if got := jsonRequest.IsHeadersSet(); got != true {
			t.Errorf("jsonRequest.IsHeadersSet() = %v, want %v", got, true)
		}
		if got := jsonRequest.Headers(); !reflect.DeepEqual(got, want) {
			t.Errorf("jsonRequest.Headers() = %v, want %v", got, want)
		}
	})
	t.Run("TestJsonRequest-RetryOption", func(t *testing.T) {
		want := 3

		jsonRequest := NewJsonRequestOptions("https://httpbin.org/get")
		jsonRequest.SetRetryOption(want)

		if got := jsonRequest.IsRetryOptionSet(); got != true {
			t.Errorf("jsonRequest.IsRetryOptionSet() = %v, want %v", got, true)
		}
		if got := jsonRequest.RetryOption(); got != want {
			t.Errorf("jsonRequest.RetryOption() = %v, want %v", got, want)
		}
	})
	t.Run("TestJsonRequestOptions-FollowRedirectOption", func(t *testing.T) {
		want := true

		jsonRequest := NewJsonRequestOptions("https://httpbin.org/get")
		jsonRequest.SetFollowRedirectOption(want)

		if got := jsonRequest.IsFollowRedirectOptionSet(); got != true {
			t.Errorf("jsonRequest.IsFollowRedirectOptionSet() = %v, want %v", got, true)
		}
		if got := jsonRequest.FollowRedirectOption(); got != want {
			t.Errorf("jsonRequest.FollowRedirectOption() = %v, want %v", got, want)
		}
	})
	t.Run("TestJsonRequestOptions-TimeoutOption", func(t *testing.T) {
		want := time.Second * 5

		jsonRequest := NewJsonRequestOptions("https://httpbin.org/get")
		jsonRequest.SetTimeoutOption(want)

		if got := jsonRequest.IsTimeoutOptionSet(); got != true {
			t.Errorf("jsonRequest.IsTimeoutOptionSet() = %v, want %v", got, true)
		}
		if got := jsonRequest.TimeoutOption(); got != want {
			t.Errorf("jsonRequest.TimeoutOption() = %v, want %v", got, want)
		}
	})
}

func TestDirectHttpPostFiles(t *testing.T) {
	t.Run("TestDirectHttpPostFiles", func(t *testing.T) {
		directDialOptions := NetworkRunner.NewDirectDialOptions()
		directDialOptions.SetDialTimeout(60)
		directDialOptions.SetRelayTimeout(60)
		directDialer, err := NetworkRunner.NewDirectDialer(directDialOptions)
		if err != nil {
			t.Fatal(err)
		}

		directHttpRunner, err := NewDirectHttpRunner(directDialer)
		if err != nil {
			t.Fatal(err)
		}

		txtFile, err := os.Open("file_test.bin")
		if err != nil {
			t.Fatal(err)
		}
		defer txtFile.Close()

		files := map[string]FileInfo{
			"upload_files": BuildFileInfo(txtFile),
		}

		formRequest := NewFormRequestOptions("https://httpbin.org/post")
		formRequest.SetFiles(files)

		response, err := directHttpRunner.PostForm(formRequest)
		if err != nil {
			t.Fatal(err)
		}
		if got := response.StatusCode(); got != 200 {
			t.Fatalf("response.StatusCode() = %v, want 200", got)
		}
	})
}
