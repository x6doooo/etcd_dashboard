/**
 * Created by dx.yang on 16/9/10.
 */


(function() {
    angular.module('fe')
        .directive('navBar', function() {
            return {
                restrict: 'A',
                templateUrl: 'app/components/navBar/navBar.html',
                controller: 'NavBarCtrl',
                controllerAs: 'nav'
            };
        })
})();
