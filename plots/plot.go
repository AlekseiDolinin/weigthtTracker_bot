/*package plots

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
	p.Title.Text = "График динамики веса по дням"
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
*/

package plots

import (
	"bytes"
	"fmt"
	"image/color"
	"weightTrack_bot/models"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
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
			diff := result[i-1].GetWeight() - result[i].GetWeight()
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
	// Создаем новый график с настройками
	p := plot.New()
	p.Title.Text = "Динамика веса"
	p.Title.TextStyle.Font.Size = font.Length(14)
	p.Title.Padding = 10
	p.X.Label.Text = "Дни"
	p.Y.Label.Text = "Вес (кг)"
	p.X.Padding = 5
	p.Y.Padding = 5

	// Настраиваем сетку
	p.Add(plotter.NewGrid())

	// Увеличиваем размер графика
	p.BackgroundColor = color.RGBA{R: 245, G: 245, B: 245, A: 255}
	p.Legend.TextStyle.Font.Size = font.Length(10)

	data := FromARPtoDP(result)

	// Устанавливаем максимальное значение X с запасом
	if len(data) > 0 {
		p.X.Max = data[len(data)-1].X + 0.5 // Добавляем 0.5 единицы справа
	}

	// Создаем точки для графика
	points := make(plotter.XYs, len(data))
	for i, d := range data {
		points[i].X = d.X
		points[i].Y = d.Y
	}

	// Добавляем линию с настройками
	line, err := plotter.NewLine(points)
	if err != nil {
		return nil, err
	}
	line.LineStyle.Width = vg.Points(2)
	line.LineStyle.Color = color.RGBA{R: 0, G: 102, B: 204, A: 255}
	p.Add(line)

	// Добавляем точки на график
	scatter, err := plotter.NewScatter(points)
	if err != nil {
		return nil, err
	}
	scatter.GlyphStyle.Color = color.RGBA{R: 0, G: 102, B: 204, A: 255}
	scatter.GlyphStyle.Radius = vg.Points(4)
	p.Add(scatter)

	// Добавляем подписи к точкам
	for i, d := range data {
		if i > 0 {
			label, err := plotter.NewLabels(plotter.XYLabels{
				XYs:    []plotter.XY{{X: d.X, Y: d.Y + 2}},
				Labels: []string{data[i-1].Label},
			})
			if err != nil {
				return nil, err
			}
			p.Add(label)
		}
	}

	// Настраиваем размеры графика
	width := 15 * vg.Centimeter
	height := 10 * vg.Centimeter

	// Создаем canvas в памяти
	c := vgimg.New(width, height)
	dc := draw.New(c)
	p.Draw(dc)

	// Создаем буфер для PNG изображения
	var buf bytes.Buffer
	pngWriter := vgimg.PngCanvas{Canvas: c}

	// Записываем PNG в буфер
	if _, err := pngWriter.WriteTo(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
