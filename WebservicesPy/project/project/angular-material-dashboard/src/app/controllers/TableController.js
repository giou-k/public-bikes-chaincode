(function(){

  angular
    .module('app')
    .controller('TableController', [
      '$scope',
      '$http',
      TableController
    ]);

  function TableController($scope, $http) {
    var vm = this;
    $http.get('http://localhost:8888/stations/').then(function(data){
      console.log(data)
      $scope.stations = data.data
    })
  }

})();
