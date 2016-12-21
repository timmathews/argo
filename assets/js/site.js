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

  var ws = new WebSocket("ws://localhost:8082/ws/stats");
  ws.onmessage = function(data) {
    var stats = JSON.parse(data.data);

    $('#stats').empty();

    $.each(stats, function(k, v) {
      $('#stats').append("<tr><th>" + k + "</th><td>" + v + "</td></tr>");
    });
  };
});

function submitForm() {
  var http = new XMLHttpRequest();
  http.open("POST", "/admin", true);
  http.setRequestHeader("Content-type","application/json");

  var params = {
    vessel: {
      name: $('#vesselName').val(),
      manufacturer: $('#vesselBrand').val(),
      model: $('#vesselModel').val(),
      year: $('#vesselYear').val(),
      registration: $('#registration').val(),
      mmsi: $('#mmsi').val(),
      callsign: $('#callsign').val(),
      uuid: getUUID()
    },
    connection: {
      listen: $('#listenOn').val(),
      port: Number.parseInt($('#port').val())
    },
    provider_type: $('#provider-type').val(),
    log_file: $('#provider-file').val()
  }

  console.log(params);

  http.send(JSON.stringify(params));
  http.onload = function() {
    console.log(http.response);
  }

}

function getUUID() {
  var uuid = [
    $('#uuid_0').val(),
    $('#uuid_1').val(),
    $('#uuid_2').val(),
    $('#uuid_3').val(),
    $('#uuid_4').val()
  ];

  return uuid;
}
