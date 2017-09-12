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

  $('#provider-type').change(function(e) {
    if($(e.target).val() === 'filestream') {
      $('#file-select').removeClass('hidden');
    } else {
      $('#file-select').addClass('hidden');
    }
  });

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

  if($('#provider-type').val() === 'filestream') {
    $('#file-select').removeClass('hidden');
  }

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

function submitForm() {
  var params = {
    vessel: {
      name: $('#vesselName').val(),
      manufacturer: $('#vesselBrand').val(),
      model: $('#vesselModel').val(),
      year: Number.parseInt($('#vesselYear').val()),
      registration: $('#registration').val(),
      mmsi: Number.parseInt($('#mmsi').val()),
      callsign: $('#callsign').val(),
      uuid: getUUID()
    },
    connection: {
      listen: $('#listenOn').val(),
      port: Number.parseInt($('#port').val())
    },
    provider_type: $('#provider-type').val(),
    log_file: $('#provider-file').val()
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
