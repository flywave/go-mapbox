package sprite

func NewNode(name string, width, height int) *Node {
	return &Node{Key: name, Width: width, Height: height}
}

type Node struct {
	Key    string
	Width  int
	Height int
	X      int
	Y      int
	Used   bool
	Right  *Node
	Down   *Node
}

type GrowingPacker struct {
	root *Node
}

func (p *GrowingPacker) Fit(blocks []*Node) {
	length := len(blocks)

	width, height := 0, 0

	if length > 0 {
		width, height = blocks[0].Width, blocks[0].Height
	}

	p.root = &Node{X: 0, Y: 0, Width: width, Height: height}

	for i := 0; i < length; i++ {
		block := blocks[i]
		node := findNode(p.root, block.Width, block.Height)

		if node != nil {
			fit := splitNode(node, block.Width, block.Height)
			block.X = fit.X
			block.Y = fit.Y
			continue
		}

		fit := p.growNode(block.Width, block.Height)
		block.X = fit.X
		block.Y = fit.Y
	}
}

func findNode(root *Node, width, height int) *Node {
	if root.Used {
		n := findNode(root.Right, width, height)
		if n != nil {
			return n
		}
		return findNode(root.Down, width, height)
	}
	if (width <= root.Width) && (height <= root.Height) {
		return root
	}
	return nil
}

func splitNode(node *Node, width, height int) *Node {
	node.Used = true
	node.Down = &Node{X: node.X, Y: node.Y + height, Width: node.Width, Height: node.Height - height}
	node.Right = &Node{X: node.X + width, Y: node.Y, Width: node.Width - width, Height: height}
	return node
}

func (p *GrowingPacker) growNode(width, height int) *Node {
	canGrowDown := width <= p.root.Width
	canGrowRight := height <= p.root.Height
	shouldGrowRight := canGrowRight && (p.root.Height >= (p.root.Width + width))
	shouldGrowDown := canGrowDown && (p.root.Width >= (p.root.Height + height))

	if shouldGrowRight {
		return p.growRight(width, height)
	} else if shouldGrowDown {
		return p.growDown(width, height)
	} else if canGrowRight {
		return p.growRight(width, height)
	} else if canGrowDown {
		return p.growDown(width, height)
	}
	return nil
}

func (p *GrowingPacker) growDown(width, height int) *Node {
	p.root = &Node{
		Used:   true,
		X:      0,
		Y:      0,
		Width:  p.root.Width,
		Height: p.root.Height + height,
		Down:   &Node{X: 0, Y: p.root.Height, Width: p.root.Width, Height: height},
		Right:  p.root,
	}

	node := findNode(p.root, width, height)
	if node != nil {
		return splitNode(node, width, height)
	}
	return nil
}

func (p *GrowingPacker) growRight(width, height int) *Node {
	p.root = &Node{
		Used:   true,
		X:      0,
		Y:      0,
		Width:  p.root.Width + width,
		Height: p.root.Height,
		Right:  &Node{X: p.root.Width, Y: 0, Width: width, Height: p.root.Height},
		Down:   p.root,
	}

	node := findNode(p.root, width, height)
	if node != nil {
		return splitNode(node, width, height)
	}
	return nil
}
