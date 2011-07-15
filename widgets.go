//target:github.com/droundy/gui

package gui

import (
	"fmt"
	"html"
)

func Empty() Widget {
	return Text("")
}

func Text(t string) interface {
	Widget
	String
	Clickable
} {
	return &text{<-newId, t, nil}
}

func EditText(t string) interface {
	Widget
	Changeable
	String
} {
	return &edittext{text{<-newId, t, nil}, nil}
}

func Button(t string) interface {
	Widget
	Clickable
	String
} {
	return &button{text{<-newId, t, nil}, nil}
}

func Checkbox() interface {
	Widget
	Changeable
	Bool
} {
	c := &checkbox{<-newId, false, nil}
	c.OnChange(func() Refresh {
		fmt.Println("I am toggling", c)
		return NeedsRefresh
	})
	return c
}

func Table(rows ...[]Widget) Widget {
	return &table{<-newId, rows}
}

func Column(gui ...Widget) Widget {
	ws := make([][]Widget, len(gui))
	for i := range ws {
		ws[i] = []Widget{gui[i]}
	}
	return &table{<-newId, ws}
}

func Paragraphs(gui ...Widget) Widget {
	return &paragraphs{<-newId, gui}
}

func Row(gui ...Widget) Widget {
	return &table{<-newId, [][]Widget{gui}}
}

/////////////////////////////////////////
// Here is the event-handling stuff... //
/////////////////////////////////////////


type Refresh bool

const (
	NeedsRefresh Refresh = true
	StillClean   Refresh = false
)

type Hook func() Refresh

func (r Refresh) String() string {
	if r {
		return "NeedsRefresh"
	}
	return "StillClean"
}

///////////////////////////////////////
// Everything below this is private! //
///////////////////////////////////////

type paragraphs struct {
	Id
	ws []Widget
}

func (ps *paragraphs) Private__html() (html string, extra []extraCommand) {
	for _, w := range ps.ws {
		whtml, e := w.Private__html()
		extra = append(extra, e...)
		html += "<p>" + whtml + "</p>\n"
	}
	return
}
func (t *paragraphs) Private__getChildren() []Widget {
	return t.ws
}

type table struct {
	Id
	ws [][]Widget
}

func (t *table) Private__html() (html string, extra []extraCommand) {
	html = "<table>\n"
	for i, r := range t.ws {
		class := "even"
		switch {
		case i == 0:
			class = "even first"
		case i&1 == 1:
			class = "odd"
		}
		html += `  <tr class="` + class + `">` + "\n"
		for _, w := range r {
			whtml, wextra := w.Private__html()
			extra = append(extra, wextra...)
			html += "    <td>" + whtml + "</td>\n"
		}
		html += "  </tr>\n"
	}
	html += "</table>\n"
	return
}
func (t *table) Private__getChildren() []Widget {
	out := []Widget{}
	for _, ws := range t.ws {
		out = append(out, ws...)
	}
	return out
}

type text struct {
	Id
	string
	ClickHandler
}

func (dat *text) Private__html() (string, []extraCommand) {
	return `<span onclick="say('onclick:` + string(dat.Private__getId()) + ":" +
		dat.GetString() + `')">` + html.EscapeString(dat.string) + `</span>`, nil
}
func (b *text) GetString() string {
	return b.string
}
func (b *text) SetString(newt string) {
	b.string = newt
}

type edittext struct {
	text
	ChangeHandler
}

func (dat *edittext) Private__html() (string, []extraCommand) {
	h := `<input type="text" onchange="say('onchange:` + string(dat.Private__getId()) + ":" + dat.GetString() +
		`:' + this.value)" value="` + html.EscapeString(dat.text.GetString()) + `" />`
	//fmt.Println(h)
	return h, nil
}

type button struct {
	text
	ClickHandler
}

func (dat *button) Private__html() (string, []extraCommand) {
	return `<input type="submit" onclick="say('onclick:` + string(dat.Private__getId()) + ":" + dat.GetString() + `')" value="` +
		html.EscapeString(dat.GetString()) + `" />`, nil
}

type checkbox struct {
	Id
	BoolValue
	ChangeHandler
}

func (dat *checkbox) Private__html() (string, []extraCommand) {
	checked := ""
	if dat.GetBool() {
		checked = "checked='checked' "
	}
	h := `<input type="checkbox" onchange="say('onchange:` + string(dat.Id) + `')" ` + checked + `" />`
	//fmt.Println(h)
	return h, nil
}

func RadioButton(v string) interface {
	Widget
	Changeable
	Bool
	String
} {
	out := &radiobutton{text{<-newId, v, nil}, false, nil, nil}
	out.OnChange(func() Refresh {
		fmt.Println("I am toggling", out)
		return NeedsRefresh
	})
	out.OnClick(func() Refresh {
		out.Toggle()
		return out.HandleChange()
	})
	return out
}

type radiobutton struct {
	text
	BoolValue
	ChangeHandler
	ClickHandler
}

func (dat *radiobutton) Private__html() (string, []extraCommand) {
	checked := ""
	if dat.GetBool() {
		checked = " checked='checked' "
	}
	return `<input type="radio" onchange="say('onchange:` + string(dat.Private__getId()) + `')" value="` +
		html.EscapeString(dat.GetString()) + `"` + checked +
		`/><span onclick="say('onclick:` + string(dat.Private__getId()) + ":" +
		dat.GetString() + `')">` + html.EscapeString(dat.string) + `</span>`, nil
}

func RadioGroup(butts ...interface {
	Changeable
	Bool
	String
}) interface {
	String
	Changeable
} {
	out := radiogroup{butts, nil}
	numselected := 0
	for i := range butts {
		b := butts[i]
		if b.GetBool() {
			numselected += 1
		}
		b.OnChange(func() Refresh {
			bval := b.GetBool()
			for _, b2 := range out.buttons {
				if b2.GetString() != b.GetString() {
					b2.SetBool(!bval)
				}
			}
			out.HandleChange()
			return NeedsRefresh
		})
	}
	for _, b := range out.buttons {
		if b.GetBool() {
			numselected -= 1
		}
		if numselected == 0 {
			b.SetBool(true)
			numselected -= 1
		} else {
			b.SetBool(false)
		}
	}
	return &out
}

type radiogroup struct {
	buttons []interface {
		Changeable
		Bool
		String
	}
	ChangeHandler
}

func (dat *radiogroup) GetString() string {
	for _, b := range dat.buttons {
		if b.GetBool() {
			return b.GetString()
		}
	}
	panic("Bug: radio group should always have one button selected!")
}

func (dat *radiogroup) SetString(v string) {
	foundstring := false
	for _, b := range dat.buttons {
		if b.GetString() == v {
			foundstring = true
		}
	}
	if foundstring == false {
		panic("Cannot SetString to " + v + " in radio group:  no such option!")
	}
	for _, b := range dat.buttons {
		if b.GetString() == v {
			b.SetBool(true)
		} else {
			b.SetBool(false)
		}
	}
}

func Menu(opts ...string) interface {
	Widget
	String
	Changeable
} {
	return &menu{<-newId, opts, 0, nil}
}

type menu struct {
	Id
	options  []string
	selected int
	ChangeHandler
}

func (m *menu) Private__html() (string, []extraCommand) {
	out := `<select onchange="say('onchange:` + string(m.Private__getId()) + ":" + m.GetString() +
		`:' + this.value)" value="` + html.EscapeString(m.GetString()) + `" />`
	for i, v := range m.options {
		if i == m.selected {
			out += "\n<option value=\"" + v + `" selected='selected'>` + v + "</option>"
		} else {
			out += "\n<option value=\"" + v + `">` + v + "</option>"
		}
	}
	return out + "\n</select>", nil
}

func (m *menu) GetString() string {
	if m.selected >= 0 && m.selected < len(m.options) {
		return m.options[m.selected]
	}
	return ""
}

func (m *menu) SetString(v string) {
	for i, s := range m.options {
		if s == v {
			m.selected = i
		}
	}
}
