(function(){

  angular
    .module('app')
    .controller('LoginController', [
      '$http',
      '$state',
      LoginController
    ]);

  function LoginController($http, $state) {
    var vm = this;

    vm.user = {
      email: '',
      password: ''
    };

    vm.ui = {};
    vm.ui.errors = '';

    vm.login = function() {
      $http.post('http://localhost:8888/auth/login/', vm.user)
      .then(function successCallback(response) {
          console.log(response)
          $state.go('home.admin')
          $state.Params = vm.user
        }, function errorCallback(response) {
          vm.ui.errors = "Login failed"
          vm.user = {
            email: '',
            password: ''
          };
        });
    }
  }

})();
