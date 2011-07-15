package main

import (
	"http"
	"fmt"
	"github.com/droundy/gui"
)

// FIXME: There is no way to make the following return something that
// is both a gui.Bool and a gui.HasText.
type labelledcheckbox struct {
	gui.Widget
	gui.String
	gui.Changeable
	gui.Bool
}
func LabelledCheckbox(l string) interface { gui.Widget; gui.String; gui.Changeable; gui.Bool } {
	cb := gui.Checkbox()
	label := gui.Text(l)
	table := gui.Row(cb, label)
	label.OnClick(func() gui.Refresh {
		cb.Toggle()
		return cb.HandleChange()
	})
	out := labelledcheckbox{ table, label, cb, cb }
	return &out
}

type radiobuttons struct {
	gui.Widget
	gui.String
	gui.Changeable
}
	
func RadioButtons(vs... string) interface{ gui.Widget; gui.String; gui.Changeable } {
	var bs []interface{ gui.Changeable; gui.Bool; gui.String }
	var ws []gui.Widget
	for _,v := range vs {
		b := gui.RadioButton(v)
		bs = append(bs, b)
		ws = append(ws, b)
	}
	col := gui.Column(ws...)
	grp := gui.RadioGroup(bs...)
	return &radiobuttons{ col, grp, grp }
}


func main() {
	http.HandleFunc("/style.css", styleServer)
	buttonA := gui.Button("A")
	buttonB := gui.Button("B")
	buttonA.OnClick(func() gui.Refresh {
		fmt.Println("I clicked on button A")
		buttonB.SetString(buttonB.GetString() + buttonB.GetString())
		return gui.StillClean
	})
	buttonB.OnClick(func() gui.Refresh {
		fmt.Println("I clicked on button A")
		t := buttonB.GetString()
		buttonB.SetString(t[:len(t)/2+1])
		return gui.StillClean
	})
	iscool := gui.Checkbox()
	name := gui.EditText("Enter name here")
	hello := gui.Text("Hello world!")
	name.OnChange(func() gui.Refresh {
		hello.SetString("Hello " + name.GetString() + "!")
		return gui.StillClean
	})
	testing_checkbox := LabelledCheckbox("testing")
	testing_checkbox.OnChange(func() gui.Refresh {
		fmt.Println("Hello world")
		if testing_checkbox.GetBool() {
			testing_checkbox.SetString("this test is currently true")
		} else {
			testing_checkbox.SetString("this test is now false")
		}
		return gui.NeedsRefresh
	})

	// Now let's test out a set of radio buttons
	radio := RadioButtons("apples", "lemons", "oranges", "pears")
	radio_report := gui.Text("I like to eat tasty fruit")
	menu := gui.Menu("apples", "lemons", "oranges", "pears")
	radio.OnChange(func() gui.Refresh {
		menu.SetString(radio.GetString());
		radio_report.SetString("I like to eat " + radio.GetString())
		return gui.NeedsRefresh
	})

	menu.OnChange(func() gui.Refresh {
		radio.SetString(menu.GetString())
		radio_report.SetString("I like to eat " + radio.GetString())
		return radio.HandleChange()
	})

	err := gui.Run(12346,
		gui.Column(
		iscool,
		testing_checkbox,
		gui.Row(buttonA, buttonB),
		gui.Row(gui.Text("Name:"), name),
		hello,
		gui.Text("Goodbye world!"),
		radio,
		radio_report,
		menu,
		))
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}

func styleServer(c http.ResponseWriter, req *http.Request) {
	c.Header().Set("Content-Type", "text/css")
	fmt.Fprint(c, `
html {
    margin: 0;
    padding: 0;
}

body {
    margin: 0;
    padding: 0;
    background: #ffffff;
    font-family: arial,helvetica,"sans serif";
    font-size: 12pt;
}
h1 {
font-family: verdana,helvetica,"sans serif";
font-weight: bold;
font-size: 16pt;
}
h2 { font-family: verdana,helvetica,"sans serif";
font-weight: bold;
font-size: 14pt;
}
p {
font-family: arial,helvetica,"sans serif";
font-size:12pt;
}
li {
  font-family: arial,helvetica,"sans serif";
  font-size: 12pt;
}
a {
  color: #555599;
}
`)
}
