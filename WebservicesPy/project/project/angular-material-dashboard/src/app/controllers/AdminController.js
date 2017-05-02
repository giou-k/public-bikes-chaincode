(function(){

  angular
    .module('app')
    .controller('AdminController', [
      '$http',
      'leafletData',
      '$scope',
      AdminController
    ])
    .controller('UploadController', [
      '$scope',
      '$http',
      UploadController
    ]);


  function AdminController($http, leafletData,$scope) {
    $scope.center = {
      lat: 37.987,
      lng: 23.692,
      zoom: 8,
      station: '',
      passname: ''
    };
    $scope.$watch('files.length',function(newVal,oldVal){
        console.log($scope.files);
    });

    $scope.saveStation = function() {
    $http.post('http://localhost:8888/stationinput/', $scope.center)
        console.log(response)
        $scope.center = {
          lat: 37.987,
          lng: 23.692,
          zoom: 8,
          station: '',
          passname: ''
        };
      }, function errorCallback(response) {
        $scope.center = {
          lat: 37.987,
          lng: 23.692,
          zoom: 8,
          station: '',
          passname: ''
        };
    };
  }

function UploadController($scope, $http){
      $scope.submit = function(){
        console.log ("Upload Single Image to Stamplay successful!");
          var formData = new FormData();
          angular.forEach($scope.files,function(obj){
              formData.append('files[]', obj.lfFile, 'file1');
          });
          formData.append('station_name', $scope.station_name)
          formData.append('rupos', $scope.rupos)
          $http.post('http://localhost:8888/upload/', formData, {
              transformRequest: angular.identity,
              headers: {'Content-Type': undefined}
          }).then(function(result){
              // do sometingh
          },function(err){
            console.log ("Upload Single Image to Stamplay successful!");
              // do sometingh
          });
      };
  }

})();
