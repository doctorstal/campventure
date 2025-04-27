package animations

type Animation interface {
	// returns true if finished
	Update() (remove bool)
	Frame() int
}

type BasicAnimation struct {
	first        int
	last         int
	step         int     // how many indices do we move per frame
	speedInTps   float32 // how many ticks before next frame
	frameCounter float32
	frame        int
}

type LoopAnimation struct {
	*BasicAnimation
}

func (a *LoopAnimation) Update() (remove bool) {
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.speedInTps
		a.frame += a.step
		if a.frame > a.last {
			a.frame = a.first
		}
	}
	return false
}

func (a *LoopAnimation) Frame() int {
	return a.frame
}

func NewLoopAnimation(first, last, step int, speed float32) Animation {
	return &LoopAnimation{&BasicAnimation{
		first,
		last,
		step,
		speed,
		speed,
		first,
	}}
}

type SingleFrameAnimation struct {
	frame int
}

// Frame implements Animation.
func (s *SingleFrameAnimation) Frame() int {
	return s.frame
}

// Update implements Animation.
func (s *SingleFrameAnimation) Update() bool {
	return false
}

func NewSingleFrameAnimation(frame int) Animation {
	return &SingleFrameAnimation{
		frame: frame,
	}
}

type OneTimeAnimation struct {
	*BasicAnimation
	stopped     bool
	removeAfter bool
}

func (a *OneTimeAnimation) Update() (remove bool) {
	if a.stopped {
		return a.removeAfter
	}
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.speedInTps
		a.frame += a.step
		if a.frame >= a.last {
			a.stopped = true
		}
	}
	return false
}

func (a *OneTimeAnimation) Frame() int {
	return a.frame
}

func NewOneTimeAnimation(first, last, step int, speed float32, removeAfter bool) Animation {
	return &OneTimeAnimation{
		&BasicAnimation{
			first,
			last,
			step,
			speed,
			speed,
			first,
		},
		false,
		removeAfter,
	}
}

type CallBackAnimation struct {
	*LoopAnimation
	callback func(frame int) bool
}

func NewCallBackAnimation(first, last, step int, speed float32, callback func(frame int) bool) Animation {
	return &CallBackAnimation{
		&LoopAnimation{&BasicAnimation{
			first,
			last,
			step,
			speed,
			speed,
			first,
		}},
		callback,
	}
}

func (c *CallBackAnimation) Update() bool {
	c.frameCounter -= 1.0
	if c.frameCounter < 0.0 {
		c.frameCounter = c.speedInTps
		c.frame += c.step
		if c.frame >= c.last {
			c.frame = c.first
		}
		return c.callback(c.frame)
	}
	return false
}
