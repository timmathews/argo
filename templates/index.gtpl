<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>Argo :: A Signal K Server</title>

    <!-- Bootstrap core CSS -->
    <link href="/assets/css/bootstrap.min.css" rel="stylesheet">

    <!-- Custom styles for this template -->
    <link href="/assets/css/site.css" rel="stylesheet">
  </head>

  <body>
    <header>
      <div class="nav navbar-inverse navbar-fixed-top">
        <div class="container-fluid">
          <div class="navbar-header">
            <h1 class="navbar-brand">Argo <small>A Signal K Server</small></h1>
          </div>
        </div>
      </div>
    </header>
    <div class="container-fluid">
      <div class="row">
        <div class="col-sm-2 col-sm-offset-1">
          <h2>Statistics</h2>
          <table class="table table-condensed table-striped">
            <thead>
              <tr><th>PGN</th><th>Count</th></tr>
            </thead>
            <tbody id="stats"></tbody>
          </table>
        </div>
        <div class="col-sm-4">
          <form onsubmit="return submitForm()">
            <div class="tab-content">
              <div role="tabpanel" class="tab-pane active" id="basic-setup">
                <div class="form-group">
                  <label for="vesselName">Vessel Name</label>
                  <input class="form-control" id="vesselName" placeholder="Name" value="">
                </div>
                <div class="row">
                  <div class="col-xs-4">
                    <div class="form-group">
                      <label for="vesselBrand">Make</label>
                      <input class="form-control" id="vesselBrand" placeholder="Manufacturer"
                      value="">
                    </div>
                  </div>
                  <div class="col-xs-4">
                    <div class="form-group">
                      <label for="vesselModel">Model</label>
                      <input class="form-control" id="vesselModel" placeholder="Model"
                      value="">
                    </div>
                  </div>
                  <div class="col-xs-4">
                    <div class="form-group">
                      <label for="vesselYear">Year</label>
                      <input class="form-control" id="vesselYear" placeholder="Year" type="number"
                      value="">
                    </div>
                  </div>
                </div>
                <div class="row">
                  <div class="col-xs-4">
                    <div class="form-group">
                      <label for="mmsi">MMSI</label>
                      <span class="glyphicon glyphicon-question-sign" data-toggle="tooltip" data-placement="top"
                        title="Maritime Mobile Service Identity, assigned regionally, leave blank if you do not have
                        one"></span>
                      <input class="form-control" id="mmsi" placeholder="MMSI" value="">
                    </div>
                  </div>
                  <div class="col-xs-4">
                    <div class="form-group">
                      <label for="callsign">Callsign</label>
                      <span class="glyphicon glyphicon-question-sign" data-toggle="tooltip" data-placement="top"
                        title="Your radio callsign"></span>
                      <input class="form-control" id="callsign" placeholder="Callsign"
                      value="">
                    </div>
                  </div>
                  <div class="col-xs-4">
                    <div class="form-group">
                      <label for="registration">Registration</label>
                      <span class="glyphicon glyphicon-question-sign" data-toggle="tooltip" data-placement="top"
                        title="Your boat's registration number, typically displayed on the side of the hull"></span>
                      <input class="form-control" id="registration" placeholder="Registration"
                      value="">
                    </div>
                  </div>
                </div>
                <div class="form-group">
                  <label for="port">UUID</label><a class="pull-right" id="genUUID" href="#">Generate new UUID</a>
                  <span class="glyphicon glyphicon-question-sign" data-toggle="tooltip" data-placement="top"
                    title="The UUID is used to uniquely identify your boat in Signal K, if you have one from a
                    previous installation of Signal K enter it here, otherwise click 'Generate new UUID to have one
                    created for you"></span>
                  <div class="row narrow">
                    <div class="col-xs-3">
                      <input class="form-control" id="uuid_0" value="">
                    </div>
                    <div class="col-xs-2">
                      <input class="form-control" id="uuid_1" value="">
                    </div>
                    <div class="col-xs-2">
                      <input class="form-control" id="uuid_2" value="">
                    </div>
                    <div class="col-xs-2">
                      <input class="form-control" id="uuid_3" value="">
                    </div>
                    <div class="col-xs-3">
                      <input class="form-control" id="uuid_4" value="">
                    </div>
                  </div>
                </div>
              </div> <!-- basic-setup -->
              <div role="tabpanel" class="tab-pane" id="network">
                <div class="row">
                  <div class="col-xs-8">
                    <div class="form-group">
                      <label for="listenOn">Listen On</label>
                      <input class="form-control" id="listenOn" placeholder="IP Address" value="">
                    </div>
                  </div>
                  <div class="col-xs-4">
                    <div class="form-group">
                      <label for="port">Port</label>
                      <input class="form-control" id="port" type="number" placeholder="Port" value="">
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
          <ul class="nav nav-pills nav-stacked" role="tablist">
            <li role="presentation" class="active">
              <a href="#basic-setup" aria-controls="basic-setup" role="tab" data-toggle="tab">Vessel Info</a>
            </li>
            <li role="presentation">
              <a href="#network" aria-controls="network" role="tab" data-toggle="tab">Network</a>
            </li>
            <li role="presentation">
              <a href="#providers" aria-controls="providers" role="tab" data-toggle="tab">Providers</a>
            </li>
          </ul>
        </div>
      </div>
    </div>

    <footer>
      <div class="navbar navbar-inverse navbar-fixed-bottom">
        <div class="container-fluid">
          <ul class="nav navbar-nav">
            <li><a href="https://github.com/timmathews/argo">GitHub</a></li>
            <li><a href="http://signalk.org">Signal K</a></li>
          </li>
          <p class="navbar-text">
            &copy; 2016 Tim Mathews, Licensed
            <a href="http://www.gnu.org/licenses/gpl-3.0.html" class="navbar-link">GPLv3</a>
          </p>
        </div>
      </div>
    </footer>

    <!-- Bootstrap core JavaScript
    ================================================== -->
    <!-- Placed at the end of the document so the pages load faster -->
    <script src="/assets/js/jquery.min.js"></script>
    <script src="/assets/js/bootstrap.min.js"></script>
    <!-- IE10 viewport hack for Surface/desktop Windows 8 bug -->
    <script src="/assets/js/site.js"></script>
  </body>
</html>
