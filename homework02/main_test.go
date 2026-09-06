package main

import (
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want int
	}{
		{name: "positive", in: 5, want: 15},
		{name: "zero", in: 0, want: 10},
		{name: "negative", in: -3, want: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in
			add(&got)
			if got != tt.want {
				t.Fatalf("add(%d) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestDouble(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{name: "multiple values", in: []int{2, 7, 10}, want: []int{4, 14, 20}},
		{name: "zero and negative", in: []int{0, -3, 1}, want: []int{0, -6, 2}},
		{name: "empty", in: []int{}, want: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := make([]*int, len(tt.in))
			for i := range tt.in {
				v := tt.in[i]
				values[i] = &v
			}

			double(values)

			for i := range values {
				if *values[i] != tt.want[i] {
					t.Fatalf("double(%v) = %v, want %v", tt.in, extractValues(values), tt.want)
				}
			}
		})
	}
}

func extractValues(values []*int) []int {
	result := make([]int, len(values))
	for i, v := range values {
		result[i] = *v
	}
	return result
}

func TestSchedulerSubmitAfterShutdownDoesNotLeak(t *testing.T) {
	scheduler := NewScheduler(1, 2)
	scheduler.Shutdown()

	err := scheduler.Submit(func(ctx context.Context) (int, error) {
		return 1, nil
	})
	if err == nil {
		t.Fatal("Submit after Shutdown() should return an error")
	}

	done := make(chan struct{})
	go func() {
		scheduler.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Wait() blocked after failed Submit after Shutdown()")
	}
}

func TestSchedulerSubmitWhenContextCancelledDoesNotLeak(t *testing.T) {
	scheduler := NewScheduler(1, 1)
	scheduler.cancel()

	err := scheduler.Submit(func(ctx context.Context) (int, error) {
		return 1, nil
	})
	if err == nil {
		t.Fatal("Submit when ctx is already cancelled should return an error")
	}

	done := make(chan struct{})
	go func() {
		scheduler.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Wait() blocked after context cancellation before Submit")
	}
}

func TestAtomicCase(t *testing.T) {
	output := captureOutput(AtomicCase)

	if output != "10000\n" {
		t.Fatalf("AtomicCase() output = %q, want %q", output, "10000\n")
	}
}

func TestMutexProtectAdd(t *testing.T) {
	output := captureOutput(MutexProtectAdd)

	if output != "10000\n" {
		t.Fatalf("MutexProtectAdd() output = %q, want %q", output, "10000\n")
	}
}

func TestGoroutineCommunication(t *testing.T) {
	output := captureOutput(GoroutineCommunication)
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 100 {
		t.Fatalf("GoroutineCommunication() printed %d values, want 100", len(lines))
	}
	for i, line := range lines {
		want := "Received: " + strconv.Itoa(i)
		if line != want {
			t.Fatalf("GoroutineCommunication() line %d = %q, want %q", i, line, want)
		}
	}
}

func TestGoroutineCommunicationNoBuffer(t *testing.T) {
	output := captureOutput(GoroutineCommunicationNoBuffer)
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 10 {
		t.Fatalf("GoroutineCommunicationNoBuffer() printed %d values, want 10", len(lines))
	}
	for i, line := range lines {
		want := "Received: " + strconv.Itoa(i)
		if line != want {
			t.Fatalf("GoroutineCommunicationNoBuffer() line %d = %q, want %q", i, line, want)
		}
	}
}

func captureOutput(function func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	original := os.Stdout
	os.Stdout = writer
	function()
	writer.Close()
	os.Stdout = original

	output, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	reader.Close()
	return string(output)
}
