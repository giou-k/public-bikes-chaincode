/**
 * Maps Controller
 */

angular.module('RDash')
    .controller('MapsCtrl', ['$scope', '$cookieStore', MapsCtrl]);

function MapsCtrl($scope, $cookieStore) {
$scope.map = { center: { latitude: 45, longitude: -73 }, zoom: 8 };
}
