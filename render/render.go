package render

import (
	"math"

	"osm-map/astar"
	"osm-map/data"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func convert(lat, lon float64, minLat, maxLat, minLon, maxLon float64, width, height float64) (float64, float64) {
	x := ((lon - minLon) / (maxLon - minLon)) * width
	y := ((lat - minLat) / (maxLat - minLat)) * height
	return x, y
}

func closestNode(pos pixel.Vec, nodes map[int64]pixel.Vec) int64 {
	var closestID int64 = -1
	minDist := math.MaxFloat64
	for id, nodePos := range nodes {
		dist := pos.To(nodePos).Len()
		if dist < minDist {
			minDist = dist
			closestID = id
		}
	}
	return closestID
}

func run() {
	wWidth, wHeight := 1024.0, 768.0
	cfg := pixelgl.WindowConfig{
		Title:  "Map",
		Bounds: pixel.R(0, 0, wWidth, wHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	d, err := data.LoadOSMData()
	if err != nil {
		panic(err)
	}

	var minLat, maxLat = 90.0, -90.0
	var minLon, maxLon = 180.0, -180.0
	for _, el := range d.Elements {
		if el.Type != "node" {
			continue
		}
		if el.Lat < minLat {
			minLat = el.Lat
		}
		if el.Lat > maxLat {
			maxLat = el.Lat
		}
		if el.Lon < minLon {
			minLon = el.Lon
		}
		if el.Lon > maxLon {
			maxLon = el.Lon
		}
	}

	nodes := make(map[int64]pixel.Vec)
	graph := &astar.Graph{Nodes: make(map[int64]*astar.Node)}

	for _, el := range d.Elements {
		if el.Type != "node" {
			continue
		}
		x, y := convert(el.Lat, el.Lon, minLat, maxLat, minLon, maxLon, wWidth, wHeight)
		nodes[el.ID] = pixel.V(x, y)
		graph.Nodes[el.ID] = &astar.Node{
			ID:  el.ID,
			Lat: el.Lat,
			Lon: el.Lon,
		}
	}

	for _, el := range d.Elements {
		if el.Type != "way" {
			continue
		}
		for i := 0; i < len(el.Nodes)-1; i++ {
			a := el.Nodes[i]
			b := el.Nodes[i+1]
			graph.Nodes[a].Adj = append(graph.Nodes[a].Adj, b)
			graph.Nodes[b].Adj = append(graph.Nodes[b].Adj, a)
		}
	}

	startID := int64(-1)
	goalID := int64(-1)
	var path []int64
	pathIndex := 1
	pathFound := false
	animTimer := 0.0

	imd := imdraw.New(nil)

	for !win.Closed() {
		win.Clear(colornames.White)
		imd.Clear()

		// Draw roads
		for _, el := range d.Elements {
			if el.Type != "way" {
				continue
			}
			hwyType, ok := el.Tags["highway"]
			if !ok {
				continue
			}

			switch hwyType {
			case "motorway", "unclassified", "pedestrian", "footway", "steps", "cycleway", "path", "busway", "raceway", "corridor":
				continue
			case "primary", "primary_link", "secondary", "secondary_link":
				imd.Color = colornames.Black
			default:
				imd.Color = colornames.Darkgray
			}

			for i := 0; i < len(el.Nodes)-1; i++ {
				a, aOk := nodes[el.Nodes[i]]
				b, bOk := nodes[el.Nodes[i+1]]
				if aOk && bOk {
					imd.Push(a, b)
					switch hwyType {
					case "primary", "primary_link":
						imd.Line(5)
					case "secondary", "secondary_link":
						imd.Line(3.5)
					default:
						imd.Line(3)
					}
				}
			}
		}

		// Handle click
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			mousePos := win.MousePosition()
			clickedNodeID := closestNode(mousePos, nodes)
			if startID == -1 {
				startID = clickedNodeID
			} else if goalID == -1 && clickedNodeID != startID {
				goalID = clickedNodeID
				path = astar.Algorithm(graph, startID, goalID)
				pathFound = true
				pathIndex = 1
				animTimer = 0
			}
		}

		// Draw path
		if pathFound && pathIndex < len(path) {
			animTimer += 1 / 1000.0
			if animTimer >= 0.01 {
				animTimer = 0
				aID := path[pathIndex-1]
				bID := path[pathIndex]
				aPos, aOk := nodes[aID]
				bPos, bOk := nodes[bID]
				if aOk && bOk {
					imd.Color = colornames.Blue
					imd.Push(aPos, bPos)
					imd.Line(5)
				}
				pathIndex++
			}
		} else if pathFound {
			for i := 0; i < len(path)-1; i++ {
				a, aOk := nodes[path[i]]
				b, bOk := nodes[path[i+1]]
				if aOk && bOk {
					imd.Color = colornames.Blue
					imd.Push(a, b)
					imd.Line(5)
				}
			}
		}
		if pos, ok := nodes[startID]; ok {
			imd.Color = colornames.Green
			imd.Push(pos)
			imd.Circle(3, 0)
		}
		if pos, ok := nodes[goalID]; ok {
			imd.Color = colornames.Red
			imd.Push(pos)
			imd.Circle(3, 0)
		}

		imd.Draw(win)
		win.Update()
	}
}

func Init() {
	pixelgl.Run(run)
}
