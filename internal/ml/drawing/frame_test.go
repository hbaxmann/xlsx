// Copyright (c) 2017 Andrey Gayvoronsky <plandem@gmail.com>
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package drawing_test

import (
	"bytes"
	"encoding/xml"
	"github.com/plandem/ooxml/drawing/dml"
	"github.com/plandem/ooxml/drawing/dml/chart"
	"github.com/plandem/ooxml/ml"
	"github.com/plandem/xlsx/internal/ml/drawing"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestFrame(t *testing.T) {
	type Entity struct {
		XMLName xml.Name       `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing entity"`
		DMLName dml.Name       `xml:",attr"`
		Frame   *drawing.Frame `xml:"graphicFrame"`
	}

	data := strings.NewReplacer("\t", "", "\n", "").Replace(`
	<xdr:entity xmlns:xdr="http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
		<xdr:graphicFrame macro="">
			<xdr:nvGraphicFramePr>
				<xdr:cNvPr id="2" name="Chart 1"></xdr:cNvPr>
				<xdr:cNvGraphicFramePr></xdr:cNvGraphicFramePr>
			</xdr:nvGraphicFramePr>
			<xdr:xfrm>
				<a:off x="1" y="2"></a:off>
				<a:ext cx="3" cy="4"></a:ext>
			</xdr:xfrm>
			<a:graphic>
				<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/chart">
					<c:chart xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" r:id="rId1"></c:chart>
				</a:graphicData>
			</a:graphic>
		</xdr:graphicFrame>
	</xdr:entity>
`)

	decoder := xml.NewDecoder(bytes.NewReader([]byte(data)))
	entity := &Entity{}
	err := decoder.DecodeElement(entity, nil)
	require.Nil(t, err)

	frame := &drawing.Frame{
		NonVisual: &drawing.FrameNonVisual{
			DrawingProperties: &dml.NonVisualCommonProperties{
				ID:   2,
				Name: "Chart 1",
			},
			FrameProperties: &dml.NonVisualFrameProperties{},
		},
		Graphic: &dml.GraphicalObject{
			Data: &dml.GraphicalObjectData{
				Uri: "http://schemas.openxmlformats.org/drawingml/2006/chart",
				Chart: &chart.Ref{
					RID: "rId1",
				},
			},
		},
		Transform: &dml.Transform2D{
			Offset: &dml.Point2D{
				X: 1,
				Y: 2,
			},
			Size: &dml.PositiveSize2D{
				Height: 3,
				Width:  4,
			},
		},
		ReservedAttributes: ml.ReservedAttributes{
			Attrs: []xml.Attr{
				{
					Name: xml.Name{
						Local: "macro",
					},
				},
			},
		},
	}

	require.Equal(t, &Entity{
		XMLName: xml.Name{
			Space: "http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing",
			Local: "entity",
		},
		Frame: frame,
	}, entity)

	//encode data should be same as original
	encode, err := xml.Marshal(entity)
	require.Nil(t, err)
	require.Equal(t, strings.NewReplacer("xdr:", "", ":xdr", "").Replace(data), string(encode))
}

func TestFrameNonVisual(t *testing.T) {
	type Entity struct {
		XMLName   xml.Name                `xml:"http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing entity"`
		DMLName   dml.Name                `xml:",attr"`
		NonVisual *drawing.FrameNonVisual `xml:"nvGraphicFramePr"`
	}

	data := strings.NewReplacer("\t", "", "\n", "").Replace(`
	<xdr:entity xmlns:xdr="http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">
		<xdr:nvGraphicFramePr>
			<xdr:cNvPr id="2" name="Chart 1"></xdr:cNvPr>
			<xdr:cNvGraphicFramePr></xdr:cNvGraphicFramePr>
		</xdr:nvGraphicFramePr>
	</xdr:entity>
`)

	decoder := xml.NewDecoder(bytes.NewReader([]byte(data)))
	entity := &Entity{}
	err := decoder.DecodeElement(entity, nil)
	require.Nil(t, err)

	object := &drawing.FrameNonVisual{
		DrawingProperties: &dml.NonVisualCommonProperties{
			ID:   2,
			Name: "Chart 1",
		},
		FrameProperties: &dml.NonVisualFrameProperties{},
	}

	require.Equal(t, &Entity{
		XMLName: xml.Name{
			Space: "http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing",
			Local: "entity",
		},
		NonVisual: object,
	}, entity)

	//encode data should be same as original
	encode, err := xml.Marshal(entity)
	require.Nil(t, err)
	require.Equal(t, strings.NewReplacer("xdr:", "", ":xdr", "").Replace(data), string(encode))
}