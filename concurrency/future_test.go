package concurrency

import (
	"errors"
	"testing"
	"time"
)

func TestFuture_Get(t *testing.T) {
	tests := []struct {
		name        string
		action      func() (int, error)
		expectedVal int
		expectedErr error
	}{
		{
			name: "Success case",
			action: func() (int, error) {
				return 42, nil
			},
			expectedVal: 42,
			expectedErr: nil,
		},
		{
			name: "Error case",
			action: func() (int, error) {
				return 0, errors.New("test error")
			},
			expectedVal: 0,
			expectedErr: errors.New("test error"),
		},
		{
			name: "Delayed execution",
			action: func() (int, error) {
				time.Sleep(100 * time.Millisecond)
				return 100, nil
			},
			expectedVal: 100,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			future := NewFuture(tt.action)
			val, err := future.Get()

			if val != tt.expectedVal {
				t.Errorf("Expected value %v, got %v", tt.expectedVal, val)
			}

			if (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr == nil) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			} else if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error message %v, got %v", tt.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestFuture_Concurrency(t *testing.T) {
	// Test that multiple futures can run concurrently
	start := time.Now()

	future1 := NewFuture(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 1, nil
	})

	future2 := NewFuture(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 2, nil
	})

	val1, err1 := future1.Get()
	val2, err2 := future2.Get()

	elapsed := time.Since(start)

	if err1 != nil || err2 != nil {
		t.Errorf("Unexpected errors: %v, %v", err1, err2)
	}

	if val1 != 1 || val2 != 2 {
		t.Errorf("Expected values 1 and 2, got %v and %v", val1, val2)
	}

	// Both futures should complete in roughly 200ms, not 400ms
	if elapsed >= 390*time.Millisecond {
		t.Errorf("Futures did not run concurrently. Took %v", elapsed)
	}
}
