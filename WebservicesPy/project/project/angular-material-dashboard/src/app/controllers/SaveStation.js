(function(){

  angular
    .module('app')
    .controller('SaveStation', [
      '$http',
      SaveStation
    ]);

  function SaveStation($http) {
    var vm = this;

    vm.station = {
      name: '',
      passname: '',
      lat: '',
      long: ''
    };

    vm.ui = {};

    vm.saveStation = function() {
      $http.post('http://localhost:8888/stationinput/', vm.station)
      .then(function successCallback(response) {
          console.log(response)
        }, function errorCallback(response) {
          vm.station = {
            name: '',
            passname: '',
            lat: '',
            long: ''
          };
        });
    }
  }

})();
