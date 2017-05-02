(function(){

  angular
    .module('app')
    .controller('TableUsersController', [
      '$scope',
      '$http',
      TableUsersController
    ]);

  function TableUsersController($scope, $http) {
    var vm = this;
    $http.get('http://localhost:8888/users/').then(function(data){
      console.log(data)
      $scope.users = data.data
    })
  }

})();
