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
	exampleHTML              = []byte(exampleHTMLStr)
	exampleLong              = []byte(exampleLongStr)

	expectedVerySimple        = "<div><p>Panda<p></div>"
	expectedSimple            = "<div><p>Hello Panda, how are you doing today?</p></div>"
	expectedInjectionA        = "<div>&lt;script src=&#34;badsite.com/injectionPage&#34;/&gt;</div>"
	expectedInjectionB        = "<div>&#60;script src=&#34;badsite.com/injectionPage&#34;/&#62;</div>"
	expectedApprovedInjection = "<head><script src=\"goodsite.com/apistuff\"/></head>"

	m = map[string]string{
		"name":              "Panda",
		"question":          "how are you doing today",
		"injection":         "<script src=\"badsite.com/injectionPage\"/>",
		"approvedInjection": "<script src=\"goodsite.com/apistuff\"/>",
	}

	//map[string]map[string]User{
	dm = map[string]interface{}{
		"title": "Hello world!",
		"basic": map[string]interface{}{
			"users": []Aficionado{
				&User{"Hello", "Joe"},
				&User{"Greetings", "David"},
				&User{"What's up", "Derpson"},
			},
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

func (u *User) MarshalMustache(p *parser) error {
	p.ForEach(func(key string) (val interface{}) {
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

func TestBasic(t *testing.T) {
	if err := testSuite(); err != nil {
		t.Error(err)
		return
	}
}

func TestHTML(t *testing.T) {
	if err := parse(exampleHTML, dm, func(b []byte) {
		outputB = b
	}); err != nil {
		t.Error(err)
		return
	}

	fmt.Println("What do we have here?", string(outputB))
}

func TestHoisie(t *testing.T) {
	if err := testHoisieSuite(); err != nil {
		t.Error(err)
		return
	}
}

func testSuite() (err error) {
	var str string

	if err = parse(exampleVerySimple, m, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}

	if str = string(outputB); str != expectedVerySimple {
		return errInvalidOutput
	}

	if err = parse(exampleSimple, m, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}

	if str = string(outputB); str != expectedSimple {
		return errInvalidOutput
	}

	if err = parse(exampleInjection, m, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}

	if str = string(outputB); str != expectedInjectionA && str != expectedInjectionB {
		return errInvalidOutput
	}

	if err = parse(exampleApprovedInjection, m, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}

	if str = string(outputB); str != expectedApprovedInjection {
		return errInvalidOutput
	}

	if err = parse(exampleHTML, dm, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}

	return
}

func runSuite() (err error) {
	if err = parse(exampleVerySimple, m, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}

	if err = parse(exampleSimple, m, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}

	if err = parse(exampleInjection, m, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}

	if err = parse(exampleApprovedInjection, m, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}

	return
}

func TestLong(t *testing.T) {
	var err error
	if err = parse(exampleLong, m, func(b []byte) {
		outputB = b
	}); err != nil {
		return
	}
}

func runHoisieSuite() (err error) {
	outputStr = hmust.Render(exampleVerySimpleStr, m)
	outputStr = hmust.Render(exampleSimpleStr, m)
	outputStr = hmust.Render(exampleInjectionStr, m)
	outputStr = hmust.Render(exampleApprovedInjectionStr, m)

	return
}

func testHoisieSuite() (err error) {
	var str string
	if str = hmust.Render(exampleVerySimpleStr, m); str != expectedVerySimple {
		fmt.Println("No match?", str)
		return errInvalidOutput
	}

	if str = hmust.Render(exampleSimpleStr, m); str != expectedSimple {
		fmt.Println("No match?", str)

		return errInvalidOutput
	}

	if str = hmust.Render(exampleInjectionStr, m); str != expectedInjectionA && str != expectedInjectionB {
		fmt.Println("No match?", str)

		return errInvalidOutput
	}

	if str = hmust.Render(exampleApprovedInjectionStr, m); str != expectedApprovedInjection {
		fmt.Println("No match?", str)

		return errInvalidOutput
	}

	return
}

func BenchmarkSuite(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		if err = runSuite(); err != nil {
			b.Error(err)
			return
		}
	}

	b.ReportAllocs()
}

func BenchmarkHoisieSuite(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		if err = runHoisieSuite(); err != nil {
			b.Error(err)
			return
		}
	}

	b.ReportAllocs()
}

func BenchmarkLong(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		if err = parse(exampleLong, m, func(b []byte) {
			outputB = b
		}); err != nil {
			return
		}
	}

	b.ReportAllocs()
}

func BenchmarkHoisieLong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		outputStr = hmust.Render(exampleLongStr, m)
	}

	b.ReportAllocs()
}

func BenchmarkHTML(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		if err = parse(exampleHTML, dm, func(b []byte) {
			outputB = b
		}); err != nil {
			return
		}
	}

	b.ReportAllocs()
}

func BenchmarkHoisieHTML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		outputStr = hmust.Render(exampleHTMLStr, dm)
	}

	b.ReportAllocs()
}
