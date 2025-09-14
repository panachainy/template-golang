package integration

import (
	"testing"
)

func TestExample(t *testing.T) {
	// Arrange
	expected := "hello"

	// Act
	actual := "hello"

	// Assert
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
