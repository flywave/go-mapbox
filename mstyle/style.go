package mapboxglstyle

import (
	"encoding/json"
	"image/color"
	"io"

	"github.com/jamesrr39/goutil/errorsx"
	"github.com/jamesrr39/ownmap-app/ownmap"
	"github.com/jamesrr39/ownmap-app/ownmapdal"
	"github.com/jamesrr39/ownmap-app/styling"
)

const (
	ObjectTypePoint      = "Point"
	ObjectTypeLineString = "LineString"
	ObjectTypePolygon    = "Polygon"

	backgroundLayerID = "background"
)

func (s *MapboxGLStyle) GetStyleID() string {
	return s.style.ID
}

func (s *MapboxGLStyle) GetNodeStyle(node *ownmap.OSMNode, zoomLevel ownmap.ZoomLevel) (*styling.NodeStyle, errorsx.Error) {
	for layerIndex, layer := range s.style.Layers {
		if layer.ID == backgroundLayerID {
			continue
		}

		if layer.MinZoom != nil && *layer.MinZoom > float64(zoomLevel) {
			continue
		}

		if layer.MaxZoom != nil && *layer.MaxZoom < float64(zoomLevel) {
			continue
		}

		nodeStyle := layer.GetLayerNodeStyle(node, zoomLevel, layerIndex)
		if nodeStyle != nil {
			return nodeStyle, nil
		}
	}

	// not shown
	return nil, nil
}

func (s *MapboxGLStyle) GetWayStyle(tags []*ownmap.OSMTag, zoomLevel ownmap.ZoomLevel) (*styling.WayStyle, errorsx.Error) {
	// go through each layer

	for layerIndex, layer := range s.style.Layers {
		if layer.ID == backgroundLayerID {
			continue
		}

		if layer.MinZoom != nil && *layer.MinZoom > float64(zoomLevel) {
			continue
		}

		if layer.MaxZoom != nil && *layer.MaxZoom < float64(zoomLevel) {
			continue
		}

		wayStyle := layer.GetLayerWayStyle(tags, zoomLevel, layerIndex)
		if wayStyle != nil {
			return wayStyle, nil
		}
	}

	// not shown
	return nil, nil
}

func (s *MapboxGLStyle) GetRelationStyle(relationData *ownmap.RelationData, zoomLevel ownmap.ZoomLevel) (*styling.RelationStyle, errorsx.Error) {
	wayStyle, err := s.GetWayStyle(relationData.Tags, zoomLevel)
	if err != nil {
		return nil, err
	}

	if wayStyle == nil {
		return nil, nil
	}

	return &styling.RelationStyle{ZIndex: wayStyle.ZIndex}, nil
}

func (s *MapboxGLStyle) GetBackground() color.Color {
	return s.backgroundColor
}

func (s *MapboxGLStyle) GetWantedObjects(zoomLevel ownmap.ZoomLevel) []*ownmapdal.TagKeyWithType {
	var objects []*ownmapdal.TagKeyWithType
	for _, layer := range s.style.Layers {
		if layer.ID == backgroundLayerID {
			continue
		}

		if layer.MinZoom != nil && *layer.MinZoom > float64(zoomLevel) {
			continue
		}

		if layer.MaxZoom != nil && *layer.MaxZoom < float64(zoomLevel) {
			continue
		}

		tagKeysToFetch := layer.Filter.GetTagKeysToFetch(layer.SourceLayer)
		objects = append(objects, tagKeysToFetch...)
	}
	return objects
}

func (s *styleType) calculateBackgroundColor() (color.Color, errorsx.Error) {
	if len(s.Layers) == 0 {
		return color.White, nil
	}

	backgroundLayer := s.Layers[0]
	if backgroundLayer.ID != backgroundLayerID {
		return color.White, nil
	}

	return backgroundLayer.Paint.BackgroundColor.GetColorAtZoomLevel(0), nil
}

// MapboxGLStyle represents a stylesheet represented in the MapboxGL format.
// https://docs.mapbox.com/mapbox-gl-js/style-spec/
// https://maputnik.github.io/editor
type MapboxGLStyle struct {
	style           *styleType
	backgroundColor color.Color
}

type styleType struct {
	Bearing    float64    `json:"bearing"`
	Center     []float64  `json:"center"`
	Glyphs     string     `json:"glyphs"`
	Layers     []*Layer   `json:"layers"`
	Light      Light      `json:"light"`
	Metadata   Metadata   `json:"metadata"`
	Name       string     `json:"name"`
	Pitch      float64    `json:"pitch"`
	Sources    Sources    `json:"sources"`
	Sprite     string     `json:"sprite"`
	Transition Transition `json:"transition"`
	Version    int        `json:"version"`
	Zoom       float64    `json:"zoom"`
	ID         string     `json:"id"`
}

func Parse(reader io.Reader) (*MapboxGLStyle, errorsx.Error) {
	s := new(styleType)
	dec := json.NewDecoder(reader)
	// dec.DisallowUnknownFields()
	err := dec.Decode(s)
	if err != nil {
		return nil, errorsx.Wrap(err)
	}

	bgColor, err := s.calculateBackgroundColor()
	if err != nil {
		return nil, errorsx.Wrap(err)
	}

	return &MapboxGLStyle{s, bgColor}, nil
}

func (s *styleType) Validate() errorsx.Error {
	const expectedVersion = 8
	if s.Version != expectedVersion {
		return errorsx.Errorf("version: expected %d but was %d", expectedVersion, s.Version)
	}

	return nil
}
