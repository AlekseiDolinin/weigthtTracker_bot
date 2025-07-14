package plots

import (
	"bytes"
	"fmt"
	"weightTrack_bot/models"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

// Структура для построения графика
type DataPoint struct {
	X     float64
	Y     float64
	Label string
}

// Конструктор структуры
func NewDataPoint(x float64, y float64, l string) DataPoint {
	return DataPoint{X: x, Y: y, Label: l}
}

// Преобразует данные из формата слайса структур []models.AvgRecordsPeriod в []DataPoint
func FromARPtoDP(result []models.AvgRecordsPeriod) []DataPoint {
	var data []DataPoint
	xAxis := 1.0
	for i := len(result) - 1; i >= 0; i-- {
		var diffStr string
		if i > 0 {
			diff := result[i].GetWeight() - result[i-1].GetWeight()
			diffStr = fmt.Sprintf("%+.2f", diff)
		}
		newPoint := NewDataPoint(xAxis, result[i].GetWeight(), diffStr)
		data = append(data, newPoint)
		xAxis += 1.0
	}
	return data
}

// Создает график
func MakePlot(result []models.AvgRecordsPeriod) ([]byte, error) {
	// Создаем новый график
	p := plot.New()
	p.Title.Text = "График иизмерения веса по дням"
	p.X.Label.Text = "Дни"
	p.Y.Label.Text = "Вес в килограммах"

	data := FromARPtoDP(result)

	// Создаем точки для графика
	points := make(plotter.XYs, len(data))
	for i, d := range data {
		points[i].X = d.X
		points[i].Y = d.Y
	}

	// Добавляем точки на график
	line, err := plotter.NewLine(points)
	if err != nil {
		return nil, err
	}
	p.Add(line)

	// Создаем canvas в памяти
	c := vgimg.New(10*vg.Centimeter, 10*vg.Centimeter)
	p.Draw(draw.New(c))

	// Создаем буфер для PNG изображения
	var buf bytes.Buffer
	pngWriter := vgimg.PngCanvas{Canvas: c}

	// Записываем PNG в буфер
	if _, err := pngWriter.WriteTo(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
