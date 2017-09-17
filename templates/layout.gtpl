{{define "layout"}}
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
          <li class="nav-item">
            <a class="nav-link" href="/">Configuration</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="/apps">Apps</a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#">Plugins</a>
          </li>
        </ul>
      </div>
    </nav>
    <div class="container">
      {{template "content" .}}
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
{{end}}
