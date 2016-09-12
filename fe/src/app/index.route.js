(function () {
    'use strict';

    angular
        .module('fe')
        .config(routerConfig);

    /** @ngInject */
    function routerConfig($stateProvider, $urlRouterProvider) {
        $stateProvider
            .state('home', {
                controller: 'HomeCtrl',
                controllerAs: 'home',
                url: '/',
                templateUrl: 'app/components/home/home.html'
            })
            .state('endpoints', {
                abstact: true,
                controller: 'EndpointsIndexCtrl',
                controllerAs: 'edpts',
                url: '/endpoints',
                templateUrl: 'app/components/endpoints/index.html'
            })
            .state('endpoints.list', {
                url: '/list',
                views: {
                    subView: {
                        controller: 'EndpointsListCtrl',
                        controllerAs: 'vm',
                        templateUrl: 'app/components/endpoints/subViews/list.html'
                    }
                }
            })
            // .state('manage', {
            //     abstract: true,
            //     controller: 'MainController',
            //     controllerAs: 'main',
            //     url: '/m',
            //     templateUrl: 'app/main/main.html',
            // })
            // .state('manage.list', {
            //     url: '/list',
            //     views: {
            //         subView: {
            //             controller: 'MainListController',
            //             controllerAs: 'vm',
            //             templateUrl: 'app/main-list/list.html'
            //         }
            //     }
            // });

        $urlRouterProvider.otherwise('/');
    }

})();
