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
