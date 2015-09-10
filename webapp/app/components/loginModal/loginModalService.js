var loginModalService = angular.module('loginModalService', ['ui.bootstrap', 'loginModalController']);

loginModalService.service('loginModalService', ['$modal', '$rootScope', function ($modal, $rootScope) {

  function assignCurrentUser (user) {
    $rootScope.currentUser = user;
    return user;
  }

  return function() {
    var instance = $modal.open({
      templateUrl: 'app/components/loginModal/loginModalView.html',
      controller: 'loginModalController',
      controllerAs: 'loginModalController'
    })

    return instance.result.then(assignCurrentUser);
  };

}]);