goImbApp.config(function ($stateProvider, $urlRouterProvider, $authProvider) {
  $authProvider.loginUrl = '/auth/login';
  $urlRouterProvider.otherwise('/welcome');
  $stateProvider
    .state('welcome', {
      url: '/welcome',
      templateUrl: 'app/components/welcome/welcomeView.html',
      data: {
        requireLogin: false
      }
    })
    .state('app', {
      abstract: true,
      url: '/app',
      template: '<ui-view/>',
      data: {
        requireLogin: true
      }
    })
    .state('app.thread', {
      url: '/thread',
      templateUrl: 'app/components/thread/threadView.html'
    });
});
