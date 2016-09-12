/**
 * Created by dx.yang on 16/9/10.
 */


(function() {

    angular.module('fe').controller('NavBarCtrl', NavBarCtrl);

    /** @ngInject */
    function NavBarCtrl() {
        var nav = this;
        nav.list = [{
            title: 'Endpoints',
            href: 'endpoints'
        }, {
            title: 'Data',
            href: 'data'
        }, {
            title: 'Monitor',
            href: 'monitor'
        }, {
            title: 'Users',
            href: 'users'
        }];

        nav.go = function(n) {

        };
    }

})();