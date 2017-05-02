(function(){

  angular
    .module('app')
    .controller('UserController', [
      '$http',
      '$state',
      '$scope',
      UserController
    ]);

  function UserController($http, $state, $scope) {
    $scope.user = {email: 'anz@test.com'};
    var config = {
      params: $scope.user,
      headers : {'Accept' : 'application/json'}
    };
    $http.get('http://localhost:8888/user/', config)
        console.log(response),function errorCallback(response) {
      };

    $scope.tableData = [];
    $scope.getItems = function() {
    $http.post('http://localhost:8888/user/', $scope.user)
      }, function errorCallback(response) {
    };

    tableService
      .loadAllItems()
      .then(function(tableData) {
        $scope.tableData = [].concat(tableData);
      });


  }

})();
