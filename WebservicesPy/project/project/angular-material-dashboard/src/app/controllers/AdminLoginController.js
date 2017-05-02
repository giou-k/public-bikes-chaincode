(function(){

  angular
    .module('app')
    .controller('AdminLoginController', [
      '$http',
      '$state',
      AdminLoginController
    ]);

  function AdminLoginController($http, $state) {
    var vm = this;

    vm.user = {
      email: '',
      password: ''
    };

    vm.ui = {};
    vm.ui.errors = '';

    vm.login = function() {
      $http.post('http://localhost:8888/adminlogin/', vm.user)
      .then(function successCallback(response) {
          console.log(response)
          $state.go('home.admin')
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
