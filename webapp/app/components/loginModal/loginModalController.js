var loginModalController = angular.module('loginModalController', []);

loginModalController.controller('loginModalController', ['$scope', '$auth',  function ($scope, $auth) {

  this.cancel = $scope.$dismiss;

  this.submit = function (email, password) {
	var credentials = {
      'email': email,
      'password': password
    };
    return $auth.login(credentials).then(function(data) {
      $scope.$close(data);
    }).catch(function() {
    	$scope.$dismiss()
    });
  };

}]);