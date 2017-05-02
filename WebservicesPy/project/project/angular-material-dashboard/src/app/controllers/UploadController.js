// (function(){
//
//   angular
//     .module('app')
//     .controller('UploadController', [
//       '$scope',
//       '$http',
//       UploadController
//     ]);
//
//
// function UploadController($scope, $http){
//       $scope.submit = function(){
//         console.log ("Upload!");
//           var formData = new FormData();
//           angular.forEach($scope.files,function(obj){
//               formData.append('files[]', obj.lfFile, 'file1');
//           });
//           formData.append('station_name', $scope.station_name)
//           formData.append('rupos', $scope.rupos)
//           $http.post('http://localhost:8888/upload/', formData, {
//               transformRequest: angular.identity,
//               headers: {'Content-Type': undefined}
//           }).then(function(result){
//               // do sometingh
//           },function(err){
//             console.log ("Upload!");
//               // do sometingh
//           });
//       };
//   }
//  })();
