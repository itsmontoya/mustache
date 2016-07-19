package mustache

import (
	"errors"
	"fmt"
	"testing"

	hmust "github.com/hoisie/mustache"
)

var (
	exampleVerySimpleStr        = "<div><p>{{ name }}<p></div>"
	exampleSimpleStr            = "<div><p>Hello {{ name }}, {{ question }}?</p></div>"
	exampleInjectionStr         = "<div>{{ injection }}</div>"
	exampleApprovedInjectionStr = "<head>{{{approvedInjection}}}</head>"
	exampleSimpleHTMLStr        = "<li>{{ greeting }} {{ name }}!</li>"
	exampleArrayHTMLStr         = "{{# . }}<ul><li>{{ greeting }} {{ name }}!</li></ul>{{/ . }}"
	exampleHTMLStr              = `
<html>
	<head>
		<title>{{ title }}</title>
	</head>
	<body>
		{{# basic }}
		<ul>
			{{# users }}<li>{{ greeting }} {{ name }}!</li>{{/ users }}
		</ul>
		{{/ basic }}
	</body>
</html>
`
	exampleLongStr = `
	But I must explain to you {{ name }} how all this mistaken idea of denouncing pleasure and praising pain was born and I will give you a complete account of the system, and expound the actual teachings of the great explorer of the truth, the master-builder of human happiness. May I ask you, {{ question }}? No one rejects, dislikes, or avoids pleasure itself, because it is pleasure, but because those who do not know how to pursue pleasure rationally encounter consequences that are extremely painful.

	I must say, {{ name }} Nor again is there anyone who loves or pursues or desires to obtain pain of itself, because it is pain, but because occasionally circumstances occur in which toil and pain can procure him some great pleasure. To take a trivial example, which of us ever undertakes laborious physical exercise, except to obtain some advantage from it?

	But who has any right to find fault with a man who chooses to enjoy a pleasure that has no annoying consequences, or one who avoids a pain that produces no resultant pleasure? How about this? {{ injection }}.. haha, you like that?
`

	exampleVerySimple        = []byte(exampleVerySimpleStr)
	exampleSimple            = []byte(exampleSimpleStr)
	exampleInjection         = []byte(exampleInjectionStr)
	exampleApprovedInjection = []byte(exampleApprovedInjectionStr)
	exampleSimpleHTML        = []byte(exampleSimpleHTMLStr)
	exampleArrayHTML         = []byte(exampleArrayHTMLStr)
	exampleHTML              = []byte(exampleHTMLStr)
	exampleLong              = []byte(exampleLongStr)

	expectedVerySimple        = "<div><p>Panda<p></div>"
	expectedSimple            = "<div><p>Hello Panda, how are you doing today?</p></div>"
	expectedInjectionA        = "<div>&lt;script src=&#34;badsite.com/injectionPage&#34;/&gt;</div>"
	expectedInjectionB        = "<div>&#60;script src=&#34;badsite.com/injectionPage&#34;/&#62;</div>"
	expectedApprovedInjection = "<head><script src=\"goodsite.com/apistuff\"/></head>"

	m = map[string][]byte{
		"name":              []byte("Panda"),
		"question":          []byte("how are you doing today"),
		"injection":         []byte("<script src=\"badsite.com/injectionPage\"/>"),
		"approvedInjection": []byte("<script src=\"goodsite.com/apistuff\"/>"),
	}

	users = []Aficionado{
		&User{"Hello", "Joe"},
		&User{"Greetings", "David"},
		&User{"What's up", "Derpson"},
	}

	//map[string]map[string]User{
	dm = map[string]interface{}{
		"title": "Hello world!",
		"basic": map[string]interface{}{
			"users": users,
		},
	}

	errInvalidOutput = errors.New("invalid output")
)

type Page struct {
	Title string
	Basic Basic
}

type Basic struct {
	Users []*User
}

type User struct {
	Greeting string
	Name     string
}

func (u *User) MarshalMustache(r *Renderer) error {
	r.ForEach(func(key string) (val interface{}) {
		switch key {
		case "greeting":
			val = u.Greeting
		case "name":
			val = u.Name
		}

		return
	})
	return nil
}

var (
	outputB   []byte
	outputStr string
)

func test(tmpl []byte, d interface{}) (err error) {
	var t *Template
	if t, err = Parse(tmpl); err != nil {
		fmt.Println("Parse error", err)
		return
	}

	if err = t.Render(d, func(b []byte) {
		//return
		fmt.Println(string(b))
	}); err != nil {
		fmt.Println("Render error")
		return
	}

	return
}

func TestVerySimple(t *testing.T) {
	if err := test(exampleVerySimple, m); err != nil {
		t.Error(err)
	}
}

func TestSimple(t *testing.T) {
	if err := test(exampleSimple, m); err != nil {
		t.Error(err)
	}
}

func TestApprovedInjection(t *testing.T) {
	if err := test(exampleApprovedInjection, m); err != nil {
		t.Error(err)
	}
}

func TestSimpleHTML(t *testing.T) {
	if err := test(exampleSimpleHTML, users[0]); err != nil {
		t.Error(err)
	}
}

func TestArrayHTML(t *testing.T) {
	if err := test(exampleArrayHTML, users); err != nil {
		t.Error(err)
	}
}

func TestHTML(t *testing.T) {
	if err := test(exampleHTML, dm); err != nil {
		t.Error(err)
	}
}

func TestLong(t *testing.T) {
	var err error
	if err = Render(exampleLong, m, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}
}

func BenchmarkVerySimple(b *testing.B) {
	benchmark(b, exampleVerySimple, m)
}

func BenchmarkSimple(b *testing.B) {
	benchmark(b, exampleSimple, m)
}

func BenchmarkInjection(b *testing.B) {
	benchmark(b, exampleInjection, m)
}

func BenchmarkApprovedInjection(b *testing.B) {
	benchmark(b, exampleApprovedInjection, m)
}

func BenchmarkSimpleHTML(b *testing.B) {
	benchmark(b, exampleSimpleHTML, users[0])
}

func BenchmarkArrayHTML(b *testing.B) {
	benchmark(b, exampleArrayHTML, users)
}

func BenchmarkHTML(b *testing.B) {
	benchmark(b, exampleHTML, dm)
}

func BenchmarkLong(b *testing.B) {
	benchmark(b, exampleLong, m)
}

func benchmark(b *testing.B, tgt []byte, data interface{}) {
	b.StopTimer()
	var (
		tmpl *Template
		err  error
	)

	if tmpl, err = Parse(tgt); err != nil {
		b.Error(err)
		return
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if err = tmpl.Render(data, func(b []byte) {
			outputB = b
		}); err != nil {
			b.Error(err)
			return
		}
	}

	b.ReportAllocs()
}

func benchmarkHoisie(b *testing.B, tgt string, data interface{}) {
	b.StopTimer()
	var (
		tmpl *hmust.Template
		err  error
	)

	if tmpl, err = hmust.ParseString(tgt); err != nil {
		b.Error(err)
		return
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		outputStr = tmpl.Render(data)
	}

	b.ReportAllocs()
}

func BenchmarkHoisieVerySimple(b *testing.B) {
	benchmarkHoisie(b, exampleVerySimpleStr, m)
}

func BenchmarkHoisieSimple(b *testing.B) {
	benchmarkHoisie(b, exampleSimpleStr, m)
}

func BenchmarkHoisieInjection(b *testing.B) {
	benchmarkHoisie(b, exampleInjectionStr, m)
}

func BenchmarkHoisieApprovedInjection(b *testing.B) {
	benchmarkHoisie(b, exampleApprovedInjectionStr, m)
}

func BenchmarkHoisieSimpleHTML(b *testing.B) {
	benchmarkHoisie(b, exampleSimpleHTMLStr, users[0])
}

func BenchmarkHoisieArrayHTML(b *testing.B) {
	benchmarkHoisie(b, exampleArrayHTMLStr, users)
}

func BenchmarkHoisieHTML(b *testing.B) {
	benchmarkHoisie(b, exampleHTMLStr, dm)
}

func BenchmarkHoisieLong(b *testing.B) {
	benchmarkHoisie(b, exampleLongStr, m)
}
