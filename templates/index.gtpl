{{define "content"}}
<div class="row">
  <div class="col-3">
    <h2>Statistics</h2>
    <table class="table table-sm table-striped">
      <thead>
        <tr><th>PGN</th><th>Count</th></tr>
      </thead>
      <tbody id="stats"></tbody>
    </table>
  </div>
  <div class="col-7">
    <form id="settings" onsubmit="return submitForm()">
      <div class="tab-content">
        <div role="tabpanel" class="tab-pane active" id="basic-setup">
          <div class="form-group">
            <label for="vesselName">Vessel Name</label>
            <input class="form-control" id="vesselName" placeholder="Name" value="{{.Vessel.Name}}">
          </div>
          <div class="form-row">
            <div class="col">
              <div class="form-group">
                <label for="vesselManufacturer">Manufacturer</label>
                <input class="form-control" id="vesselManufacturer" placeholder="Manufacturer"
                value="{{.Vessel.Manufacturer}}">
              </div>
            </div>
            <div class="col">
              <div class="form-group">
                <label for="vesselModel">Model</label>
                <input class="form-control" id="vesselModel" placeholder="Model" value="{{.Vessel.Model}}">
              </div>
            </div>
            <div class="col">
              <div class="form-group">
                <label for="vesselYear">Year</label>
                <input class="form-control" id="vesselYear" placeholder="Year" type="number" value="{{.Vessel.Year}}">
              </div>
            </div>
          </div>
          <div class="form-row">
            <div class="col">
              <div class="form-group">
                <label for="mmsi">MMSI</label>
                <img src="/assets/svg/info.svg" class="octicon" data-toggle="tooltip" data-placement="top"
                  title="Maritime Mobile Service Identity, assigned regionally, leave blank if you do not have
                  one" />
                <input class="form-control" id="mmsi" placeholder="MMSI" value="{{.Vessel.Mmsi}}">
              </div>
            </div>
            <div class="col">
              <div class="form-group">
                <label for="callsign">Callsign</label>
                <img src="/assets/svg/info.svg" class="octicon" data-toggle="tooltip" data-placement="top"
                  title="Your radio callsign" />
                <input class="form-control" id="callsign" placeholder="Callsign" value="{{.Vessel.Callsign}}">
              </div>
            </div>
            <div class="col">
              <div class="form-group">
                <label for="registration">Registration</label>
                <img src="/assets/svg/info.svg" class="octicon" data-toggle="tooltip" data-placement="top"
                  title="Your boat's registration number, typically displayed on the side of the hull" />
                <input class="form-control" id="registration" placeholder="Registration"
                  value="{{.Vessel.Registration}}">
              </div>
            </div>
          </div>
          <div class="form-group">
            <label for="port">UUID</label><a class="float-right" id="genUUID" href="#">Generate new UUID</a>
            <img src="/assets/svg/info.svg" class="octicon" data-toggle="tooltip" data-placement="top"
              title="The UUID is used to uniquely identify your boat in Signal K, if you have one from a
              previous installation of Signal K enter it here, otherwise click 'Generate new UUID to have one
              created for you" />
            <div class="form-row">
              <div class="col-3">
                <input class="form-control" id="uuid_0" value="{{.Vessel.Uuid0}}">
              </div>
              <div class="col-2">
                <input class="form-control" id="uuid_1" value="{{.Vessel.Uuid1}}">
              </div>
              <div class="col-2">
                <input class="form-control" id="uuid_2" value="{{.Vessel.Uuid2}}">
              </div>
              <div class="col-2">
                <input class="form-control" id="uuid_3" value="{{.Vessel.Uuid3}}">
              </div>
              <div class="col-3">
                <input class="form-control" id="uuid_4" value="{{.Vessel.Uuid4}}">
              </div>
            </div>
          </div>
        </div> <!-- basic-setup -->
        <div role="tabpanel" class="tab-pane" id="network">
          <div class="form-row">
            <div class="col-8">
              <div class="form-group">
                <label for="listenOn">Listen On</label>
                <input class="form-control" id="listenOn" placeholder="IP Address" value="{{.Server.ListenOn}}">
              </div>
            </div>
            <div class="col-4">
              <div class="form-group">
                <label for="port">Port</label>
                <input class="form-control" id="port" type="number" placeholder="Port" value="{{.Server.Port}}">
              </div>
            </div>
          </div>
          <div class="form-check">
            <label class="form-check-label">
              <input class="form-check-input" id="enableWebsockets" type="checkbox"
                {{if .Server.EnableWebsockets}} checked {{end}}>
              Enable WebSockets
            </label>
          </div>
          <div class="form-check">
            <label class="form-check-label">
              <input class="form-check-input" id="useTls" type="checkbox" {{if .Server.UseTls}} checked {{end}}>
              Use TLS
            </label>
          </div>
          <div class="form-group">
            <label for="certificate">Install Certificate</label>
            <input class="form-control-file" id="certificate" type="file" placeholder="Certificate">
          </div>
          <hr>
          <div class="form-check">
            <label class="form-check-label">
              <input class="form-check-input" id="mqttEnable" type="checkbox" {{if .Mqtt.Enable}} checked {{end}}>
              Enable MQTT
            </label>
          </div>
          <div class="form-row">
            <div class="col-8">
              <div class="form-group">
                <label for="mqttHost">MQTT Server</label>
                <input class="form-control" id="mqttHost" placeholder="Hostname" value="{{.Mqtt.Host}}">
              </div>
            </div>
            <div class="col-4">
              <div class="form-group">
                <label for="mqttPort">Port</label>
                <input class="form-control" id="mqttPort" type="number" placeholder="Port" value="{{.Mqtt.Port}}">
              </div>
            </div>
          </div>
          <div class="form-row">
            <div class="col">
              <div class="form-group">
                <label for="mqttClientId">Client ID</label>
                <input class="form-control" id="mqttClientId" placeholder="Client ID" value="{{.Mqtt.ClientId}}">
              </div>
            </div>
            <div class="col">
              <div class="form-group">
                <label for="mqttUsername">Username</label>
                <input class="form-control" id="mqttUsername" placeholder="Username" value="{{.Mqtt.Username}}">
              </div>
            </div>
            <div class="col">
              <div class="form-group">
                <label for="mqttPassword">Password</label>
                <input class="form-control" id="mqttPassword" type="password" placeholder="Password"
                  value="{{.Mqtt.Password}}">
              </div>
            </div>
          </div>
          <div class="form-check">
            <label class="form-check-label">
              <input class="form-check-input" id="mqttUseTls" type="checkbox" {{if .Mqtt.UseTls}} checked {{end}}>
              Use TLS
            </label>
          </div>
        </div> <!-- network -->
        <div role="tabpanel" class="tab-pane" id="providers">
          <table class="table table-sm table-striped" id="providerList">
            <thead>
              <tr>
                <th>Type</th>
                <th>Path</th>
                <th>Speed</th>
              </tr>
            </thead>
            <tbody id="providerListBody">
            </tbody>
          </table>
          <div class="from-group">
            <button class="btn btn-primary" data-toggle="modal" data-target="#interfaceModal">
              Add Interface
            </button>
          </div>
        </div> <!-- providers -->
      </div> <!-- tab-content -->
      <hr>
      <div class="form-group">
        <button class="btn btn-success btn-lg btn-block">Update Settings</button>
      </div>
    </form>
  </div>
  <div class="col-sm-2">
    <label>Configuration</label>
    <ul class="nav nav-pills flex-column" role="tablist">
      <li role="presentation" class="nav-item">
        <a href="#basic-setup" aria-controls="basic-setup" role="tab" data-toggle="tab" class="nav-link active">
          Vessel Info
        </a>
      </li>
      <li role="presentation" class="nav-item">
        <a href="#network" aria-controls="network" role="tab" data-toggle="tab" class="nav-link">Network</a>
      </li>
      <li role="presentation" class="nav-item">
        <a href="#providers" aria-controls="providers" role="tab" data-toggle="tab" class="nav-link">
          Providers
        </a>
      </li>
    </ul>
  </div>
</div>

<div id="interfaceModal" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="interfaceModalLabel"
  aria-hidden="true">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="interfaceModalLabel">Add Interface</h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <form id="addInterface">
          <div class="form-group">
            <label class="control-label" for="interfaceType">
              Interface Type
            </label>
            <select class="form-control" id="interfaceType">
              <option>Select Type</option>
              <option></option>
              <option value="filestreamNmea0183">NMEA 0183 Log File</option>
              <option value="filestreamNmea2000">NMEA 2000 Log File</option>
              <option value="gpsd">GPSd Server</option>
              <option value="actisense">Actisense NGT-1</option>
              <option value="canusb">Lawicel CAN-USB</option>
            </select>
          </div>
          <div class="form-group" id="devicePathGroup" style="display: hidden">
            <label for="devicePath">Source Device</label>
            <input class="form-control" id="devicePath" type="text" />
          </div>
          <div class="form-group" id="fileSelectGroup" style="display: hidden">
            <label for="filename">Read From</label>
            <select class="form-control" id="filename" value="">
              <option>Select File</option>
              <option></option>
              <option value=""></option>
            </select>
          </div>
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
        <button type="button" class="btn btn-primary" data-dismiss="modal" onClick="addInterface()">
          Add Interface
        </button>
      </div>
    </div>
  </div>
</div>

<div id="pgnModal" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="pgnModalLabel" aria-hidden="true">
  <div class="modal-dialog" role="document">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" id="pgnModalLabel"></h5>
        <button type="button" class="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="modal-body">
        <p>PGN: <span class='pgn'></span></p>
        <p>Category: <span class='category'></span></p>
        <p>Fields</p>
        <ul class='fields'></ul>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button>
      </div>
    </div>
  </div>
</div>
{{end}}
