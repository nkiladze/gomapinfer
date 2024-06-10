package common

import (
    "github.com/dhconnelly/rtreego"
    "math"
)

// Function to create a Rect, returning a value not a pointer
func RtreegoRect(r Rectangle) rtreego.Rect {
    dx := math.Max(0.00000001, r.Max.X - r.Min.X)
    dy := math.Max(0.00000001, r.Max.Y - r.Min.Y)
    rect, err := rtreego.NewRect(rtreego.Point{r.Min.X, r.Min.Y}, []float64{dx, dy})
    if err != nil {
        panic(err)
    }
    return *rect
}

type edgeSpatial struct {
    edge *Edge
    rect rtreego.Rect // Store rect as a value, not a pointer
}

// Bounds must return a value to satisfy rtreego.Spatial
func (e *edgeSpatial) Bounds() rtreego.Rect {
    if e.rect.Min == (rtreego.Point{}) && e.rect.Max == (rtreego.Point{}) {
        r := e.edge.Src.Point.Rectangle()
        r = r.Extend(e.edge.Dst.Point)
        e.rect = RtreegoRect(r)
    }
    return e.rect
}

type Rtree struct {
    tree *rtreego.Rtree
}

func (rtree *Rtree) Search(rect Rectangle) []*Edge {
    spatials := rtree.tree.SearchIntersect(RtreegoRect(rect))
    edges := make([]*Edge, len(spatials))
    for i, spatial := range spatials {
        es, ok := spatial.(*edgeSpatial)
        if !ok {
            panic("type assertion to *edgeSpatial failed")
        }
        edges[i] = es.edge
    }
    return edges
}

func (graph *Graph) Rtree() *Rtree {
    rtree := rtreego.NewTree(2, 25, 50)
    for _, edge := range graph.Edges {
        es := &edgeSpatial{edge: edge} // Create a pointer to edgeSpatial
        rtree.Insert(es)               // Insert the pointer
    }
    return &Rtree{tree: rtree}
}
