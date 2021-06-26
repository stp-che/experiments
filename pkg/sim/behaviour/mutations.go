package behaviour

import (
	"fmt"
	"math/rand"
)

type iMutation interface {
	Apply(*Brain) *Brain
}

type simpleMutation map[int]uint8

func (m simpleMutation) Apply(b *Brain) *Brain {
	fmt.Printf("%v\n", b.Content)
	newBrain := copyBrain(b)
	fmt.Printf("%v\n", newBrain.Content)
	for i, v := range m {
		if i >= 0 && i < len(newBrain.Content) {
			newBrain.Content[i] = v
		}
	}
	newBrain.NormalizeContent()
	return newBrain
}

func randomMutation(changes, contentSize int) iMutation {
	return randomSimpleMutation(changes, contentSize)
}

func randomSimpleMutation(changes, contentSize int) iMutation {
	m := simpleMutation{}
	for i := 0; i < changes; i++ {
		m[rand.Intn(contentSize)] = uint8(rand.Intn(256))
	}
	return m
}

func copyBrain(b *Brain) *Brain {
	newBrain := &Brain{
		Structure: b.Structure,
		Content:   make([]uint8, len(b.Content)),
	}
	copy(newBrain.Content, b.Content)
	return newBrain
}
