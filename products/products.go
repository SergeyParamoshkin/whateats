package products

import (
	"bytes"
	"html/template"
	"log"
)

type Menu struct {
	Data     map[string]map[string][]string `json:"menu"`
	Template *template.Template
}

type Diet struct {
	Diet     map[string][]string
	Template *template.Template
}

type RMenu struct {
	Day       string
	Ingestion []string
}

func (d *Diet) GetBreakfast() []string {
	return d.Diet["breakfast"]
}

func (d *Diet) GetBrunch() []string {
	return d.Diet["brunch"]
}

func (d *Diet) GetLunch() []string {
	return d.Diet["lunch"]
}

func (d *Diet) GetAfternoonSnack() []string {
	return d.Diet["afternoonSnack"]
}

func (d *Diet) GetDinner() []string {
	return d.Diet["dinner"]
}

func (d *Diet) GetSupper() []string {
	return d.Diet["supper"]
}

func (d *Diet) GetAll() string {
	type meal struct {
		Meal        string
		Description string
	}

	listMeals := []meal{
		{
			"breakfast",
			"завтрак",
		}, {
			"brunch",
			"перекус",
		}, {
			"lunch",
			"обед",
		},
		{
			"afternoonSnack",
			"перекус",
		},
		{
			"dinner",
			"полдник",
		},
		{
			"supper",
			"ужин",
		},
	}

	if len(d.Diet) == 0 {
		return "Обратитесь к @sergeyparamoshkin за составлением меню"
	}
	r := make([]RMenu, 0)
	for _, m := range listMeals {
		r = append(r, RMenu{
			m.Description,
			d.Diet[m.Meal],
		})
	}
	buf := new(bytes.Buffer)
	err := d.Template.Execute(buf, r)
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

func (m *Menu) DayOfWeek(dayofweek string) *Diet {
	return &Diet{
		m.Data[dayofweek],
		m.Template,
	}
}

func newTemplagete() (*template.Template, error) {
	return template.New("telegram markdown").Parse(`
{{range $s := . }}
*{{ .Day }}*{{ $length := len .Ingestion }}{{ if gt $length 0 }}{{ range $i := .Ingestion }}
__{{ $i }}__{{ end }}{{ end }}{{ end }}`)
}

func NewMenu(userName string) *Menu {
	m, err := fromFile(userName + ".json")
	if err != nil {
		log.Println(err)
		return &Menu{}
	}
	m.Template, err = newTemplagete()
	if err != nil {
		log.Println(err)
	}
	return m
}
