var loginModalController = angular.module('loginModalController', ['satellizer']);

loginModalController.controller('loginModalController', ['$scope', '$auth', '$http',  function ($scope, $auth, $http) {

  this.cancel = $scope.$dismiss;

  this.submit = function (email, password) {
	var credentials = {
      'email': email,
      'password': password
    };
    return $auth.login(credentials).then(function(data) {
      return $http({ method: 'GET', url: '/api/me' });
    }).then(function(data) {
      if (200 == data.status) {
        $scope.$close(data.data);
      } else {
        $scope.$dismiss();  
      }
    }).catch(function() {
    	$scope.$dismiss();
    });
  };
}]);