package shadow2d

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var input *inputHandler

func init() {
	input = newInputHandler()
}

type inputHandler struct {
	mouse *mouse
}

func newInputHandler() *inputHandler {
	return &inputHandler{mouse: newMouse()}
}

func (i *inputHandler) update() {
	i.mouse.update()
}

type mouse struct {
	x, y int

	isLeftButtonPressed, isRightButtonPressed   bool
	isLeftButtonReleased, isRightButtonReleased bool
}

func newMouse() *mouse {
	return &mouse{
		x:                     0,
		y:                     0,
		isLeftButtonPressed:   false,
		isRightButtonPressed:  false,
		isLeftButtonReleased:  false,
		isRightButtonReleased: false,
	}
}

func (m *mouse) update() {
	m.x, m.y = ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !m.isLeftButtonPressed {
		// log.Println(fmt.Sprintf("Right button pressed on x:[%d] y:[%d]", m.x, m.y))
		m.isLeftButtonPressed = true
		return
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && !m.isRightButtonPressed {
		m.isRightButtonPressed = true
		// log.Println(fmt.Sprintf("Right button pressed on x:[%d] y:[%d]", m.x, m.y))
		return
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && m.isLeftButtonPressed {
		m.isLeftButtonPressed = false
		m.isLeftButtonReleased = true
		// log.Println(fmt.Sprintf("Left button released on x:[%d] y:[%d]", m.x, m.y))
		return
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && m.isRightButtonPressed {
		m.isRightButtonPressed = false
		m.isRightButtonReleased = true
		// log.Println(fmt.Sprintf("Right button released on x:[%d] y:[%d]", m.x, m.y))
		return
	}

	if m.isLeftButtonReleased {
		m.isLeftButtonReleased = false
		return
	}

	if m.isRightButtonReleased {
		m.isRightButtonReleased = false
		return
	}
}

func (m *mouse) getCoordinates() (x, y int) {
	return m.x, m.y
}
