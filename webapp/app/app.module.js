var goImbApp = angular.module('goImbApp', [
  'ui.router',
  'ui.bootstrap', 
  'satellizer',
  'loginModalService'
]);

goImbApp.config(['$httpProvider', function ($httpProvider) {
  $httpProvider.interceptors.push(function ($timeout, $q, $injector) {
    var loginModalService, $http, $state;

    // this trick must be done so that we don't receive
    // `Uncaught Error: [$injector:cdep] Circular dependency found`
    $timeout(function () {
      loginModalService = $injector.get('loginModalService');
      $http = $injector.get('$http');
      $state = $injector.get('$state');
    });

    return {
      responseError: function (rejection) {
        if ((rejection.status !== 401) || (rejection.config.url === '/auth/login')) {
          return rejection;
        }

        var deferred = $q.defer();

        loginModalService()
          .then(function () {
            deferred.resolve( $http(rejection.config) );
          })
          .catch(function () {
            $state.go('welcome');
            deferred.reject(rejection);
          });

        return deferred.promise;
      }
    };
  });
}]);


goImbApp.run(['$rootScope', '$state', '$auth', 'loginModalService', function($rootScope, $state, $auth, loginModalService) {
  $rootScope.logout = function() {
    $auth.logout();
    console.log('lolo');
    $rootScope.currentUser = null;
    $state.go('welcome');
  }

  $rootScope.$on('$stateChangeStart', function (event, toState, toParams) {
    var requireLogin = toState.data.requireLogin;

    if (requireLogin && ((typeof $rootScope.currentUser === 'undefined') 
      || (null == $rootScope.currentUser))) {
      event.preventDefault();

      loginModalService()
        .then(function () {
          return $state.go(toState.name, toParams);
        })
        .catch(function () {
          return $state.go('welcome');
        });
    }
  });
}]);


