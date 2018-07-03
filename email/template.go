package email

import (
	"bytes"
	"html/template"

	"github.com/atssteve/ec2query/apis"
)

var tmpl = `
<body>
{{range .}}
    {{range $key, $value :=  .}}
        {{if $value}}<h1>{{$key}}</h1>
            <table border=1>
                <th>Instance Id</th><th>Type</th><th>State</th><th>Launch Time</th><th>Key</th>
                {{range $k, $i := $value}}
                    <tr><td>{{$i.Instanceid}}</td><td>{{$i.Type}}</td><td>{{$i.State}}</td><td>{{$i.LaunchTime}}</td><td>{{$i.KeyName}}</td></tr>
                {{end}}
            </table>
        {{end}}
    {{end}}
{{end}}
</body>
</html>
`

var htmlbody bytes.Buffer

//BuildInstanceTemplate takes an Instance API Map using text templates build HTML email body
func BuildInstanceTemplate(instances []apis.APIMap) string {
	t := template.New("email")
	t, _ = t.Parse(tmpl)
	err := t.Execute(&htmlbody, instances)

	if err != nil {
		panic(err)
	}

	return htmlbody.String()
}
