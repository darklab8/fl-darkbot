package records

import "testing"

func TestStamped(t *testing.T) {
	objects := StampedObjects[int]{}
	objects.Add(1)
}
