(function(){

  angular
    .module('app')
    .controller('RegisterController', [
      '$http',
      '$state',
      RegisterController
    ]);

  function RegisterController($http, $state) {
    var vm = this;

    vm.user = {
      email: '',
      password: ''
    };

    vm.ui = {};
    vm.ui.errors = '';

    vm.register = function() {
      $http.post('http://localhost:8888/register/', vm.user)
      .then(function successCallback(response) {
          console.log(response)
          $state.go('login')
        }, function errorCallback(response) {
          vm.ui.errors = "Registration failed"
        });
    }
  }

})();
