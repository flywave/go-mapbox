package mapboxglstyle

import (
	"image/color"
	"log"

	"github.com/jamesrr39/goutil/errorsx"
	"github.com/jamesrr39/ownmap-app/ownmap"
	"github.com/jamesrr39/ownmap-app/styling"
)

type LayerType string

const (
	LayerTypeBackground    LayerType = "background"
	LayerTypeFill          LayerType = "fill"
	LayerTypeLine          LayerType = "line"
	LayerTypeSymbol        LayerType = "symbol"
	LayerTypeRaster        LayerType = "raster"
	LayerTypeCircle        LayerType = "circle"
	LayerTypeFillExtrusion LayerType = "fill-extrusion"
	LayerTypeHeatmap       LayerType = "heatmap"
	LayerTypeHillshade     LayerType = "hillshade"
)

type Layer struct {
	Filter      FilterContainer `json:"filter"`
	ID          LayerType       `json:"id"`
	Layout      Layout          `json:"layout"`
	MaxZoom     *float64        `json:"maxzoom"`
	Metadata    Metadata        `json:"metadata"`
	MinZoom     *float64        `json:"minzoom"`
	Paint       *Paint          `json:"paint"`
	Source      string          `json:"source"`
	SourceLayer string          `json:"source-layer"`
	Type        LayerType       `json:"type"`
}

func (l *Layer) Validate() errorsx.Error {
	if l.MaxZoom != nil && l.MinZoom != nil {
		if *l.MaxZoom < *l.MinZoom {
			return errorsx.Errorf("max zoom is smaller than min zoom")
		}
	}

	if l.MaxZoom != nil && *l.MaxZoom < 0 || *l.MaxZoom > 24 {
		return errorsx.Errorf("max zoom must be between 0 and 24 (inclusive) but was %f", *l.MaxZoom)
	}

	if l.MinZoom != nil && *l.MinZoom < 0 || *l.MinZoom > 24 {
		return errorsx.Errorf("min zoom must be between 0 and 24 (inclusive) but was %f", *l.MinZoom)
	}

	return nil
}

func (l *Layer) GetLayerNodeStyle(node *ownmap.OSMNode, zoomLevel ownmap.ZoomLevel, layerIndex int) *styling.NodeStyle {
	switch l.Type {
	case "symbol":
		switch l.SourceLayer {
		case SourceLayerPlace:

			shown := l.Filter.IsObjectShown(l.SourceLayer, node.Tags, ownmap.ObjectTypeNode)
			if !shown {
				return nil
			}

			return &styling.NodeStyle{
				TextSize:  int(l.Layout.TextSize.GetValueAtZoomLevel(zoomLevel)),
				TextColor: l.Paint.TextColor.GetColorAtZoomLevel(zoomLevel),
				ZIndex:    layerIndex,
			}
		}
	default:
		log.Printf("unhandled node layer type: %q", l.Type)
	}

	return nil
}

func (l *Layer) GetLayerWayStyle(tags []*ownmap.OSMTag, zoomLevel ownmap.ZoomLevel, layerIndex int) *styling.WayStyle {
	switch l.Type {
	case "line", "symbol":
		tagsInSourceLayer := areTagsInSourceLayer(l.SourceLayer, tags)
		if !tagsInSourceLayer {
			// OSM Way doesn't "belong" in this sourceLayer, skip everything
			return nil
		}

		shown := l.Filter.IsObjectShown(l.SourceLayer, tags, ownmap.ObjectTypeWay)
		if !shown {
			return nil
		}

		fillColor := l.Paint.FillColor.GetColorAtZoomLevel(zoomLevel)
		lineColor := l.Paint.LineColor.GetColorAtZoomLevel(zoomLevel)
		lineWidth := l.Paint.LineWidth.GetValueAtZoomLevel(zoomLevel)

		if (lineColor == nil || lineWidth == 0) && fillColor == nil {
			// there is neither a line or fill color, so don't show this item
			return nil
		}

		return &styling.WayStyle{
			FillColor:      fillColor,
			LineColor:      lineColor,
			LineDashPolicy: l.Paint.LineDashArray,
			LineWidth:      lineWidth,
			ZIndex:         layerIndex,
		}
	default:
		log.Printf("unhandled way layer type: %q", l.Type)
	}

	return nil
}

type Light struct {
	Anchor    string      `json:"anchor"`
	Color     color.Color `json:"color"`
	Intensity float64     `json:"intensity"`
	Position  []float64   `json:"position"`
}

type Source struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Sources map[string]Source

type Transition struct {
	Delay    int `json:"delay"`    // milliseconds
	Duration int `json:"duration"` // milliseconds
}
