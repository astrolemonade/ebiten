// Copyright 2022 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ui

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Foundation -framework UIKit
//
// #import <UIKit/UIKit.h>
//
// static double devicePixelRatio() {
//   return [[UIScreen mainScreen] nativeScale];
// }
import "C"

import (
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/metal"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl"
)

type graphicsDriverCreatorImpl struct {
}

func (g *graphicsDriverCreatorImpl) newAuto() (graphicsdriver.Graphics, GraphicsLibrary, error) {
	m, err1 := g.newMetal()
	if err1 == nil {
		return m, GraphicsLibraryMetal, nil
	}
	o, err2 := g.newOpenGL()
	if err2 == nil {
		return o, GraphicsLibraryMetal, nil
	}
	return nil, GraphicsLibraryUnknown, fmt.Errorf("ui: failed to choose graphics drivers: Metal: %v, OpenGL: %v", err1, err2)
}

func (g *graphicsDriverCreatorImpl) newOpenGL() (graphicsdriver.Graphics, error) {
	return opengl.NewGraphics()
}

func (*graphicsDriverCreatorImpl) newDirectX() (graphicsdriver.Graphics, error) {
	return nil, errors.New("ui: DirectX is not supported in this environment")
}

func (g *graphicsDriverCreatorImpl) newMetal() (graphicsdriver.Graphics, error) {
	return metal.NewGraphics()
}

func (*graphicsDriverCreatorImpl) newPlayStation5() (graphicsdriver.Graphics, error) {
	return nil, errors.New("ui: PlayStation 5 is not supported in this environment")
}

func (u *UserInterface) SetUIView(uiview uintptr) error {
	select {
	case err := <-u.errCh:
		return err
	case <-u.graphicsLibraryInitCh:
	}

	// This function should be called only when the graphics library is Metal.
	if g, ok := u.graphicsDriver.(interface{ SetUIView(uintptr) }); ok {
		g.SetUIView(uiview)
	}
	return nil
}

func (u *UserInterface) IsGL() (bool, error) {
	select {
	case err := <-u.errCh:
		return false, err
	case <-u.graphicsLibraryInitCh:
	}

	return u.GraphicsLibrary() == GraphicsLibraryOpenGL, nil
}

func deviceScaleFactorImpl() float64 {
	// TODO: Can this be called from non-main threads?
	return float64(C.devicePixelRatio())
}
