package progress

import "testing"

func Test_progress(t *testing.T) {
	t.Run("should have a progress of 0 with defined values", func(t *testing.T) {
		p := &Progress{
			Max:     100,
			Current: 0,
			Steps:   1,
		}

		if p.Get() != 0 {
			t.Errorf("expected %d, got %d", 0, p.Get())
		}
	})

	t.Run("should have a progress of 0 with default values", func(t *testing.T) {
		p := &Progress{
			Max: 100,
		}

		if p.Get() != 0 {
			t.Errorf("expected %d, got %d", 0, p.Get())
		}
	})

	t.Run("should have a progress of 50", func(t *testing.T) {
		p := &Progress{
			Max:     100,
			Current: 50,
		}

		if p.Get() != 50 {
			t.Errorf("expected %d, got %d", 50, p.Get())
		}
	})

	t.Run("should have a progress of 25", func(t *testing.T) {
		p := &Progress{
			Max:     200,
			Current: 50,
		}

		if p.Get() != 25 {
			t.Errorf("expected %d, got %d", 25, p.Get())
		}
	})

	t.Run("should have a progress of 100", func(t *testing.T) {
		p := &Progress{
			Max:     100,
			Current: 100,
		}

		if p.Get() != 100 {
			t.Errorf("expected %d, got %d", 100, p.Get())
		}
	})

	t.Run("should have a progress of 100 with default steps", func(t *testing.T) {
		p := &Progress{
			Max:     100,
			Current: 0,
		}

		p.Add(100)

		if p.Get() != 100 {
			t.Errorf("expected %d, got %d", 100, p.Get())
		}
	})

	t.Run("should have a progress of 50 with steps", func(t *testing.T) {
		p := &Progress{
			Max:     100,
			Current: 0,
			Steps:   2,
		}

		p.Add(50)

		if p.Get() != 100 {
			t.Errorf("expected %d, got %d", 100, p.Get())
		}
	})

	t.Run("should have a progress of 50 with steps and max 60", func(t *testing.T) {
		p := &Progress{
			Max:     60,
			Current: 15,
			Steps:   5,
		}

		p.Add(3)

		if p.Get() != 50 {
			t.Errorf("expected %d, got %d", 50, p.Get())
		}
	})

	t.Run("should have a progress of 75", func(t *testing.T) {
		p := &Progress{
			Max:     60,
			Current: 15,
			Steps:   5,
		}

		p.Set(45)

		if p.Get() != 75 {
			t.Errorf("expected %d, got %d", 75, p.Get())
		}
	})
}
