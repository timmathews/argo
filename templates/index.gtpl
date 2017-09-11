<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <title>Argo :: A Signal K Server</title>

    <!-- Bootstrap core CSS -->
    <link href="/assets/css/bootstrap.min.css" rel="stylesheet">

    <!-- Custom styles for this template -->
    <link href="/assets/css/site.css" rel="stylesheet">
  </head>

  <body>
    <nav class="navbar navbar-expand-md navbar-dark bg-dark fixed-top">
      <a class="navbar-brand" href="/">Argo <small>A Signal K Server</small></a>
      <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarLinks"
        aria-controls="navbarLinks" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarLinks">
        <ul class="navbar-nav mr-auto">
          <li class="nav-item active">
            <a class="nav-link" href="#">Configuration <span class="sr-only">(current)</span></a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#">Apps</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#">Plugins</a>
          </li>
        </ul>
      </div>
    </nav>
    <div class="container">
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
          <form onsubmit="return submitForm()">
            <div class="tab-content">
              <div role="tabpanel" class="tab-pane active" id="basic-setup">
                <div class="form-group">
                  <label for="vesselName">Vessel Name</label>
                  <input class="form-control" id="vesselName" placeholder="Name" value="{{.Name}}">
                </div>
                <div class="form-row">
                  <div class="col">
                    <div class="form-group">
                      <label for="vesselBrand">Make</label>
                      <input class="form-control" id="vesselBrand" placeholder="Manufacturer" value="{{.Make}}">
                    </div>
                  </div>
                  <div class="col">
                    <div class="form-group">
                      <label for="vesselModel">Model</label>
                      <input class="form-control" id="vesselModel" placeholder="Model" value="{{.Model}}">
                    </div>
                  </div>
                  <div class="col">
                    <div class="form-group">
                      <label for="vesselYear">Year</label>
                      <input class="form-control" id="vesselYear" placeholder="Year" type="number" value="{{.Year}}">
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
                      <input class="form-control" id="mmsi" placeholder="MMSI" value="{{.Mmsi}}">
                    </div>
                  </div>
                  <div class="col">
                    <div class="form-group">
                      <label for="callsign">Callsign</label>
                      <img src="/assets/svg/info.svg" class="octicon" data-toggle="tooltip" data-placement="top"
                        title="Your radio callsign" />
                      <input class="form-control" id="callsign" placeholder="Callsign" value="{{.Callsign}}">
                    </div>
                  </div>
                  <div class="col">
                    <div class="form-group">
                      <label for="registration">Registration</label>
                      <img src="/assets/svg/info.svg" class="octicon" data-toggle="tooltip" data-placement="top"
                        title="Your boat's registration number, typically displayed on the side of the hull" />
                      <input class="form-control" id="registration" placeholder="Registration"
                        value="{{.Registration}}">
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
                      <input class="form-control" id="uuid_0" value="{{.Uuid}}">
                    </div>
                    <div class="col-2">
                      <input class="form-control" id="uuid_1" value="{{.Uuid}}">
                    </div>
                    <div class="col-2">
                      <input class="form-control" id="uuid_2" value="{{.Uuid}}">
                    </div>
                    <div class="col-2">
                      <input class="form-control" id="uuid_3" value="{{.Uuid}}">
                    </div>
                    <div class="col-3">
                      <input class="form-control" id="uuid_4" value="{{.Uuid}}">
                    </div>
                  </div>
                </div>
              </div> <!-- basic-setup -->
              <div role="tabpanel" class="tab-pane" id="network">
                <div class="form-row">
                  <div class="col-8">
                    <div class="form-group">
                      <label for="listenOn">Listen On</label>
                      <input class="form-control" id="listenOn" placeholder="IP Address" value="">
                    </div>
                  </div>
                  <div class="col-4">
                    <div class="form-group">
                      <label for="port">Port</label>
                      <input class="form-control" id="port" type="number" placeholder="Port" value="">
                    </div>
                  </div>
                </div>
                <div class="form-check">
                  <label class="form-check-label">
                    <input class="form-check-input" type="checkbox">
                    Enable WebSockets
                  </label>
                </div>
                <div class="form-check">
                  <label class="form-check-label">
                    <input class="form-check-input" type="checkbox">
                    Enable TLS
                  </label>
                </div>
                <div class="form-group">
                  <label for="certificate">Install Certificate</label>
                  <input class="form-control-file" id="certificate" type="file" placeholder="Certificate">
                </div>
                <hr>
                <div class="form-check">
                  <label class="form-check-label">
                    <input class="form-check-input" type="checkbox">
                    Enable MQTT
                  </label>
                </div>
                <div class="form-row">
                  <div class="col-8">
                    <div class="form-group">
                      <label for="mqttHost">MQTT Server</label>
                      <input class="form-control" id="mqttHost" placeholder="Hostname" value="">
                    </div>
                  </div>
                  <div class="col-4">
                    <div class="form-group">
                      <label for="mqttPort">Port</label>
                      <input class="form-control" id="mqttPort" type="number" placeholder="Port" value="">
                    </div>
                  </div>
                </div>
                <div class="form-row">
                  <div class="col">
                    <div class="form-group">
                      <label for="mqttClientId">Client ID</label>
                      <input class="form-control" id="mqttClientId" placeholder="Client ID" value="">
                    </div>
                  </div>
                  <div class="col">
                    <div class="form-group">
                      <label for="mqttUsername">Username</label>
                      <input class="form-control" id="mqttUsername" placeholder="Client ID" value="">
                    </div>
                  </div>
                  <div class="col">
                    <div class="form-group">
                      <label for="mqttPassword">Password</label>
                      <input class="form-control" id="mqttPassword" type="password" placeholder="Password" value="">
                    </div>
                  </div>
                </div>
                <div class="form-check">
                  <label class="form-check-label">
                    <input class="form-check-input" type="checkbox">
                    Use TLS
                  </label>
                </div>
                <hr>
                <div class="form-row">
                  <div class="col-4">
                    <div class="form-group">
                      <label for="mqttPort">Port</label>
                      <input class="form-control" id="mqttPort" type="number" placeholder="Port" value="">
                    </div>
                  </div>
                </div>
              </div> <!-- network -->
              <div role="tabpanel" class="tab-pane" id="providers">
                <div class="form-group">
                  <label for="provider-type">Provider Type</label>
                  <select class="form-control" id="provider-type" value="">
                    <option>Select Type</option>
                    <option></option>
                    <option value="filestream-nmea0183">NMEA 0183 Log File</option>
                    <option value="filestream-nmea2000">NMEA 2000 Log File</option>
                    <option value="gpsd">GPSd Server</option>
                    <option value="actisense">Actisense NGT-1</option>
                  </select>
                </div>
                <div class="form-group hidden" id="file-select">
                  <label for="provider-file">Read From</label>
                  <select class="form-control" id="provider-file" value="">
                    <option>Select File</option>
                    <option></option>
                    <option value=""></option>
                  </select>
                </div>
              </div> <!-- providers -->
            </div> <!-- tab-content -->
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
    </div>

    <nav class="navbar navbar-expand-md navbar-dark bg-dark fixed-bottom">
      <div class="navbar-nav mr-auto">
        <a class="nav-item nav-link" href="https://github.com/timmathews/argo">GitHub</a>
        <a class="nav-item nav-link" href="http://signalk.org">Signal K</a>
      </div>
      <span class="navbar-text">
        &copy; 2016 Tim Mathews, Licensed
        <a href="http://www.gnu.org/licenses/gpl-3.0.html">GPLv3</a>
      </span>
    </nav>

    <div id="pgnModal" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="pgnModalLabel"
      aria-hidden="true">
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
            <button type="button" class="btn
              btn-secondary"
              data-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Bootstrap core JavaScript
    ================================================== -->
    <!-- Placed at the end of the document so the pages load faster -->
    <script src="/assets/js/jquery.min.js"></script>
    <script src="/assets/js/popper.min.js"></script>
    <script src="/assets/js/bootstrap.min.js"></script>
    <!-- IE10 viewport hack for Surface/desktop Windows 8 bug -->
    <script src="/assets/js/site.js"></script>
  </body>
</html>
