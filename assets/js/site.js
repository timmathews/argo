/*!
 * IE10 viewport hack for Surface/desktop Windows 8 bug
 * Copyright 2014-2015 Twitter, Inc.
 * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)
 */

// See the Getting Started docs for more information:
// http://getbootstrap.com/getting-started/#support-ie10-width

(function () {
  'use strict';

  if (navigator.userAgent.match(/IEMobile\/10\.0/)) {
    var msViewportStyle = document.createElement('style')
    msViewportStyle.appendChild(
      document.createTextNode(
        '@-ms-viewport{width:auto!important}'
      )
    )
    document.querySelector('head').appendChild(msViewportStyle)
  }
})();

$(function() {
  $('[data-toggle="tooltip"]').tooltip();

  $('#interfaceType').change(handleInterfaceTypeChange);

  $('.appInstall').on('submit', handleAppInstall);

  $('#genUUID').click(function(e) {
    e.preventDefault();

    $.get('/admin/uuid', function(data) {
      var uuid = data.uuid;
      console.log(uuid);
      $('#uuid_0').val(uuid[0]);
      $('#uuid_1').val(uuid[1]);
      $('#uuid_2').val(uuid[2]);
      $('#uuid_3').val(uuid[3]);
      $('#uuid_4').val(uuid[4]);
    });
  });

  var host = window.location.host;
  var proto = window.location.protocol === 'http:' ? 'ws' : 'wss';
  var ws = new WebSocket(proto + "://" + host + "/ws/stats");
  ws.onmessage = function(data) {
    var stats = JSON.parse(data.data);

    $('#stats').empty();

    $.each(stats, function(k, v) {
      $('#stats').append('<tr><th><a href="#" data-pgn="' + k + '">' + k + "</th><td>" + v + "</td></tr>");
    });

    $('a[data-pgn]').on('click', function(e) { showPgnData($(e.target).data('pgn')); });
  };
});

function showPgnData(pgn) {
  $.getJSON('/signalk/v1/api/messages/' + pgn, function(data) {
    console.log(data.FieldList.map(x => x.Name));
    $('#pgnModal .modal-title').text(data.Description);
    $('#pgnModal .modal-body .pgn').text(data.Pgn);
    $('#pgnModal .modal-body .category').text(data.Category);
    $('#pgnModal .fields').text('');
    data.FieldList.forEach(function(x) {
      $('#pgnModal .fields').append('<li>' + x.Name + '</li>');
    });
    $('#pgnModal').modal();
  });
}

function handleInterfaceTypeChange(e) {
  $('#fileSelectGroup').hide();
  $('#devicePathGroup').hide();
  $('#deviceBaudGroup').hide();
  $('#gpsdPortGroup').hide();

  var intType = $(e.target).val();
  console.log(intType);

  if (intType === 'filestreamNmea0183' || intType === 'filestreamNmea2000') {
    $('#fileSelectGroup').show();
  } else if (intType === 'gpsd') {
    $('#gpsdPortGroup').show();
  } else if (intType === 'actisense' || intType === 'canusb') {
    $('#devicePathGroup').show();
    $('#deviceBaudGroup').show();
  }
}

function addInterface() {
  var intType = $('#interfaceType').val();
  var intPath = $('#devicePath').val();
  var intFile = $('#fileName').val();
  var intBaud = $('#baudRate').val();

  var intName = '';

  if (intType === 'filestreamNmea0183') {
    intName = 'NMEA 0183 Recording';
  } else if (intType === 'filestreamNmea2000') {
    intName = 'NMEA 2000 Recording';
  } else if (intType === 'gpsd') {
    intName = 'GPSd';
  } else if (intType === 'actisense') {
    intName = 'Actisense NGT-1';
  } else if (intType === 'canusb') {
    intName = 'Lawicel CAN-USB';
  }

  var list = document.getElementById('providerListBody');

  var r = document.createElement('tr');
  var d1 = document.createElement('td');
  d1.appendChild(document.createTextNode(intName));
  var d2 = document.createElement('td');
  d2.appendChild(document.createTextNode(intPath));
  var d3 = document.createElement('td');
  d3.appendChild(document.createTextNode(intFile));

  r.appendChild(d1);
  r.appendChild(d2);
  r.appendChild(d3);

  list.appendChild(r);
}

function submitForm() {
  var params = {
    Vessel: {
      Name: $('#vesselName').val(),
      Manufacturer: $('#vesselManufacturer').val(),
      Model: $('#vesselModel').val(),
      Year: Number.parseInt($('#vesselYear').val()),
      Registration: $('#registration').val(),
      Mmsi: Number.parseInt($('#mmsi').val()),
      Callsign: $('#callsign').val(),
      Uuid: getUUID()
    },
    Server: {
      ListenOn: $('#listenOn').val(),
      Port: Number.parseInt($('#port').val()),
      UseTls: $('#useTls').is(':checked'),
      EnableWebsockets: $('#enableWebsockets').is(':checked'),
      Certificate: $('#certificate').val() // Get Base64 encoded file??
    },
    Mqtt: {
      Enable: $('#mqttEnable').is(':checked'),
      UseTls: $('#mqttUseTls').is(':checked'),
      Host: $('#mqttHost').val(),
      Port: Number.parseInt($('#mqttPort').val()),
      ClientId: $('#mqttClientId').val(),
      Username: $('#mqttUsername').val(),
      Password: $('#mqttPassword').val(),
      Channel: $('#mqttChannel').val()
    }
  };

  var http = new XMLHttpRequest();
  http.open("POST", "/admin", true);
  http.setRequestHeader("Content-type", "application/json");
  http.send(JSON.stringify(params));

  return false;
}

function getUUID() {
  var uuid = [
    $('#uuid_0').val(),
    $('#uuid_1').val(),
    $('#uuid_2').val(),
    $('#uuid_3').val(),
    $('#uuid_4').val()
  ];

  return uuid.join('-');
}

function handleAppInstall(e) {
  e.preventDefault();
  var p = $(e.target).find('input[name="name"]').val();
  var v = $(e.target).find('input[name="version"]').val();

  var pkg = {package: p, version: v};

  var http = new XMLHttpRequest();
  http.open("POST", "/apps/install", true);
  http.setRequestHeader("Content-type", "application/json");
  http.send(JSON.stringify(pkg));
}
