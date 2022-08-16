package main

import (
	"testing"
)

func TestHello(t *testing.T) {

	helloTest := []struct {
		name  string
		hello string
		want  string
	}{
		{name: "hello everyone default english", hello: Hello("", ""), want: "Hello, World"},
		{name: "hello chris english", hello: Hello("Chris", "English"), want: "Hello, Chris"},
		{name: "hello everyone english", hello: Hello("", "English"), want: "Hello, World"},
		{name: "hello chris spanish", hello: Hello("Chris", "Spanish"), want: "Hola, Chris"},
		{name: "hello everyone spanish", hello: Hello("", "Spanish"), want: "Hola, World"},
		{name: "hello chris french", hello: Hello("Chris", "French"), want: "Bonjor, Chris"},
		{name: "hello everyone french", hello: Hello("", "French"), want: "Bonjor, World"},
	}

	for _, tt := range helloTest {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hello
			if got != tt.want {
				t.Errorf("%#v got %q want %q", tt.hello, got, tt.want)
			}
		})
	}

}
