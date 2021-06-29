{{define "content"}}
<h1>N2K Debugging and Reverse Engineering</h1>
<table class="table table-sm table-striped">
  <thead>
    <tr>
      <th>Timestamp</th>
      <th>PGN</th>
      <th>Source</th>
      <th>Destination</th>
      <th>Priority</th>
      <th>Length</th>
      <th>Data</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>2021-06-25 15:09:30</td>
      <td>130824: Simnet Unknown</td>
      <td>12</td>
      <td>255</td>
      <td>3</td>
      <td>76</td>
      <td>
<pre><code>7d 99 64 20 00 00 0d 21   00 00 69 20 00 00 d3 20
00 00 81 40 00 00 00 00   0a 21 00 00 0c 21 00 00
35 20 2a da 7d 20 78 00   53 20 45 1b 75 40 20 6c
fb ff 09 41 00 00 00 00   0b 41 00 00 00 00 d0 40
00 00 00 00 1d 21 00 00   9d 20 b7 7a</code></pre>
        </td>
    </tr>
  {{range $idx, $msg := .}}
  {{end}}
  </tbody>
</table>
<hr>
<table class="table table-sm table-hover table-striped">
  <thead>
    <tr>
      <th>Field</th>
      <th>Width</th>
      <th>Resolution</th>
      <th>Signed</th>
      <th>Units</th>
      <th>Value</th>
      <th></th>
    </tr>
  </thead>
  <tbody id="details">
  </tbody>
</table>
<hr>

<script type="text/html" id="row-template">
  <tr id="row-{%= id %}">
    <td>{%= field %}</td> 
    <td><input id="width" type="number" class="form-control form-control-sm"></td> 
    <td>
      <select id="resolution" class="form-control form-control-sm">
        <option>Select</option>
        <option>---</option>
        <option value="1">Not Used</option>
        <option value="2">Degrees</option>
      </select>
    </td> 
    <td>
      <div class="form-check form-check-inline">
        <input type="checkbox" id="signed" value="true">
        <input id="signed" value="false" type="hidden">
      </div>
    </td> 
    <td><input id="units" class="form-control form-control-sm" placeholder="Units"></td> 
    <td></td> 
    <td>
      <div class="btn-group">
        <button class="btn btn-sm btn-danger" onClick="removeRow({%= id %})">
          <img src="/assets/svg/dash.svg" class="octicon">
        </button>
        <button class="btn btn-sm btn-success" onClick="addRow({%= id %})">
          <img src="/assets/svg/plus.svg" class="octicon">
        </button>
      </div>
    </td> 
  </tr>
</script>
{{end}}