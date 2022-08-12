package main

import "testing"

func assertCorrectMessage(t *testing.T, got string, want string) {
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestHello(t *testing.T) {
	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Chris", "English")
		want := "Hello, Chris"
		assertCorrectMessage(t, got, want)
	})
	t.Run("say 'Hello, World' when an empty string is supplied", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})
	t.Run("Say 'Hola, Chris' in Spanish", func(t *testing.T) {
		got := Hello("Chris", "Spanish")
		want := "Hola, Chris"
		assertCorrectMessage(t, got, want)
	})
	t.Run("Say hola everyone in Spanish", func(t *testing.T) {
		got := Hello("", "Spanish")
		want := "Hola, World"
		assertCorrectMessage(t, got, want)
	})
	t.Run("Say Bonjor to everyone French", func(t *testing.T) {
		got := Hello("", "French")
		want := "Bonjor, World"
		assertCorrectMessage(t, got, want)
	})
	t.Run("Say 'Bonjor, Chris'", func(t *testing.T) {
		got := Hello("Chris", "French")
		want := "Bonjor, Chris"
		assertCorrectMessage(t, got, want)
	})
}
