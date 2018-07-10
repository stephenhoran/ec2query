package email

import (
	"bytes"
	"html/template"
)

// var tmpl = `
// {{range index . "instances"}}
// 	{{range $k, $v := .}}
// 		{{$k}}
// 	{{end}}
// {{end}}
// `

//Email Templating done here.
var tmpl = `
<body>
<h1>EC2 Instances</h1>
{{range index . "instances"}}
	{{range $key, $value :=  .}}
		{{if $value}}<h3>{{$key}}</h1>
			<table border=1>
				<th>Instance Id</th><th>Type</th><th>State</th><th>Launch Time</th><th>Key</th>
				{{range $k, $i := $value}}
				<tr><td>{{$i.Instanceid}}</td><td>{{$i.Type}}</td><td>{{$i.State}}</td><td>{{$i.LaunchTime}}</td><td>{{$i.KeyName}}</td></tr>
				{{end}}
			</table>
		{{end}}
	{{end}}
{{end}}
<br>
<h1>S3 Buckets</h1>
<table border=1>
	<th>S3 Bucket</th><th>Location</th><th>Size</th>
	{{range index . "s3"}}
		<tr><td>{{.Name}}</td><td>{{.Location}}</td><td>{{.Size}}</tr>
	{{end}}
</table>
</body>
</html>
`

var htmlbody bytes.Buffer

//BuildInstanceTemplate takes an Instance API Map using text templates build HTML email body
//
// Using a bytes buffer as it implements the writer interface requires for the template execution. Make sure to convert back to string.
func BuildInstanceTemplate(instances interface{}) string {
	t := template.New("email")
	t, _ = t.Parse(tmpl)
	err := t.Execute(&htmlbody, instances)

	if err != nil {
		panic(err)
	}

	return htmlbody.String()
}
