package style

import (
	"image/color"

	"github.com/pkg/errors"
)

type ItemStyle interface {
	GetZIndex() int
}

type WayStyle struct {
	FillColor      color.Color
	LineColor      color.Color
	LineDashPolicy []float64
	LineWidth      float64
	ZIndex         int
}

func (ws *WayStyle) GetZIndex() int {
	return ws.ZIndex
}

type NodeStyle struct {
	TextSize  int
	TextColor color.Color
	ZIndex    int
}

func (ns *NodeStyle) GetZIndex() int {
	return ns.ZIndex
}

type RelationStyle struct {
	ZIndex int
}

func (rs *RelationStyle) GetZIndex() int {
	return rs.ZIndex
}

type StyleSet struct {
	stylesMap      map[string]Style
	defaultStyleID string
}

func NewStyleSet(styles []Style, defaultStyleID string) (*StyleSet, error) {
	styleSet := &StyleSet{
		stylesMap:      make(map[string]Style),
		defaultStyleID: defaultStyleID,
	}

	defaultIDFound := false

	for _, style := range styles {
		styleID := style.ID
		_, ok := styleSet.stylesMap[styleID]
		if ok {
			return nil, errors.Errorf("duplicate style ID found: %q", styleID)
		}

		styleSet.stylesMap[styleID] = style

		if defaultStyleID == styleID {
			defaultIDFound = true
		}
	}

	if !defaultIDFound {
		return nil, errors.Errorf("default ID %q not found in any supplied styles", defaultStyleID)
	}

	return styleSet, nil
}

func (s *StyleSet) GetStyleByID(id string) Style {
	return s.stylesMap[id]
}

func (s *StyleSet) GetDefaultStyle() Style {
	return s.stylesMap[s.defaultStyleID]
}

func (s *StyleSet) GetAllStyleIDs() []string {
	var styleIDs []string

	for id := range s.stylesMap {
		styleIDs = append(styleIDs, id)
	}

	return styleIDs
}
