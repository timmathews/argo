{{define "content"}}
<h1>Web Apps</h1>
<table class="table table-sm table-striped">
  <thead>
    <tr>
      <th></th>
      <th>Installed Version</th>
      <th>Name</th>
      <th>Description</th>
      <th>Authors</th>
      <th>Link</th>
    </tr>
  </thead>
  <tbody>
  {{range $idx, $app := .}}
    <tr>
      <td>
        <form class="appInstall">
          <input type="hidden" name="name" value="{{$app.Package.Name}}">
          <input type="hidden" name="version" value="{{$app.Package.Version}}">
          <button class="btn btn-light btn-sm btn-block">
            Install {{$app.Package.Version}}
          </button>
        </form>
      </td>
      <td>{{$app.Version}}</td>
      <td>
        {{if $app.Path}}
          <a href="{{$app.Path}}">
            {{$app.Package.Name}}
          </a>
        {{else}}
          {{$app.Package.Name}}
        {{end}}
      </td>
      <td>{{$app.Package.Description}}</td>
      <td>{{$app.Package.Author.Name}}</td>
      <td><a href="{{$app.Package.Links.Npm}}">npm</a></td>
    </tr>
  {{end}}
  </tbody>
</table>
{{end}}
