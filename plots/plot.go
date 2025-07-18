package plots

import (
	"bytes"
	"fmt"
	"image/color"
	"strconv"
	"time"
	"weightTrack_bot/backup"
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
	Date  time.Time
}

func (d DataPoint) GetDate() (t time.Time) {
	return d.Date
}

// Конструктор структуры
func NewDataPoint(x float64, y float64, l string, date time.Time) DataPoint {
	return DataPoint{X: x, Y: y, Label: l, Date: date}
}

// Кастомный Tick.Marker для подмены чисел на даты
type customTicks struct {
	Labels []time.Time
	Format string
}

func (ct customTicks) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := 0; i < len(ct.Labels) && float64(i+1) <= max; i++ {
		ticks = append(ticks, plot.Tick{
			Value: float64(i + 1),                   // Позиция тика (1, 2, 3...)
			Label: strconv.Itoa(ct.Labels[i].Day()), // Подпись (дата)
		})
	}
	return ticks
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
		newPoint := NewDataPoint(xAxis, result[i].GetWeight(), diffStr, result[i].GetTime())
		data = append(data, newPoint)
		xAxis += 1.0
	}
	return data
}

func GetDatesFromDataPoints(data []DataPoint) (dates []time.Time) {
	for _, d := range data {
		dates = append(dates, d.GetDate())
	}
	return
}

// Создает график
func MakePlot(result []models.AvgRecordsPeriod) ([]byte, error) {
	// Создаем новый график с настройками
	p := plot.New()
	p.Title.Text = "Динамика веса"
	p.Title.TextStyle.Font.Size = font.Length(14)
	p.Title.Padding = 10
	p.X.Label.Text = "Дата"
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

	dates := GetDatesFromDataPoints(data)

	// Настраиваем кастомные метки для оси X
	p.X.Tick.Marker = customTicks{
		Labels: dates,
		Format: "2006-01-02", // Формат даты: "15 Jan", "20 Feb" и т. д.
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
				msg := fmt.Sprintf("Ошибка добавления подписей к точкам графика %v", err)
				backup.WriteLog(msg)
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
	if i, err := pngWriter.WriteTo(&buf); err != nil {
		msg := fmt.Sprintf("Ошибка записи PNG размером %d в буфер %v", i, err)
		backup.WriteLog(msg)
		return nil, err
	}

	return buf.Bytes(), nil
}
