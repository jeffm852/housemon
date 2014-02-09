ng = angular.module 'myApp'

ng.config ($stateProvider, navbarProvider) ->
  $stateProvider.state 'jeeboot',
    url: '/jeeboot'
    templateUrl: 'jeeboot/view.html'
    controller: 'JeeBootCtrl'
  navbarProvider.add '/jeeboot', 'JeeBoot', 10

ng.controller 'JeeBootCtrl', ($scope, $timeout, jeebus) ->
  # TODO rewrite these example to use the "hm" service i.s.o. "jeebus"

  # TODO this delay seems to be required to avoid an error with WS setup - why?
  $timeout ->
    $scope.hwid = jeebus.attach '/jeeboot/hwid/'
    $scope.$on '$destroy', -> jeebus.detach '/jeeboot/hwid/'
    $scope.firmware = jeebus.attach '/jeeboot/firmware/'
    $scope.$on '$destroy', -> jeebus.detach '/jeeboot/firmware/'
  , 100

  $scope.onFileDrop = (x) ->
    lastId = Object.keys($scope.firmware).sort().pop() | 0
    lastId = 999  if lastId < 999
    for f in x
      r = new FileReader()
      r.onload = (e) ->
        jeebus.rpc 'savefile', "firmware/#{f.name}", e.target.result
        jeebus.store "/jeeboot/firmware/#{++lastId}", file: f.name
      r.readAsText f

  $scope.fwDel = (swid) ->
    jeebus.store "/jeeboot/firmware/#{swid}"

  $scope.hwDel = (hwid) ->
    jeebus.store "/jeeboot/hwid/#{hwid}"

  $scope.hwSave = (id, field, value) ->
    row = $scope.hwid[id]
    row[field] = value
    jeebus.store "/jeeboot/hwid/#{id}", row

# see http://docs.angularjs.org/guide/forms
ng.directive 'contenteditable', ($parse) ->
  restrict: 'A'
  link: (scope, elm, attr) ->
    if attr.onBlur
      elm.on 'blur', ->
        scope.$apply ->
          fn = $parse attr.onBlur
          fn scope, $value: elm.text()

# see also github.com/danialfarid/angular-file-upload
ng.directive 'onFileDrop', ($parse) ->
  restrict: 'A'
  link: (scope, elem, attr) ->

    elem[0].addEventListener 'dragover', (evt) ->
      evt.stopPropagation()
      evt.preventDefault()
      elem.addClass 'dragActive'

    elem[0].addEventListener 'dragleave', (evt) ->
      elem.removeClass 'dragActive'

    elem[0].addEventListener 'drop', (evt) ->
      evt.stopPropagation()
      evt.preventDefault()
      elem.removeClass 'dragActive'

      fn = $parse attr.onFileDrop
      fl = (x for x in evt.dataTransfer.files)
      fn scope, $files: fl, $event: evt
